package main

import (
	"context"
	"time"

	"github.com/renproject/phi"
)

func main() {
	opts := phi.Options{Cap: 1}
	a := NewDestA("Alice")
	aTask := phi.New(&a, opts)
	b := NewDestB("Bob")
	bTask := phi.New(&b, opts)
	c := NewDestC("Charlie")
	cTask := phi.New(&c, opts)
	router := NewRouter(aTask, bTask, cTask)
	resolver := phi.NewResolver(&router)
	user := NewUser(resolver)
	opts.Cap = 3
	userTask := phi.New(&user, opts)

	done := context.Background()
	go aTask.Run(done)
	go bTask.Run(done)
	go cTask.Run(done)
	go userTask.Run(done)

	var ok bool
	_, ok = userTask.Send(MessageA{})
	if !ok {
		panic("could not send message to router")
	}
	_, ok = userTask.Send(MessageB{})
	if !ok {
		panic("could not send message to router")
	}
	_, ok = userTask.Send(MessageC{})
	if !ok {
		panic("could not send message to router")
	}
	time.Sleep(10 * time.Millisecond)
}
