package main

type Hub struct {
	rooms        map[string]*Channel
	clients      map[string]*Client
	addChannel   chan *Channel
	removeRoom   chan *Channel
	addClient    chan *Client
	removeClient chan *Client
}

func newHub() *Hub {
	return &Hub{
		make(map[string]*Channel),
		make(map[string]*Client),
		make(chan *Channel),
		make(chan *Channel),
		make(chan *Client),
		make(chan *Client),
	}
}

func (h *Hub) routeMessage(id string) chan Message {
	if room, ok := h.rooms[id]; ok {
		return room.send
	}
	if client, ok := h.clients[id]; ok {
		return client.send
	}
	return nil
}

func (h *Hub) run() {
	for {
		select {
		case room := <-h.addChannel:
			h.rooms[room.ID] = room

		case room := <-h.removeRoom:
			delete(h.rooms, room.ID)

		case client := <-h.addClient:
			h.clients[client.ID] = client

		case client := <-h.removeClient:
			delete(h.clients, client.ID)
			for _, v := range h.rooms {
				if _, ok := v.clients[client]; ok {
					v.unsubscribe <- client
				}
			}
		}
	}
}

func (h *Hub) roomSubscribe(id string) chan *Client {
	if r, ok := h.rooms[id]; ok {
		return r.subscribe
	}
	return nil
}

/*func (h *Hub) broadcast() {
	for _, c := range h.clients {

	}
}*/
