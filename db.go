package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type dbchannel struct {
	ID   string
	Name string
}

var database struct {
	sync.Mutex
	db *sql.DB
}

func loadDB(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("Can't load database:", err)
	}
	database.Lock()
	database.db = db
	database.Unlock()
}

func knownID(id string) bool {
	database.Lock()
	defer database.Unlock()
	db := database.db
	var res string
	err := db.QueryRow("SELECT ID FROM Users WHERE ID = ?", id).Scan(&res)
	if err != nil {
		return false
	}
	return true
}

func getChannelsList() []dbchannel {
	database.Lock()
	defer database.Unlock()
	channels := make([]dbchannel, 0)
	db := database.db
	rows, err := db.Query("SELECT ID, Name FROM Channels")
	if err != nil {
		log.Println("Error getting channels list:", err)
		return channels
	}
	for rows.Next() {
		var channel dbchannel
		err = rows.Scan(&(channel.ID), &(channel.Name))
		if err != nil {
			log.Println("Error parsing row:", err)
			continue
		}
		channels = append(channels, channel)
	}
	return channels
}
