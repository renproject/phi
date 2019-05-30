package main

import "github.com/renproject/phi"

type Player struct {
	num, id, players, currentMax uint
	seen                         map[uint]uint
}

func NewPlayer(id, num, players uint) Player {
	seen := map[uint]uint{}
	seen[id] = num
	return Player{num: num, id: id, players: players, currentMax: num, seen: seen}
}

func (player *Player) ID() uint {
	return player.id
}

func (player *Player) Reduce(message phi.Message) phi.Message {
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
				return MessageBatch{[]phi.Message{message, done}}
			}
			return message
		}
		return nil
	default:
		panic("unexpected message type")
	}
}
