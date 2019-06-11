// This example runs a simple ping pong interaction.
package main

import (
	"context"
	"time"

	"github.com/renproject/phi"
)

func main() {
	// Create the pinger and ponger tasks
	ponger := NewPonger()
	pongerTask := phi.NewTask(&ponger, 1)
	pinger := NewPerpetualPinger(pongerTask)
	pingerTask := phi.NewTask(&pinger, 1)

	// Run the tasks
	done := context.Background()
	go pingerTask.Run(done)
	go pongerTask.Run(done)

	// Start the communication and run for some time
	pingerTask.Send(Begin{})
	time.Sleep(10 * time.Second)
}
