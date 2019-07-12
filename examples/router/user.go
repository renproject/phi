package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

// User represents a basic task that sends messages to a resolver and waits for
// the response from the associated destination task that the message was
// intended for. The task is considered to have completed successfully if it
// receives three unique responses after receiving one each of MessageA,
// MessageB and MessageC.
type User struct {
	resolver      phi.Sender
	responsesSeen map[string]struct{}
	success       chan bool
	terminated    bool
}

// NewUser creates a new user with the given resolver.
func NewUser(resolver phi.Sender) (User, <-chan bool) {
	success := make(chan bool, 1)
	return User{
		resolver:      resolver,
		responsesSeen: map[string]struct{}{},
		success:       success,
		terminated:    false,
	}, success
}

// Handle implements the `phi.Handler` interface.
func (user *User) Handle(self phi.Task, message phi.Message) {
	if user.terminated {
		return
	}

	switch message := message.(type) {
	case MessageA:
		user.sendAsync(self, message, message.Responder)
	case MessageB:
		user.sendAsync(self, message, message.Responder)
	case MessageC:
		user.sendAsync(self, message, message.Responder)
	case Response:
		fmt.Printf("received response from %v\n", message.msg)
		if _, ok := user.responsesSeen[message.msg]; ok {
			// If we have already seen a message, this is a failure
			user.success <- false
			user.terminated = true
			return
		} else {
			user.responsesSeen[message.msg] = struct{}{}
			if len(user.responsesSeen) == 3 {
				// Receiving three unique messages is a success
				user.success <- true
				user.terminated = true
				return
			}
		}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
	return
}

// sendAsync sends a message and asynchronously waits for the response. It will
// ensure that the message is sent.
func (user *User) sendAsync(self phi.Task, message phi.Message, responder chan phi.Message) {
	go func() {
		ok := user.resolver.Send(message)
		// Ensure that the message is sent
		for !ok {
			time.Sleep(10 * time.Millisecond)
			ok = user.resolver.Send(message)
		}
		m := <-responder
		ok = self.Send(m)
		// Ensure that the responses get received
		for !ok {
			time.Sleep(10 * time.Millisecond)
			ok = self.Send(m)
		}
	}()
}
