package main

import (
	"time"
)

type Userdata struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	Power        int       `json:"power"`
	LastActiveAt time.Time `json:"lastActiveAt"`
	CreatedAt    time.Time `json:"createdAt"`
	RegisteredAt time.Time `json:"registeredAt"`
}

func usernameInUse(name string) bool {
	var username string
	err := db.QueryRow("SELECT name FROM users WHERE name = ?", name).Scan(&username)
	if err == nil {
		if username == name {
			return true
		}
	}
	return false
}

func useridInUse(userid string) bool {
	var id string
	err := db.QueryRow("SELECT id FROM users WHERE id = ?", userid).Scan(&id)
	if err == nil {
		if id == userid {
			return true
		}
	}
	return false
}

func validateUserCredentials(userid string, username string, password string) bool {
	var id string
	var name string
	var pass string
	err := db.QueryRow("SELECT id, name, password FROM users WHERE id = ? AND name = ? AND password = ?", userid, username, password).Scan(&id, &name, &pass)
	if err == nil {
		if id == userid {
			return true
		}
	}
	return false
}

func selectUser(id string, public bool) Userdata {
	var data Userdata
	if public {
		db.QueryRow("SELECT id, name, power, lastActiveAt FROM users WHERE id = ?", id).Scan(&data.Id, &data.Name, &data.Power, &data.LastActiveAt)
	} else {
		db.QueryRow("SELECT id, name, password, power, lastActiveAt, createdAt FROM users WHERE id = ?", id).Scan(&data.Id, &data.Name, &data.Password, &data.Power, &data.LastActiveAt, &data.CreatedAt)
	}
	return data
}

func insertUser(id string, name string, password string, power int, lastActiveAt time.Time, createdAt time.Time) {
	_, err := db.Exec("INSERT INTO users (id, name, password, power, lastActiveAt, createdAt) VALUES (?, ?, ?, ?, ?, ?)", id, name, password, power, lastActiveAt, createdAt)
	if err == nil {
		println("user inserted")
	} else {
		println("user insertion failed")
	}
}
