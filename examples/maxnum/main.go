// This example implements a number consensus algorithm. The aim of the
// algorithm is for a set of players, each with their own internal number, to
// all agree on the maximum number of all the players. They do this by
// communicating with eachother using messages in a connected network. When a
// player decides it knows what the maximum number is, it outputs this number
// and terminates. When all of the players terminate, if they all output the
// same number, and that number was indeed the maximum of the internal numbers,
// the algorithm is considered to have executed correctly. Otherwise, the
// execution was a failure.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/renproject/phi"
)

func main() {
	// Seend RNG
	rand.Seed(time.Now().UTC().UnixNano())

	// Parameters for the algorithm
	numPlayers := uint(3)
	max := uint(1000)

	// Make players
	playerMax := uint(0)
	players := make([]Player, numPlayers)
	playerTasks := make([]phi.Task, numPlayers)
	for i := uint(0); i < numPlayers; i++ {
		num := uint(rand.Intn(int(max)))

		// Keep track of the actual max number
		if num > playerMax {
			playerMax = num
		}

		players[i] = NewPlayer(i, num, numPlayers)
		playerTasks[i] = phi.NewTask(&players[i], 2*int(numPlayers))
	}

	// Make router
	playerMap := map[uint]phi.Task{}
	for i, player := range players {
		playerMap[player.ID()] = playerTasks[i]
	}
	router, results := NewRouter(ringTopology(numPlayers), playerMap)
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

	// Read and check the result
	result := <-results
	if result.Success && result.Max == playerMax {
		fmt.Printf("Success: %v players reached consensus on a maximum value of %v\n", result.Players, result.Max)
	} else {
		fmt.Println("Failed!")
	}
}

func ringTopology(numPlayers uint) map[uint][]uint {
	routeTable := map[uint][]uint{}
	for i := uint(0); i < numPlayers; i++ {
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
	return routeTable
}
