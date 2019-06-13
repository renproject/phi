package main

// Begin is a signal for the router start the execution.
type Begin struct{}

// IsMessage implements the `phi.Message` interface.
func (Begin) IsMessage() {}

// PlayerNum is the message that players pass between eachother to disseminate
// what numbers each player has.
type PlayerNum struct {
	from, player, num uint
}

// IsMessage implements the `phi.Message` interface.
func (PlayerNum) IsMessage() {}

// Done is the message that a player will send when they are confident they
// know what the maximum number is.
type Done struct {
	player, max uint
}

// IsMessage implements the `phi.Message` interface.
func (Done) IsMessage() {}
