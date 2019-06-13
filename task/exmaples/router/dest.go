package main

import (
	"fmt"

	"github.com/renproject/phi"
)

// DestA represents one type of destination task.
type DestA struct {
	name string
}

// NewDestA creates a new `DestA` with the given name.
func NewDestA(name string) DestA {
	return DestA{name: name}
}

// Reduce implements the `phi.Reducer` interface.
func (destA *DestA) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case MessageA:
		return Response{msg: destA.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

// DestB represents one type of destination task.
type DestB struct {
	name string
}

// NewDestB creates a new `DestB` with the given name.
func NewDestB(name string) DestB {
	return DestB{name: name}
}

// Reduce implements the `phi.Reducer` interface.
func (destB *DestB) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case MessageB:
		return Response{msg: destB.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

// DestC represents one type of destination task.
type DestC struct {
	name string
}

// NewDestC creates a new `DestC` with the given name.
func NewDestC(name string) DestC {
	return DestC{name: name}
}

// Reduce implements the `phi.Reducer` interface.
func (destC *DestC) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case MessageC:
		return Response{msg: destC.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
