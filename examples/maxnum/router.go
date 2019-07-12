package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

// Result contains the information relevant to a completed execution of the
// tasks.
type Result struct {
	Max, Players uint
	Success      bool
}

// A Router is responsible for routing messages between the players.
type Router struct {
	players             map[uint]phi.Task
	routeTable          map[uint][]uint
	result, resultsSeen uint
	resultWriter        chan Result
	terminated          bool
}

// NewRouter returns a new `Router` along with a channel that the result of the
// execution will be written to. The `routeTable` represents the topology of
// the network; `routeTable[from]` will be the list of indices that the `from`
// player is connected to. These connections are directed: if index 1 is in the
// list `routeTable[0]`, then messages from index 0 will be directly sent to
// index 1, but it is not necessarily the case that message from index 1 will
// be sent to index 0. For this to be the case, index 0 will need to be an
// element of `routeTable[1]`. For the algorithm to terminate, it is required
// that the network is connected.
func NewRouter(routeTable map[uint][]uint, players map[uint]phi.Task) (Router, chan Result) {
	resultWriter := make(chan Result, 1)
	return Router{
		players:      players,
		routeTable:   routeTable,
		result:       0,
		resultsSeen:  0,
		resultWriter: resultWriter,
		terminated:   false,
	}, resultWriter
}

// Handle implements the `phi.Handler` interface.
func (router *Router) Handle(self phi.Task, message phi.Message) {
	// A nil message means that nothing needs to be routed
	if message == nil || router.terminated {
		return
	}

	switch message := message.(type) {
	case BeginRouter:
		for _, player := range router.players {
			responder := make(chan phi.Message, 1)
			router.sendAsync(self, player, Begin{Responder: responder}, responder)
		}
	case PlayerNum:
		for _, to := range router.routeTable[message.from] {
			responder := make(chan phi.Message, 1)
			router.sendAsync(self, router.players[to], PlayerNum{
				from:      message.from,
				player:    message.player,
				num:       message.num,
				Responder: responder,
			}, responder)
		}
	case Done:
		router.resultsSeen++
		if router.resultsSeen == 1 {
			router.result = message.max
		} else {
			if message.max != router.result {
				router.resultWriter <- Result{Success: false}
				router.terminated = true
			} else if router.resultsSeen == uint(len(router.players)) {
				router.resultWriter <- Result{Max: router.result, Players: uint(len(router.players)), Success: true}
				router.terminated = true
			}
		}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

// sendAsync sends a message and asynchronously waits for the response. It will
// ensure that the message is sent.
func (router *Router) sendAsync(self phi.Task, player phi.Task, message phi.Message, responder chan phi.Message) {
	go func() {
		ok := player.Send(message)
		// Ensure that the message is sent
		for !ok {
			time.Sleep(10 * time.Millisecond)
			ok = player.Send(message)
		}
		m := <-responder
		if messages, ok := m.(phi.Messages); ok {
			for _, m := range messages {
				ok = self.Send(m)
				// Ensure that the responses get received
				for !ok {
					time.Sleep(10 * time.Millisecond)
					ok = self.Send(m)
				}
			}
		} else {
			ok = self.Send(m)
			// Ensure that the responses get received
			for !ok {
				time.Sleep(10 * time.Millisecond)
				ok = self.Send(m)
			}
		}
	}()
}
