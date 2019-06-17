package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

// User represents a basic task that sends messages to a resolver and waits for
// the response from the associated destination task that the message was
// intended for.
type User struct {
	resolver phi.Sender
}

// NewUser creates a new user with the given resolver.
func NewUser(resolver phi.Sender) User {
	return User{resolver: resolver}
}

// Reduce implements the `phi.Reducer` interface.
func (user *User) Reduce(self phi.Task, message phi.Message) phi.Message {
	switch message := message.(type) {
	case MessageA, MessageB, MessageC:
		user.sendAsync(self, message)
	case Response:
		fmt.Printf("received response from %v\n", message.msg)
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}
	return nil
}

// sendAsync sends a message and asynchronously waits for the response. It will
// ensure that the message is sent.
func (user *User) sendAsync(self phi.Task, message phi.Message) {
	go func() {
		responder, ok := user.resolver.Send(message)
		// Ensure that the message is sent
		for !ok {
			time.Sleep(10 * time.Millisecond)
			responder, ok = user.resolver.Send(message)
		}
		messages := <-responder
		for _, m := range messages {
			_, ok = self.Send(m)
			// Ensure that the responses get received
			for !ok {
				time.Sleep(10 * time.Millisecond)
				_, ok = self.Send(m)
			}
		}
	}()
}
