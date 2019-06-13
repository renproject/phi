package main

import (
	"fmt"
	"time"

	"github.com/renproject/phi"
)

type LB struct {
	id int
}

func (LB) Reduce(_ phi.Task, _ phi.Message) phi.Message {
	fmt.Println("Doing something!")
	time.Sleep(time.Second)
	return nil
}
