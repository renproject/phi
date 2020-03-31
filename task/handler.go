package task

// A Handler is any type that can receive a Message and mutate its own internal
// state. It is also expected to respond to the Message.
type Handler interface {
	Handle(Message)
}

// A HandlerFunc is a function that directly implements the Handler interface.
type HandlerFunc func(Message)

// Handle implements the Handler interface for the HandlerFunc type.
func (f HandlerFunc) Handle(message Message) {
	f(message)
}
