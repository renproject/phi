package main

import (
	"time"

	"github.com/renproject/phi"
)

// LB is a simple load balancing task.
type LB struct {
	id int
}

// Reduce implements the `phi.Reducer` interface. We simulate a slow task by
// simply sleeping for a time before returning.
func (LB) Reduce(_ phi.Task, _ phi.Message) phi.Message {
	time.Sleep(time.Second)
	return Done{}
}
