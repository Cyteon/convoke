package utils

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var session *r.Session

func LoadDB() *r.Session {
	Log("Connection to DB", "")

	config := LoadConfig("config.yaml")

	session, err := r.Connect(r.ConnectOpts{
		Address:  config.DbUrl,
		Database: "convoke",
		Username: config.DbUser,
		Password: config.DbPass,
	})

	if err != nil {
		LogFatal("DB Connection Error, "+err.Error(), "red")
	}

	if session.IsConnected() {
		Log("Connected to DB", "green")
	} else {
		LogFatal("DB Error, not connected", "red")
	}

	return session
}
