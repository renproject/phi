package task

import (
	"context"
	"sync"

	"github.com/renproject/phi/co"
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
// corresponding result. The `Task` argument is the parent task of the reducer,
// and can be used as a handle to send messages to a reducer's own task.
type Reducer interface {
	Reduce(Task, Message) Message
}

// Resolver represents something that can resolve or route messages to certain
// senders.
type Resolver interface {
	Resolve(Message) Sender
}

// Options are passed when constructing a `Task`. The `Cap` is the buffer
// capacity of the `Task`'s channel, and the `Scale` is the number of worker
// instances of the reducer for load balancing. If `Scale` is an number less
// than 2, there will only be one instance of the reducer. It is important to
// note that additional copies of the reducer will not be created for `Scale`
// >= 2; this means that reducers that have and modify their own state are not
// safe to be used at non-unity scales. Only reducers that are purely
// functional should be used with non-unity scale.
type Options struct {
	Cap, Scale int
}

// task is a basic implementation for a `Task`.
type task struct {
	// The reducer for message handling logic.
	reducer Reducer

	// The buffered channel that incoming messages are written to.
	input chan messageWithResponder

	scale int
}

// New returns a new task with the given reducer and buffer capacity. The
// buffer capacity is the number of messages that can be buffered for
// processing before the task can no longer accept more messages (until space
// in the buffer is freed up by processing messages in the buffer).
func New(reducer Reducer, opts Options) Task {
	return &task{
		reducer: reducer,
		input:   make(chan messageWithResponder, opts.Cap),
		scale:   opts.Scale,
	}
}

// Run implements the `Runner` interface (in order to implement the `Task`
// interface). This function blocks. The task will continue to run until it is
// signalled to terminate by the context.
func (task *task) Run(ctx context.Context) {
	loop := func() {
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

	// Don't spawn a go routine if there is no load balancing
	if task.scale < 2 {
		loop()
	} else {
		co.ParForAll(task.scale, func(i int) { loop() })
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
			messages = append(messages, task.reducer.Reduce(task, msg))
		}
	default:
		messages = append(messages, task.reducer.Reduce(task, message))
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

// router is an implementation of a `Sender` that is a resolver.
type router struct {
	resolverMu *sync.Mutex
	resolver   Resolver
}

// NewResolver returns a new sender that represents a resolver. The given
// resolver determines how the sender routes messages; any message `m` that is
// sent to this sender will be sent to the sender determined by the resolver
// through `Resolve(m)`.
func NewResolver(resolver Resolver) Sender {
	return &router{
		resolverMu: new(sync.Mutex),
		resolver:   resolver,
	}
}

// Send implements the `Sender` interface.
func (r *router) Send(message Message) (<-chan Messages, bool) {
	sender := func() Sender {
		r.resolverMu.Lock()
		defer r.resolverMu.Unlock()
		return r.resolver.Resolve(message)
	}()
	return sender.Send(message)
}
