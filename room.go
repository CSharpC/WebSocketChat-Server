package main

type Room struct {
	clients map[*Client]bool

	inMsg       chan Message
	subscribe   chan *Client
	unsubscribe chan *Client
}

func newRoom() *Room {
	return &Room{
		clients:     make(map[*Client]bool),
		subscribe:   make(chan *Client),
		unsubscribe: make(chan *Client),
		inMsg:       make(chan Message, 20),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.subscribe:
			r.clients[client] = true
		case client := <-r.unsubscribe:
			r.clients[client] = false
		case msg := <-r.inMsg:
			for c, s := range r.clients {
				if s {
					c.send <- msg
				}
			}
		}
	}
}
