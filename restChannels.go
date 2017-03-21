package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func channelList(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("client_id")
	if !knownID(id) {
		w.WriteHeader(403)
		return
	}
	channels := getChannelsList()
	result, err := json.Marshal(channels)
	if err != nil {
		log.Println("Error JSONing channels:", err)
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}
