package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/guglicap/chatServer/db"
	uuid "github.com/nu7hatch/gouuid"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println("Error generating id:", err)
		w.WriteHeader(500)
		return
	}
	token, err := uuid.NewV4()
	if err != nil {
		log.Println("Error generating token:", err)
		w.WriteHeader(500)
		return
	}
	var data struct {
		Name string
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	err = db.InsertUser(token.String(), id.String(), data.Name)
	if err != nil {
		log.Println("Can't insert into DB:", err)
		w.WriteHeader(500)
		return
	}
	resp, err := json.Marshal(struct {
		ID    string `json:"client_id"`
		Token string `json:"access_token"`
	}{
		id.String(),
		token.String(),
	})
	if err != nil {
		log.Println("Error marshaling UUID:", err)
		w.WriteHeader(500)
		return
	}
	w.Write(resp)
}
