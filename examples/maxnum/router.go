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
	task                phi.Task
	players             map[uint]phi.Task
	routeTable          map[uint][]uint
	result, resultsSeen uint
	resultWriter        chan Result
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
	}, resultWriter
}

// SetTask sets the reference to the `phi.Task` which has the given `Router` as
// a `phi.Reducer`. It is necessary for this method to be called for the task
// to work correctly.
func (router *Router) SetTask(task phi.Task) {
	router.task = task
}

// Reduce implements the `phi.Reducer` interface.
func (router *Router) Reduce(message phi.Message) phi.Message {
	// A nil message means that nothing needs to be routed
	if message == nil {
		return nil
	}

	switch message := message.(type) {
	case Begin:
		for _, player := range router.players {
			router.sendAsync(player, message)
		}
	case PlayerNum:
		for _, to := range router.routeTable[message.from] {
			router.sendAsync(router.players[to], message)
		}
	case Done:
		router.resultsSeen++
		if router.resultsSeen == 1 {
			router.result = message.max
		} else {
			if message.max != router.result {
				router.resultWriter <- Result{Success: false}
				router.terminate()
			} else if router.resultsSeen == uint(len(router.players)) {
				router.resultWriter <- Result{Max: router.result, Players: uint(len(router.players)), Success: true}
				router.terminate()
			}
		}
	case phi.MessageBatch:
		for _, msg := range message.Messages {
			router.Reduce(msg)
		}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}

	return nil
}

func (router *Router) terminate() {
	for _, player := range router.players {
		player.Terminate()
	}
	router.task.Terminate()
}

func (router *Router) sendAsync(player phi.Task, message phi.Message) {
	go func() {
		response, ok := player.SendSync(message)
		if ok {
			// Retry until the message is successfully sent
			for !router.task.Send(response) {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}
