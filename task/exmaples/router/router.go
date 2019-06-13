package main

import (
	"fmt"

	"github.com/renproject/phi"
)

type Router struct {
	destA, destB, destC phi.Sender
}

func NewRouter(a, b, c phi.Sender) Router {
	return Router{
		destA: a,
		destB: b,
		destC: c,
	}
}

func (router *Router) Resolve(message phi.Message) phi.Sender {
	switch message.(type) {
	case MessageA:
		return router.destA
	case MessageB:
		return router.destB
	case MessageC:
		return router.destC
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
}
