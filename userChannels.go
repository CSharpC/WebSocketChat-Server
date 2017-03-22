package main

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/guglicap/chatServer/db"
)

func userChannels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUserChannels(w, r)
	case "POST", "PATCH":
		postUserChannels(w, r)
	}
}

func postUserChannels(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	id := db.GetUserField(token, "ID")
	if len(id) == 0 {
		w.WriteHeader(403)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var chanID struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &chanID)
	chans := db.GetUserField(token, "Channels")
	q := make(chan Sendable)
	router.query <- RouterQuery{chanID.ID, TypeChannel, q}
	ch := <-q
	if ch == nil {
		w.WriteHeader(404)
		return
	}
	router.query <- RouterQuery{id, TypeClient, q}
	cl := <-q
	if cl != nil {
		ch.(*Channel).Subscribe() <- cl
	}
	err = db.UpdateUser(token, "Channels", chans+","+ch.ID())
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func getUserChannels(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !db.KnownToken(token) {
		w.WriteHeader(403)
		return
	}
	chans := db.GetUserChannels(token)
	resp, err := json.Marshal(chans)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(resp)
}
