package main

import (
	"context"
	"time"

	"github.com/renproject/phi"
)

func main() {
	pinger := NewPerpetualPinger()
	pingerTask := phi.NewTask(&pinger, 1)
	ponger := NewPonger()
	pongerTask := phi.NewTask(ponger, 1)

	pinger.CompleteSetup(pongerTask, pingerTask)
	done := context.Background()
	go pingerTask.Run(done)
	go pongerTask.Run(done)

	pingerTask.Send(Begin{})
	time.Sleep(10 * time.Second)
}
