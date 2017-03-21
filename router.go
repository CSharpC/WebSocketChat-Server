package main

import "log"

const (
	TYPE_ANY     = iota
	TYPE_CLIENT  = iota
	TYPE_CHANNEL = iota
)

type RouterQuery struct {
	ID     string
	Type   int
	result chan Sendable
}

type Router struct {
	sendables   map[string]Sendable
	addSendable chan Sendable
	query       chan RouterQuery
}

func (r *Router) add() {
	for s := range r.addSendable {
		if _, ok := r.sendables[s.ID()]; ok {
			continue
		}
		log.Println("Adding", s.ID())
		r.sendables[s.ID()] = s
		go func() {
			<-s.Disconnect()
			log.Println("Deleting", s.ID(), "from the server")
			delete(r.sendables, s.ID())
			s.Disconnect() <- true
		}()
	}
}

func (r *Router) listenQuery() {
	for q := range r.query {
		if s, ok := r.sendables[q.ID]; ok && (s.Type() == q.Type || q.Type == TYPE_ANY) {
			q.result <- s
		} else {
			q.result <- nil
		}
	}
}

func (r *Router) Run() {
	go r.add()
	go r.listenQuery()
}
