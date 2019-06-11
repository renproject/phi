package main

import (
	"fmt"

	"github.com/renproject/phi"
)

// Player represents one player in the consensus algorithm.
type Player struct {
	num, id, players, currentMax uint
	seen                         map[uint]uint
}

// NewPlayer returns a new `Player` with ID `id` and internal number `num` in
// an execution where there are `players` other players. The ID represents the
// index for routing.
func NewPlayer(id, num, players uint) Player {
	seen := map[uint]uint{}
	seen[id] = num
	return Player{num: num, id: id, players: players, currentMax: num, seen: seen}
}

// ID returns the routing ID for the given player.
func (player *Player) ID() uint {
	return player.id
}

// Reduce implements the `phi.Reducer` interface.
func (player *Player) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message := message.(type) {
	case Begin:
		return PlayerNum{from: player.id, player: player.id, num: player.num}
	case PlayerNum:
		if _, ok := player.seen[message.player]; !ok {
			player.seen[message.player] = message.num
			message.from = player.id
			if message.num > player.currentMax {
				player.currentMax = message.num
			}
			if uint(len(player.seen)) == player.players {
				done := Done{player: player.id, max: player.currentMax}
				return phi.Messages{message, done}
			}
			return message
		}
		return nil
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
