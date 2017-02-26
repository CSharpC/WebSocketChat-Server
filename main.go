package main

import (
	"log"
	"net/http"
)

func main() {
	var rooms *Room
	rooms = newRoom()
	go rooms.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandle(rooms, w, r)
	})
	err := http.ListenAndServe(":25565", nil)
	if err != nil {
		log.Fatal(err)
	}
}
