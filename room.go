package main

type Channel struct {
	ID          string
	Name        string
	clients     map[*Client]bool
	send        chan Message
	subscribe   chan *Client
	unsubscribe chan *Client
}

func newChannel(ID, name string) *Channel {
	return &Channel{
		clients:     make(map[*Client]bool),
		subscribe:   make(chan *Client),
		unsubscribe: make(chan *Client),
		ID:          ID,
		Name:        name,
		send:        make(chan Message, 20),
	}
}

func (r *Channel) run() {
	for {
		select {
		case client := <-r.subscribe:
			r.clients[client] = true
		case client := <-r.unsubscribe:
			r.clients[client] = false
			delete(r.clients, client)
		case msg := <-r.send:
			for c, s := range r.clients {
				if s {
					c.send <- msg
				}
			}
		}
	}
}
