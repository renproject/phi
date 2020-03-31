package task

import (
	"context"
)

// Request is a marker interface. It requires no explicit functionality, and
// only exists so that programmers cannot accidentally send types as requests
// that were not intended to be requests.
type Request interface {
	// IsRequest is a marker function. It should do nothing.
	IsRequest()
}

// Response is a marker interface. It requires no explicit functionality, and
// only exists so that programmers cannot accidentally send types as responses
// that were not intended to be responses.
type Response interface {
	// IsResponse is a marker function. It should do nothing.
	IsResponse()
}

// A Message is used for inter-Task communication. It carries an internal
// Request that will be delivered to the Task.
type Message interface {
	// Request returns the inner Request carried by this Message.
	Request() Request

	// Respond to this Message. Respond should be called at least once on each
	// Message, but must be called at most once on any Message. It is
	// non-blocking, and safe for concurrent use.
	Respond(Response)

	// Wait for a Response to this Message, until the Context is done. Any error
	// returned from the Context will be returned. Waiting will block the
	// current goroutine. Waiting is safe for concurrent use, but will only
	// return a non-nil Response to at most one caller.
	Wait(context.Context) (Response, error)
}
