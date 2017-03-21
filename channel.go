package main

import (
	"log"
)

type Channel struct {
	send       chan Message
	id, name   string
	members    map[string]Sendable
	sub, unsub chan Sendable
	disconnect chan bool
}

func (c *Channel) Send() chan Message {
	return c.send
}

func (c *Channel) ID() string {
	return c.id
}

func (c *Channel) Type() int {
	return TYPE_CHANNEL
}

func (c *Channel) run() {
	for m := range c.send {
		for _, s := range c.members {
			s.Send() <- m
		}
	}
}

func (c *Channel) Run() {
	go c.run()
	go c.unsubscribe()
	go c.subscribe()
}

func (c *Channel) Subscribe() chan Sendable {
	return c.sub
}

func (c *Channel) Unsubscribe() chan Sendable {
	return c.unsub
}

func (c *Channel) Disconnect() chan bool {
	return c.disconnect
}

func (c *Channel) unsubscribe() {
	for s := range c.unsub {
		if _, ok := c.members[s.ID()]; ok {
			delete(c.members, s.ID())
		}
	}
}

func (c *Channel) subscribe() {
	for s := range c.sub {
		if _, ok := c.members[s.ID()]; ok {
			return
		}
		c.members[s.ID()] = s
		go func() {
			<-s.Disconnect()
			log.Println("Deleting", s.ID(), "from", c.ID())
			delete(c.members, s.ID())
			s.Disconnect() <- true
		}()
	}
}
