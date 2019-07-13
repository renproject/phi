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

// Messages is a collection of messages. Note that handlers will never receive
// a `Messages` type because they will always be flattened before being
// processed. If a task needs to respond with more than one message, and it is
// important that these messages are processed together, then a custom
// container type should be created and handled appropriately in the handler.
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
	// Send takes a message and returns a boolean that indicates if the send was
	// successful.
	Send(Message) bool
}

// Task is the intersection of the `Runner` and `Sender` interfaces. It
// represents an entity that (when running) can be sent messages and upon
// receipt of these messages performs internal logic (which often involves
// sending messages to other tasks).
type Task interface {
	Runner
	Sender
}

// Handler defines a type that can receive a message and mutate its internal
// state. The `Task` argument is the parent task of the Handler, and can be used
// by a Handler to send messages to itself.
type Handler interface {
	Handle(Task, Message)
}

// Router represents something that can route different messages to different
// senders. Returning a nil sender from `Route` signifies that the message is
// not to be sent anywhere.
type Router interface {
	Route(Message) Sender
}

// Options are passed when constructing a `Task`. The `Cap` is the buffer
// capacity of the `Task`'s channel, and the `Scale` is the number of worker
// instances of the handler for load balancing. If `Scale` is an number less
// than 2, there will only be one instance of the handler. It is important to
// note that additional copies of the handler will not be created for `Scale`
// >= 2; this means that handlers that have and modify their own state are not
// safe to be used at non-unity scales. Only handlers that are purely
// functional should be used with non-unity scale.
type Options struct {
	Cap, Scale int
}

// task is a basic implementation for a `Task`.
type task struct {
	// The handler for message handling logic.
	handler Handler

	// The buffered channel that incoming messages are written to.
	input chan Message

	// The scale (number of workers) for the task.
	scale int
}

// New returns a new task with the given handler and buffer capacity. The
// buffer capacity is the number of messages that can be buffered for
// processing before the task can no longer accept more messages (until space
// in the buffer is freed up by processing messages in the buffer).
func New(handler Handler, opts Options) Task {
	return &task{
		handler: handler,
		input:   make(chan Message, opts.Cap),
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
				task.handle(flatten(message))
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
// currently has a full buffer and won't accept the message.
func (task *task) Send(m Message) bool {
	select {
	case task.input <- m:
		return true
	default:
		return false
	}
}

// handle a message sent to the Task. It is assumed that the message is
// flattened.
func (task *task) handle(m Message) {
	switch m := m.(type) {
	case Messages:
		for _, msg := range m {
			task.handler.Handle(task, msg)
		}
	default:
		task.handler.Handle(task, m)
	}
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
	rMu *sync.Mutex
	r   Router
}

// NewRouter returns a new sender that represents a Router. The given Router
// determines how the sender routes messages; any message `m` that is sent to
// this sender will be sent to the sender determined by the Router through
// `Route(m)`.
func NewRouter(r Router) Sender {
	return &router{
		rMu: new(sync.Mutex),
		r:   r,
	}
}

// Send implements the `Sender` interface. If the resolver returns a nil Sender,
// it signifies that the message is not to be sent anywhere.
func (r *router) Send(message Message) bool {
	sender := func() Sender {
		r.rMu.Lock()
		defer r.rMu.Unlock()
		return r.r.Route(message)
	}()
	if sender != nil {
		return sender.Send(message)
	}
	return true
}
