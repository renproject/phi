package main

import (
	"time"

	"github.com/renproject/phi"
)

// LB is a simple load balancing task.
type LB struct {
	id int
}

// Handle implements the `phi.Handlerr` interface. We simulate a slow task by
// simply sleeping for a time before returning.
func (LB) Handle(_ phi.Task, m phi.Message) {
	if m, ok := m.(Init); ok {
		time.Sleep(time.Second)
		m.Responder <- Done{}
	}
}
