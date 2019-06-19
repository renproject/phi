package main

// MessageA represents a message that has particular destination task.
type MessageA struct{}

// IsMessage implements the `phi.Message` interface.
func (MessageA) IsMessage() {}

// MessageB represents a message that has particular destination task.
type MessageB struct{}

// IsMessage implements the `phi.Message` interface.
func (MessageB) IsMessage() {}

// MessageC represents a message that has particular destination task.
type MessageC struct{}

// IsMessage implements the `phi.Message` interface.
func (MessageC) IsMessage() {}

// Response represents a response message from one of the destinations.
type Response struct {
	msg string
}

// IsMessage implements the `phi.Message` interface.
func (Response) IsMessage() {}
