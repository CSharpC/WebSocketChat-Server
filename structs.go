package main

type Sendable interface {
	Send() chan Message
	ID() string
	Type() int
	Disconnect() chan bool
}

type Message struct {
	Type    int
	To      string
	From    string
	Content string
}
