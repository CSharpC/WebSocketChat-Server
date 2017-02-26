package main

type ClientInfo struct {
	Nickname string
	ID       string
}

type Message struct {
	From    string
	To      string
	Content string
}
