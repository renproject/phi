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

// Handle implements the `phi.Handler` interface.
func (destA *DestA) Handle(_ phi.Task, message phi.Message) {
	switch m := message.(type) {
	case MessageA:
		m.Responder <- Response{msg: destA.name}
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

// Handle implements the `phi.Handler` interface.
func (destB *DestB) Handle(_ phi.Task, message phi.Message) {
	switch m := message.(type) {
	case MessageB:
		m.Responder <- Response{msg: destB.name}
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

// Handle implements the `phi.Handler` interface.
func (destC *DestC) Handle(_ phi.Task, message phi.Message) {
	switch m := message.(type) {
	case MessageC:
		m.Responder <- Response{msg: destC.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
