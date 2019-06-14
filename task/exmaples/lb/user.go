package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

// User is a task that interacts with the load balancer.
type User struct {
	lb            phi.Sender
	numResults    int
	resultsNeeded int
	done          chan struct{}
}

// NewUser returns a new `User` along with a channel that will be closed when
// the user has finished. Finishing is determined by receiving `resultsNeeded`
// responses from the load balancer.
func NewUser(lb phi.Sender, resultsNeeded int) (User, <-chan struct{}) {
	done := make(chan struct{})
	return User{
		lb:            lb,
		numResults:    0,
		resultsNeeded: resultsNeeded,
		done:          done,
	}, done
}

// Reduce implements the `phi.Reducer` interface. Upon receiving an `Init`
// message, the user will then send this to the load balancer. Once receiving
// the corresponding result, it will update the number of results it has seen.
// Once it has seen `resultsNeeded` responses, it will close the done channel,
// signalling that it has finished.
func (user *User) Reduce(self phi.Task, message phi.Message) phi.Message {
	switch message.(type) {
	case Init:
		user.sendAsync(self, message)
	case Done:
		user.numResults++
		if user.numResults >= user.resultsNeeded {
			close(user.done)
		}
	default:
		panic(fmt.Sprintf("unexpected message type %T", message))
	}

	return nil
}

// sendAsync sends a message and asynchronously waits for the response. It will
// ensure that the message is sent.
func (user *User) sendAsync(self phi.Task, message phi.Message) {
	go func() {
		responder, ok := user.lb.Send(message)
		// Ensure that the message is sent
		for !ok {
			time.Sleep(10 * time.Millisecond)
			responder, ok = user.lb.Send(message)
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
