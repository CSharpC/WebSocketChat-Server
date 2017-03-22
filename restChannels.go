package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/guglicap/chatServer/db"
)

func RESTChanList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !db.KnownToken(token) {
		w.WriteHeader(403)
		return
	}
	channels := db.GetChannelsList()
	result, err := json.Marshal(channels)
	if err != nil {
		log.Println("Error JSONing channels:", err)
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}
