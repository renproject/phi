package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

// PerpetualPinger represents an object that sends pings, but in addition will
// send another ping after each pong it receives.
type PerpetualPinger struct {
	ponger, task phi.Task
}

// NewPerpetualPinger returns a new `PerpetualPinger`.
func NewPerpetualPinger() PerpetualPinger {
	return PerpetualPinger{nil, nil}
}

// CompleteSetup givens the `PerpetualPinger` references to a `Ponger` task as
// well as its own parent task. If this method is not called, the
// `PerpetualPinger` will not function correctly.
func (pinger *PerpetualPinger) CompleteSetup(ponger, task phi.Task) {
	pinger.ponger = ponger
	pinger.task = task
}

// Reduce implements the `phi.Reducer` interface.
func (pinger *PerpetualPinger) Reduce(message phi.Message) phi.Message {
	switch message.(type) {
	case Begin:
		fmt.Println("Pinger beginning...")
		pinger.pingAsync()
		return nil
	case Pong:
		fmt.Println("Received Pong!")
		time.Sleep(500 * time.Millisecond)
		pinger.pingAsync()
		return nil
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

func (pinger *PerpetualPinger) pingAsync() {
	go func() {
		pong, ok := pinger.ponger.SendSync(Ping{})
		if !ok {
			panic("failed to send ping")
		}
		pinger.task.Send(pong)
	}()
}

// Ponger represents an object that sends a pong on receipt of a ping.
type Ponger struct{}

// NewPonger returns a new `Ponger` object.
func NewPonger() Ponger {
	return Ponger{}
}

// Reduce implements the `phi.Reducer` interface.
func (ponger *Ponger) Reduce(message phi.Message) phi.Message {
	switch message.(type) {
	case Ping:
		fmt.Println("Received Ping!")
		time.Sleep(500 * time.Millisecond)
		return Pong{}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
