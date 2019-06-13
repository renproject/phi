package main

type MessageA struct{}

func (MessageA) IsMessage() {}

type MessageB struct{}

func (MessageB) IsMessage() {}

type MessageC struct{}

func (MessageC) IsMessage() {}

type Response struct {
	msg string
}

func (Response) IsMessage() {}
