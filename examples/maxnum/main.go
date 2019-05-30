package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/renproject/phi"
)

func main() {
	numPlayers := uint(3)
	max := uint(100)
	routeTable := map[uint][]uint{}
	players := make([]Player, numPlayers)
	playerTasks := make([]phi.Task, numPlayers)

	// Make players
	for i := uint(0); i < numPlayers; i++ {
		num := uint(rand.Intn(int(max)))
		players[i] = NewPlayer(i, num, numPlayers)
		playerTasks[i] = phi.NewTask(&players[i], 2*int(numPlayers))

		// Ring topology
		var prev, next uint
		if i == 0 {
			prev, next = numPlayers-1, 1
		} else if i == numPlayers-1 {
			prev, next = numPlayers-2, 0
		} else {
			prev, next = i-1, i+1
		}
		routeTable[i] = []uint{prev, next}
	}

	// Make router
	playerMap := map[uint]phi.Task{}
	results := make(chan Result, 1)
	for i, player := range players {
		playerMap[player.ID()] = playerTasks[i]
	}
	router := NewRouter(routeTable, playerMap, results)
	routerTask := phi.NewTask(&router, int(numPlayers))
	router.SetTask(routerTask)

	// Start the tasks
	done := context.Background()
	for _, player := range playerTasks {
		go player.Run(done)
	}
	go routerTask.Run(done)

	// Send the initial message
	routerTask.Send(Begin{})

	result := <-results
	if result.success {
		fmt.Printf("Success: %v players reached consensus on a maximum value of %v\n", result.players, result.max)
	} else {
		fmt.Println("Failed!")
	}
}
