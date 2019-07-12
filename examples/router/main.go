// This example shows a basic use case of a resolver. The purpose of a resolver
// is to route incoming messages to one of several tasks depending on the type
// of message. In this example, there is a user that can send three different
// kinds of messages to its resolver. The resolver will then route the message
// to one of the three destination tasks that correspond to the message. These
// destination tasks send back a response with a message indicating where the
// response is coming from.
package main

import (
	"context"
	"os"
	"time"

	"github.com/renproject/phi"
)

func main() {
	// Construct the destination tasks
	opts := phi.Options{Cap: 1}
	a := NewDestA("Alice")
	aTask := phi.New(&a, opts)
	b := NewDestB("Bob")
	bTask := phi.New(&b, opts)
	c := NewDestC("Charlie")
	cTask := phi.New(&c, opts)

	// Construct the resolver
	router := NewRouter(aTask, bTask, cTask)
	resolver := phi.NewRouter(&router)

	// Construct the user. Use an increased channel capacity to avoid dropped
	// messages.
	opts.Cap = 3
	user, success := NewUser(resolver)
	userTask := phi.New(&user, opts)

	// Start the tasks. Notice that a resolver is just a `phi.Sender`, and
	// hence does not need to be (and indeed cannot be) run.
	ctx := context.Background()
	go aTask.Run(ctx)
	go bTask.Run(ctx)
	go cTask.Run(ctx)
	go userTask.Run(ctx)

	// Send a message that should be routed to each of the three destinations.
	var ok bool
	ok = userTask.Send(MessageA{Responder: make(chan phi.Message, 1)})
	if !ok {
		panic("could not send message to router")
	}
	ok = userTask.Send(MessageB{Responder: make(chan phi.Message, 1)})
	if !ok {
		panic("could not send message to router")
	}
	ok = userTask.Send(MessageC{Responder: make(chan phi.Message, 1)})
	if !ok {
		panic("could not send message to router")
	}

	// Wait a moment for the responses to get back to the user.
	select {
	case result := <-success:
		if !result {
			os.Exit(1)
		}
		return
	case <-time.After(10 * time.Second):
		os.Exit(1)
	}
}
