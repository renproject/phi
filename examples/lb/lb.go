package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

// LB is a simple load balancing task.
type LB struct {
	id int
}

// Handle implements the `phi.Handler` interface. We simulate a slow task by
// simply sleeping for a time before returning.
func (LB) Handle(_ phi.Task, m phi.Message) {
	init, ok := m.(Init)
	if !ok {
		panic(fmt.Errorf("unexpected message type=%T", m))
	}
	time.Sleep(time.Second)
	init.Responder <- Done{}
}
