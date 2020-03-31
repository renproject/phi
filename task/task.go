package task

import (
	"context"
	"fmt"
	"reflect"
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

	// SendAndWait is a convenience method that allows for more ergonomic use of
	// the Sender interface. It takes a message and attempts to deliver it until
	// the Context is done. Once the Message has been sent, it will wait for a
	// response and write the response to the given pointer. If the Message is
	// successfully delivered, and a response is received, before the Context is
	// done, then a nil error is returned. Otherwise, the Context error is
	// returned.
	SendAndWait(context.Context, Request, interface{}) error
}

// A Task is any type that can receive Messages while running in another
// goroutine. It implements Runner (so that it can be run in a background
// goroutine) and Sender (so that it can be sent Messages).
type Task interface {
	Runner
	Sender
}

// A simple implementation of the Task interface. It is sufficient for most
// purposes.
type simpleTask struct {
	// message handling logic.
	handler Handler

	// input channel for buffering messages that need to be handled.
	input chan Message
}

// New Task with the given Handler and buffer capacity. The buffer capacity is
// the number of Messages that can be buffered for processing before the Task
// can no longer accept more Messages.
func New(handler Handler, capacity int) Task {
	return &simpleTask{
		handler: handler,
		input:   make(chan Message, capacity),
	}
}

// Run the Task until the Context is done. It will read Messages from an
// internal channel that is filled by calling Send. This method will block the
// current goroutine until the Context is done.
func (task *simpleTask) Run(ctx context.Context) {
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
func (task *simpleTask) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case task.input <- message:
		return nil
	}
}

// SendAndWait will Send a Message to this Task, and then Wait for a response to
// the Message until the Context is done.
func (task *simpleTask) SendAndWait(ctx context.Context, request Request, response interface{}) error {
	message := NewMessage(request)
	if err := task.Send(ctx, message); err != nil {
		return err
	}
	if response == nil {
		panic(fmt.Errorf("expected pointer, got %v", response))
	}
	if reflect.TypeOf(response).Kind() != reflect.Ptr {
		panic(fmt.Errorf("expected pointer, got %T", response))
	}
	interf, err := message.Wait(ctx)
	if err != nil {
		return err
	}
	reflect.Indirect(reflect.ValueOf(response)).Set(reflect.ValueOf(interf))
	return nil
}
