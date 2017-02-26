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
	hub      *Hub
	conn     *websocket.Conn
	send     chan Message
	ID       string
	nickname string
}

/* This function is called at every connection on a new goroutine, thus readWs runs on a goroutine of its own
 * and we spawn another one for writeWs
 */
func wsHandle(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	var clInfo ClientInfo
	err = conn.ReadJSON(&clInfo)
	if err != nil || len(clInfo.ID) == 0 {
		log.Println("Couldn't decode client info. Aborting.", err)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		ID:   clInfo.ID,
		send: make(chan Message, 10),
	}
	hub.addClient <- client
	hub.roomSubscribe("room1") <- client
	go client.writeWs()
	client.readWs()
	log.Println("Client", client.nickname, "has gone.")
}

func (c *Client) readWs() {
	defer func() {
		c.hub.removeClient <- c
		c.conn.Close()
		close(c.send)
	}()
	for {
		var m Message
		err := c.conn.ReadJSON(&m)
		if err != nil {
			log.Println("We got an error reading from the client, interrupting communication.")
			break
		}
		c.hub.routeMessage(m.To) <- m
	}
}

func (c *Client) writeWs() {
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("The client send channel got closed, exiting the goroutine! Bye!")
}
