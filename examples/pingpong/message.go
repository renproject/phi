package main

import "github.com/renproject/phi"

// Begin signals the pinger to start by sending a ping to the ponger.
type Begin struct {
	Responder chan phi.Message
}

// IsMessage implements the `phi.Message` interface.
func (Begin) IsMessage() {}

// Ping represents a ping.
type Ping struct {
	Responder chan phi.Message
}

// IsMessage implements the `phi.Message` interface.
func (Ping) IsMessage() {}

// Pong represents a pong.
type Pong struct{}

// IsMessage implements the `phi.Message` interface.
func (Pong) IsMessage() {}
