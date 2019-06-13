package main

import (
	"fmt"

	"github.com/renproject/phi"
)

type DestA struct {
	name string
}

func NewDestA(name string) DestA {
	return DestA{name: name}
}

func (destA *DestA) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case MessageA:
		return Response{msg: destA.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

type DestB struct {
	name string
}

func NewDestB(name string) DestB {
	return DestB{name: name}
}

func (destB *DestB) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case MessageB:
		return Response{msg: destB.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}

type DestC struct {
	name string
}

func NewDestC(name string) DestC {
	return DestC{name: name}
}

func (destC *DestC) Reduce(_ phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case MessageC:
		return Response{msg: destC.name}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
