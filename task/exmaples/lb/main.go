// This example demonstrates basic use of a load balancing task. To create a
// task that has more than one worker, we set the `Scale` field in the
// `Options` struct to be at least 2. In this case the workers for the load
// balancer simply sleep for a period before returning a response, which
// simulates slow work. The user will wait to receive a certain number of
// results. The fact that the work is load balanced across the workers means
// that instead of having to wait for each worker to finish in turn serially
// (which would take nunWorkers * workTime) we only need to wait slightly
// longer than the time it takes for a single worker to do the work.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/renproject/phi"
)

func main() {
	// Number of workers in the load balancer, and the number of responses the
	// user will wait for
	n := 100

	// Construct the load balancer and user tasks.
	lbOpts := phi.Options{Cap: n, Scale: n}
	lb := LB{}
	lbTask := phi.New(lb, lbOpts)
	userOpts := phi.Options{Cap: n}
	user, finished := NewUser(lbTask, n)
	userTask := phi.New(&user, userOpts)

	// Run the tasks
	done := context.Background()
	go lbTask.Run(done)
	go userTask.Run(done)

	// Send requests to the user
	start := time.Now()
	for i := 0; i < n; i++ {
		_, ok := userTask.Send(Init{})
		if !ok {
			panic("message should send correctly")
		}
	}

	// Wait until the user has finished
	<-finished
	elapsed := time.Since(start)
	fmt.Printf("processed %v requests in %v\n", n, elapsed)
}
