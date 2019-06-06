package phi

import (
	"context"
)

// Message represents a message that is sent between tasks for communication.
// It is an enum style interface where a struct can be made a variant of the
// enum by implementing the interface.
type Message interface {
	IsMessage()
}

// messageWithResponder is a `Message` wrapper that also contains a channel
// where a response can be written.
type messageWithResponder struct {
	// The message being sent.
	message Message

	// The responder channel where the task will write its response.
	responder chan Messages
}

// Messages is a collection of messages. Note that reducers will never receive
// a `Messages` type because they will always be flattened before being
// processed. If a task needs to respond with more than one message, and it is
// important that these messages are processed together, then a custom
// container type should be created and handled appropriately in the reducer.
type Messages []Message

// IsMessage implements the Message interface.
func (Messages) IsMessage() {}

// Runner represents something that can be run.
type Runner interface {
	// Run by convention will be blocking. The context should be used to signal
	// the task to terminate.
	Run(context.Context)
}

// Sender represents something that can be sent messages.
type Sender interface {
	// Send takes a message and returns a channel where the response will be
	// written and also a boolean that indicates if the send was successful.
	Send(Message) (<-chan Messages, bool)
}

// Task is the intersection of the `Runner` and `Sender` interfaces. It
// represents an entity that (when running) can be sent messages and upon
// receipt of these messages performs internal logic (which often involves
// sending messages to other tasks) and can also return response messages.
type Task interface {
	Runner
	Sender
}

// Reducer represents something that can receive a message and provide a
// corresponding result.
type Reducer interface {
	Reduce(Message) Message
}

// task is a basic implementation for a `Task`.
type task struct {
	// The reducer for message handling logic.
	reducer Reducer

	// The buffered channel that incoming messages are written to.
	input chan messageWithResponder
}

// NewTask returns a new task with the given reducer and buffer capacity. The
// buffer capacity is the number of messages that can be buffered for
// processing before the task can no longer accept more messages (until space
// in the buffer is freed up by processing messages in the buffer).
func NewTask(reducer Reducer, cap int) Task {
	return &task{
		reducer: reducer,
		input:   make(chan messageWithResponder, cap),
	}
}

// Run implements the `Runner` interface (in order to implement the `Task`
// interface). This function blocks. The task will continue to run until it is
// signalled to terminate by the context.
func (task *task) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-task.input:
			message.responder <- task.reduce(flatten(message.message))
			close(message.responder)
		}
	}
}

// Send implements the `Sender` interface (in order to implement the `Task`
// interface). It returns a channel to which the (possibly nil) response will
// be written, and a bool indicating whether the message was able to be sent;
// true indicates the message was sent, and false indicates that the task
// currently has a full buffer and won't accept the message. In the latter case
// the returned channel will be nil.
func (task *task) Send(message Message) (<-chan Messages, bool) {
	responder := make(chan Messages, 1)
	m := messageWithResponder{message: message, responder: responder}
	select {
	case task.input <- m:
		return responder, true
	default:
		return nil, false
	}
}

// reduce will handle the reduction of a given message for a task. It is
// assumed that the message is flattened. It will always return a `Messages`
// type which contains the responses of the reducer for the given message(s).
// The returned `Messages` is flattened.
func (task *task) reduce(message Message) Messages {
	messages := Messages{}
	switch message := message.(type) {
	case Messages:
		for _, msg := range message {
			messages = append(messages, task.reducer.Reduce(msg))
		}
	default:
		messages = append(messages, task.reducer.Reduce(message))
	}
	return flatten(messages).(Messages)
}

// flatten takes a message and effectively flattens it out to depth 1, where
// the depth is determined by the level of nested Messages. If `flatten`
// receives a message that is not a `Messages`, it will return the same
// message. Otherwise, it will return a `Messages` type where the internal list
// of messages contain no `Messages`; these will be flattened out.
func flatten(message Message) Message {
	switch message := message.(type) {
	case Messages:
		msgs := Messages{}
		for _, msg := range message {
			m := flatten(msg)
			switch m := m.(type) {
			case Messages:
				msgs = append(msgs, m...)
			default:
				msgs = append(msgs, m)
			}
		}
		return msgs
	default:
		return message
	}
}
