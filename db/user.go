package db

import (
	"log"
	"strings"
)

func UpdateUser(token, column, value string) error {
	d.Lock()
	defer d.Unlock()
	_, err := d.Exec("UPDATE Users SET "+column+" = $1 WHERE Token = $2", value, token)
	if err != nil {
		log.Println("Can't update user:", err)
	}
	return err
}

func GetUserField(token, column string) string {
	d.Lock()
	defer d.Unlock()
	var res string
	err := d.QueryRow("SELECT "+column+" FROM Users WHERE Token = ?", token).Scan(&res)
	if err != nil {
		log.Println("Can't get row:", err)
	}
	return res
}

func GetUserChannels(token string) []struct {
	ID string
} {
	channels := make([]struct {
		ID string
	}, 0)
	list := GetUserField(token, "Channels")
	for _, c := range strings.Split(list, ",") {
		channels = append(channels, struct {
			ID string
		}{
			c,
		})
	}
	return channels
}

func InsertUser(token, id, name string) error {
	d.Lock()
	defer d.Unlock()
	_, err := d.Exec("INSERT INTO Users(Token, ID, Name, Channels) VALUES(?,?,?,?)", token, id, name, "")
	return err
}
