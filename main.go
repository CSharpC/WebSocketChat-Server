package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var router *Router

func main() {
	loadDB("./db.sqlite")
	router = &Router{
		make(map[string]Sendable),
		make(chan Sendable),
		make(chan RouterQuery),
	}
	router.Run()
	for _, c := range getChannelsList() {
		ch := &Channel{
			make(chan Message),
			c.ID,
			c.Name,
			make(map[string]Sendable),
			make(chan Sendable),
			make(chan Sendable),
			make(chan bool),
		}
		router.addSendable <- ch
		ch.Run()
	}
	http.HandleFunc("/ws", createWs)
	http.HandleFunc("/channels/list", channelList)
	err := http.ListenAndServe(":29107", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed upgrading connection", err)
	}
	id := r.URL.Query().Get("client_id")
	if !knownID(id) {
		log.Println("Unknown User:", id)
		w.WriteHeader(403)
		return
	}
	cl := &Client{
		send:       make(chan Message),
		router:     router,
		conn:       c,
		id:         id,
		disconnect: make(chan bool),
	}
	router.addSendable <- cl
	go cl.readWs()
	cl.writeWs()
	log.Println("Ending", cl.ID(), "goroutine")
}
