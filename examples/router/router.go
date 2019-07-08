package main

import (
	"fmt"

	"github.com/renproject/phi"
)

// Router is a simple resolver that will route to three possible destinations.
type Router struct {
	destA, destB, destC phi.Sender
}

// NewRouter creates a new router with the given three destination `Sender`s.
func NewRouter(a, b, c phi.Sender) Router {
	return Router{
		destA: a,
		destB: b,
		destC: c,
	}
}

// Route implements the `phi.Router` interface
func (router *Router) Route(message phi.Message) phi.Sender {
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
