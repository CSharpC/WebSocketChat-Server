package main

import (
	"log"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
)

func main() {
	var hub *Hub
	hub = newHub()
	go hub.run()
	room := newRoom("room1", getUUID())
	go room.run()
	hub.addRoom <- room
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandle(hub, w, r)
	})
	err := http.ListenAndServe(":25565", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getUUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}
