package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type IdName struct {
	ID   string
	Name string
}

var d struct {
	sync.Mutex
	*sql.DB
}

func Init(path string) {
	d.Lock()
	defer d.Unlock()
	db, err := sql.Open("sqlite3", path)
	d.DB = db
	if err != nil {
		log.Fatal("Can't load database:", err)
	}
}

func KnownToken(token string) bool {
	res := GetUserField(token, "Token")
	return !(len(res) == 0)
}

func GetChannelsList() []IdName {
	d.Lock()
	defer d.Unlock()
	channels := make([]IdName, 0)
	rows, err := d.Query("SELECT ID, Name FROM Channels")
	if err != nil {
		log.Println("Error getting channels list:", err)
		return channels
	}
	for rows.Next() {
		var channel IdName
		err = rows.Scan(&(channel.ID), &(channel.Name))
		if err != nil {
			log.Println("Error parsing row:", err)
			continue
		}
		channels = append(channels, channel)
	}
	return channels
}
