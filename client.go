package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	send       chan Message
	router     *Router
	conn       *websocket.Conn
	disconnect chan bool
	id         string
}

func (c *Client) Send() chan Message {
	return c.send
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Type() int {
	return TYPE_CLIENT
}

func (c *Client) readWs() {
	defer func() {
		c.conn.Close()
		close(c.send)
		close(c.disconnect)
	}()
	q := make(chan Sendable)
	for {
		var m Message
		err := c.conn.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.disconnect <- true
				<-c.disconnect
				return
			}
			log.Println(err)
			continue
		}
		c.router.query <- RouterQuery{m.To, TYPE_ANY, q}
		s := <-q
		if s == nil {
			log.Println("Can't find sendable:", m.To)
			continue
		}
		s.Send() <- m
	}
}

func (c *Client) Disconnect() chan bool {
	return c.disconnect
}

func (c *Client) writeWs() {
	for m := range c.send {
		err := c.conn.WriteJSON(m)
		if err != nil {
			log.Println("Can't send message:", err)
		}
	}
}
