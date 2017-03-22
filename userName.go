package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/guglicap/chatServer/db"
)

func Name(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PATCH":
		updateUserName(w, r)
	case "GET":
		getUserName(w, r)
	}
}

func updateUserName(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !db.KnownToken(token) {
		w.WriteHeader(403)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	var post struct {
		Name string
	}
	err = json.Unmarshal(body, &post)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	err = db.UpdateUser(token, "Name", post.Name)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func getUserName(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !db.KnownToken(token) {
		w.WriteHeader(403)
		return
	}
	name := db.GetUserField(token, "Name")
	if len(name) < 1 {
		w.WriteHeader(500)
		return
	}
	res, err := json.Marshal(struct {
		Name string
	}{
		name,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(res)
}
