package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	rooms    *Room
	conn     *websocket.Conn
	send     chan Message
	ID       string
	nickname string
}

func wsHandle(rooms *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	var clInfo ClientInfo
	err = conn.ReadJSON(&clInfo)
	if err != nil {
		return
	}
	client := &Client{
		rooms:    rooms,
		conn:     conn,
		ID:       clInfo.ID,
		nickname: clInfo.Nickname,
		send:     make(chan Message, 10),
	}
	log.Println(client)
	rooms.subscribe <- client
	go client.writeWs()
	client.readWs()
}

func (c *Client) readWs() {
	defer func() {
		c.conn.Close()
		c.rooms.unsubscribe <- c
		close(c.send)
	}()
	for {
		var m Message
		err := c.conn.ReadJSON(&m)
		if err != nil {
			break
		}
		m.From = c.nickname
		c.rooms.inMsg <- m
	}
}

func (c *Client) writeWs() {
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
