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
	room := newChannel(getUUID(), "room1")
	go room.run()
	hub.addChannel <- room

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandle(hub, w, r)
	})
	http.HandleFunc("/channels/list", func(w http.ResponseWriter, r *http.Request) {
		handleRESTCheckID(w, r, handleChannelsList, hub)
	})
	http.HandleFunc("/channels/create", func(w http.ResponseWriter, r *http.Request) {
		handleRESTCheckID(w, r, handleChannelCreate, hub)
	})
	http.HandleFunc("/channels/delete", func(w http.ResponseWriter, r *http.Request) {
		handleRESTCheckID(w, r, handleChannelDelete, hub)
	})
	http.HandleFunc("/user/name", func(w http.ResponseWriter, r *http.Request) {
		handleRESTCheckID(w, r, handleUserName, hub)
	})

	err := http.ListenAndServe(":21197", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getUUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}

func handleRESTCheckID(w http.ResponseWriter, r *http.Request, f RESTHandlerFunc, h *Hub) {
	id := r.URL.Query().Get("client_id")
	if id != "goo" && id != "reyth" {
		w.WriteHeader(403)
		return
	}
	f(h, w, r)
}
