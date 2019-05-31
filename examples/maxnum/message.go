package main

type Begin struct{}

func (Begin) IsMessage() {}

type PlayerNum struct {
	from, player, num uint
}

func (PlayerNum) IsMessage() {}

type Done struct {
	player, max uint
}

func (Done) IsMessage() {}
