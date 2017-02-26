package main

type ClientInfo struct {
	Nickname string
	ID       string
}

type Message struct {
	To      string
	From    string
	Content string
}
