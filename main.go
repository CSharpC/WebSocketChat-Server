package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/guglicap/chatServer/db"
)

var router *Router

func main() {
	db.Init("./db.sqlite")
	router = &Router{
		make(map[string]Sendable),
		make(chan Sendable),
		make(chan RouterQuery),
	}
	router.Run()
	for _, c := range db.GetChannelsList() {
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
	http.HandleFunc("/channels/list", RESTChanList)
	http.HandleFunc("/user/name", Name)
	http.HandleFunc("/user/channels", userChannels)
	http.HandleFunc("/signup", Signup)
	err := http.ListenAndServe(":29107", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createWs(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !db.KnownToken(token) {
		log.Println("Unknown token:", token)
		w.WriteHeader(403)
		return
	}
	id := db.GetUserField(token, "ID")
	if len(id) < 1 {
		log.Println("Unknown user.")
		return
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed upgrading connection", err)
	}
	cl := &Client{
		send:       make(chan Message),
		router:     router,
		conn:       c,
		id:         id,
		disconnect: make(chan bool),
	}
	router.addSendable <- cl
	result := make(chan Sendable)
	for _, c := range db.GetUserChannels(token) {
		router.query <- RouterQuery{c.ID, TypeChannel, result}
		ch := <-result
		if ch == nil {
			log.Println("couldn't find channel:", c.ID)
			continue
		}
		ch.(*Channel).Subscribe() <- cl
	}
	go cl.readWs()
	cl.writeWs()
	log.Println("Ending", cl.ID(), "goroutine")
}
