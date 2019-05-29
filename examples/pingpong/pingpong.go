package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

type Begin struct{}

func (Begin) IsMessage() {}

type Ping struct{}

func (Ping) IsMessage() {}

type Pong struct{}

func (Pong) IsMessage() {}

type PerpetualPinger struct {
	ponger, task phi.Task
}

func NewPerpetualPinger() PerpetualPinger {
	return PerpetualPinger{nil, nil}
}

func (pinger *PerpetualPinger) CompleteSetup(ponger, task phi.Task) {
	pinger.ponger = ponger
	pinger.task = task
}

func (pinger *PerpetualPinger) Reduce(message phi.Message) phi.Message {
	switch message.(type) {
	case Begin:
		fmt.Println("Pinger beginning...")
		go func() {
			pong, ok := pinger.ponger.SendSync(Ping{})
			if !ok {
				panic("failed to send ping")
			}
			pinger.task.Send(pong)
		}()
		return nil
	case Pong:
		fmt.Println("Received Pong!")
		time.Sleep(500 * time.Millisecond)
		go func() {
			pong, ok := pinger.ponger.SendSync(Ping{})
			if !ok {
				panic("failed to send ping")
			}
			pinger.task.Send(pong)
		}()
		return nil
	default:
		panic("unexpected message type")
	}
}

type Ponger struct{}

func NewPonger() Ponger {
	return Ponger{}
}

func (ponger Ponger) Reduce(message phi.Message) phi.Message {
	switch message.(type) {
	case Ping:
		fmt.Println("Received Ping!")
		time.Sleep(500 * time.Millisecond)
		return Pong{}
	default:
		panic("unexpected message type")
	}
}
