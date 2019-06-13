package main

import (
	"context"
	"time"

	"github.com/renproject/phi"
)

func main() {
	lb := LB{}
	opts := phi.Options{Cap: 1, Scale: 10}
	lbTask := phi.New(lb, opts)
	done := context.Background()
	go lbTask.Run(done)
	for i := 0; i < 10; i++ {
		var ok bool
		_, ok = lbTask.Send(Init{})
		for !ok {
			_, ok = lbTask.Send(Init{})
		}
	}
	time.Sleep(time.Second)
}
