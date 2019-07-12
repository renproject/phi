package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

var WAIT_MILLIS time.Duration = time.Duration(10)

// PerpetualPinger represents an object that sends pings, but in addition will
// send another ping after each pong it receives.
type PerpetualPinger struct {
	ponger                       phi.Task
	pingsReceived, pingsRequired int
	done                         chan struct{}
}

// NewPerpetualPinger returns a new `PerpetualPinger`.
func NewPerpetualPinger(ponger phi.Task, pingsRequired int) (PerpetualPinger, chan struct{}) {
	done := make(chan struct{}, 1)
	return PerpetualPinger{
		ponger:        ponger,
		pingsReceived: 0,
		pingsRequired: pingsRequired,
		done:          done,
	}, done
}

// Handle implements the `phi.Handler` interface.
func (pinger *PerpetualPinger) Handle(self phi.Task, message phi.Message) {
	switch message.(type) {
	case Begin:
		pinger.pingAsync(self)
	case Pong:
		fmt.Println("Received Pong!")
		pinger.pingsReceived++
		if pinger.pingsReceived == pinger.pingsRequired {
			close(pinger.done)
		}
		time.Sleep(WAIT_MILLIS * time.Millisecond)
		pinger.pingAsync(self)
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

// pingAsync sends a message to the ponger and asynchronously wais for the
// response. It does not retry so the message may not get sent.
func (pinger *PerpetualPinger) pingAsync(self phi.Task) {
	responder := make(chan phi.Message, 1)
	ok := pinger.ponger.Send(Ping{Responder: responder})
	if !ok {
		panic("failed to send ping")
	}
	go func() {
		for m := range responder {
			ok := self.Send(m)
			if !ok {
				panic("failed to receive pong")
			}
		}
	}()
}

// Ponger represents an object that sends a pong on receipt of a ping.
type Ponger struct{}

// NewPonger returns a new `Ponger` object.
func NewPonger() Ponger {
	return Ponger{}
}

// Handle implements the `phi.Handler` interface.
func (ponger *Ponger) Handle(_ phi.Task, message phi.Message) {
	switch message := message.(type) {
	case Ping:
		fmt.Println("Received Ping!")
		time.Sleep(WAIT_MILLIS * time.Millisecond)
		message.Responder <- Pong{}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
