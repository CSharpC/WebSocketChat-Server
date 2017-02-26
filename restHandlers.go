package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ChannelJSON struct {
	Name string
	ID   string
}

type UserJSON struct {
	Name string
	ID   string
}

type RESTHandlerFunc func(h *Hub, w http.ResponseWriter, r *http.Request)

func handleChannelsList(h *Hub, w http.ResponseWriter, r *http.Request) {
	channels := make([]ChannelJSON, 0)
	for k, v := range h.rooms {
		channel := ChannelJSON{
			v.Name,
			k,
		}
		channels = append(channels, channel)
	}
	result, err := json.Marshal(channels)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", string(result))
}

func handleChannelCreate(h *Hub, w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ChannelCreate: error reading body")
		return
	}
	var channel ChannelJSON
	err = json.Unmarshal(reqBody, &channel)
	if err != nil || len(channel.Name) < 1 {
		log.Println("ChannelCreate: error parsing json")
		return
	}
	h.addChannel <- newChannel(getUUID(), channel.Name)
	return
}

func handleChannelDelete(h *Hub, w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ChannelCreate: error reading body")
		return
	}
	var channel ChannelJSON
	err = json.Unmarshal(reqBody, &channel)
	if err != nil || len(channel.ID) < 1 {
		log.Println("ChannelCreate: error parsing json")
		return
	}
	h.removeRoom <- h.rooms[channel.ID]
	return
}

func handleUserName(h *Hub, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("client_id")
	var client *Client
	if _, ok := h.clients[id]; ok {
		client = h.clients[id]
	} else {
		log.Println("UserName: user not found.")
		return
	}
	var u UserJSON
	if r.Method == "GET" {
		u.Name = client.nickname
		resp, err := json.Marshal(&u)
		if err != nil {
			log.Println("UserName: error Marshaling")
			return
		}
		w.Write(resp)
		return
	}
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("UserName: error reading request")
		return
	}
	err = json.Unmarshal(req, &u)
	if err != nil || len(u.Name) < 1 {
		log.Println("UserName: error decoding json")
		return
	}
	h.clients[id].nickname = u.Name
	return
}
