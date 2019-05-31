package phi

import "context"

type Message interface {
	IsMessage()
}

type MessageSync struct {
	message   Message
	responder chan Message
}

type MessageBatch struct {
	Messages []Message
}

func (MessageBatch) IsMessage() {}

type Runner interface {
	Run(context.Context)
	Terminate()
}

type Sender interface {
	Send(Message) bool
	SendSync(Message) (Message, bool)
}

type Task interface {
	Runner
	Sender
}

type task struct {
	reducer   Reducer
	input     chan Message
	inputSync chan MessageSync
	done      chan struct{}
}

func NewTask(reducer Reducer, cap int) Task {
	return &task{
		reducer:   reducer,
		input:     make(chan Message, cap),
		inputSync: make(chan MessageSync, cap),
		done:      make(chan struct{}),
	}
}

func (task *task) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(task.done)
			return
		case <-task.done:
			return
		case message := <-task.input:
			task.reducer.Reduce(message)
		case message := <-task.inputSync:
			message.responder <- task.reducer.Reduce(message.message)
		}
	}
}

func (task *task) Terminate() {
	close(task.done)
}

func (task *task) Send(message Message) bool {
	select {
	case task.input <- message:
		return true
	default:
		return false
	}
}

func (task *task) SendSync(message Message) (Message, bool) {
	responder := make(chan Message, 1)
	select {
	case task.inputSync <- MessageSync{message, responder}:
		select {
		case <-task.done:
			return nil, false
		case response := <-responder:
			return response, true
		}
	default:
		return nil, false
	}
}

type Reducer interface {
	Reduce(Message) Message
}
