package main

// Init is a message that signals a worker to start work.
type Init struct{}

// IsMessage implements the `phi.Message` interface.
func (Init) IsMessage() {}

// Done is a message to signal that a worker is done.
type Done struct{}

// IsMessage implements the `phi.Message` interface.
func (Done) IsMessage() {}
