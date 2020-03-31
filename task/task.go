package task

import (
	"context"
)

// Runner represents something that can be run.
type Runner interface {
	// Run by convention will be blocking. The context should be used to signal
	// the task to terminate.
	Run(context.Context)
}

// A Sender is any type that can be sent a Message.
type Sender interface {
	// Send takes a message and attempts to deliver it until the Context is
	// done. If the Message is successfully delivered before the Context is
	// done, then a nil error is returned. Otherwise, the Context error is
	// returned.
	Send(context.Context, Message) error
}

// A Task is any type that can receive Messages while running in another
// goroutine. It implements Runner (so that it can be run in a background
// goroutine) and Sender (so that it can be sent Messages).
type Task interface {
	Runner
	Sender
}

// A Handler is any type that can receive a Message and mutate its own internal
// state. It is also expected to respond to the Message.
type Handler interface {
	Handle(Message)
}

// A simple implementation of the Task interface. It is sufficient for most
// purposes.
type task struct {
	// message handling logic.
	handler Handler

	// input channel for buffering messages that need to be handled.
	input chan Message
}

// New Task with the given Handler and buffer capacity. The buffer capacity is
// the number of Messages that can be buffered for processing before the Task
// can no longer accept more Messages.
func New(handler Handler, capacity int) Task {
	return &task{
		handler: handler,
		input:   make(chan Message, capacity),
	}
}

// Run the Task until the Context is done. It will read Messages from an
// internal channel that is filled by calling Send. This method will block the
// current goroutine until the Context is done.
func (task *task) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-task.input:
			task.handler.Handle(message)
		}
	}
}

// Send a Message to this Task.
func (task *task) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case task.input <- message:
		return nil
	}
}
