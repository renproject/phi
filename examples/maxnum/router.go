package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

type Result struct {
	max, players uint
	success      bool
}

type Router struct {
	task            phi.Task
	players         map[uint]phi.Task
	routeTable      map[uint][]uint
	result, numSeen uint
	resultWriter    chan Result
}

func NewRouter(routeTable map[uint][]uint, players map[uint]phi.Task, resultWriter chan Result) Router {
	return Router{
		players:      players,
		routeTable:   routeTable,
		result:       0,
		numSeen:      0,
		resultWriter: resultWriter,
	}
}

func (router *Router) SetTask(task phi.Task) {
	router.task = task
}

func (router *Router) Reduce(message phi.Message) phi.Message {
	// Drop nil messages
	if message == nil {
		return nil
	}

	switch message := message.(type) {
	case Begin:
		for _, player := range router.players {
			go func(player phi.Task, message phi.Message) {
				response, ok := player.SendSync(message)
				if ok {
					for !router.task.Send(response) {
						fmt.Println("WHOOPS BEGIN")
						time.Sleep(10 * time.Millisecond)
					}
				}
			}(player, message)
		}
	case PlayerNum:
		router.SendAsync(message.from, message)
	case Done:
		router.numSeen++
		if router.numSeen == 1 {
			router.result = message.max
		} else {
			if message.max != router.result {
				router.resultWriter <- Result{success: false}
				router.terminate()
			}
			if router.numSeen == uint(len(router.players)) {
				router.resultWriter <- Result{max: router.result, players: uint(len(router.players)), success: true}
				router.terminate()
			}
		}
		fmt.Printf("player %v terminated with result %v\n", message.player, message.max)
	case MessageBatch:
		for _, msg := range message.messages {
			router.Reduce(msg)
		}
	default:
		panic("unexpected message type")
	}

	return nil
}

func (router *Router) SendAsync(from uint, message phi.Message) {
	for _, to := range router.routeTable[from] {
		go func(to uint) {
			response, ok := router.players[to].SendSync(message)
			if ok {
				for !router.task.Send(response) {
					fmt.Println("WHOOPS ASYNC")
					time.Sleep(10 * time.Millisecond)
				}
			}
		}(to)
	}
}

func (router *Router) terminate() {
	for _, player := range router.players {
		player.Terminate()
	}
	router.task.Terminate()
}
