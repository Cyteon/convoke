package utils

import (
	"strings"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var Session *r.Session

func LoadDB() *r.Session {
	config := LoadConfig("config.yaml")

	Session, err := r.Connect(r.ConnectOpts{
		Address:  config.DbUrl,
		Database: "convoke",
		Username: config.DbUser,
		Password: config.DbPass,
	})

	if err != nil {
		LogFatal("DB Connection Error, "+err.Error(), "red")
	}

	return Session
}

func SetupDB() *r.Session {
	Log("Setting up DB", "")

	Log("Connecting to DB", "")
	Session := LoadDB()

	_, err := r.DBCreate("convoke").RunWrite(Session)
	if err != nil {
		// Check if the error is because the database already exists
		if !strings.Contains(err.Error(), "Database `convoke` already exists") {
			LogFatal("Error setting up DB: "+err.Error(), "red")
		} else {
			Log("Database 'convoke' already exists, continuing setup", "yellow")
		}
	} else {
		Log("Database 'convoke' created successfully", "green")
	}

	Session.Use("convoke")

	_, err = r.TableCreate("players").RunWrite(Session)
	if err != nil {
		if !strings.Contains(err.Error(), "Table `convoke.players` already exists") {
			LogFatal("Error setting up DB: "+err.Error(), "red")
		} else {
			Log("Table 'convoke.players' already exists, continuing setup", "yellow")
		}
	} else {
		Log("Table 'players' created successfully", "green")
	}

	_, err = r.TableCreate("admins").RunWrite(Session)
	if err != nil {
		if !strings.Contains(err.Error(), "Table `convoke.admins` already exists") {
			LogFatal("Error setting up DB: "+err.Error(), "red")
		} else {
			Log("Table 'convoke.admins' already exists, continuing setup", "yellow")
		}
	} else {
		Log("Table 'admins' created successfully", "green")

		Log("Creating default admin user", "green")

		password, err := HashPassword("admin")

		if err != nil {
			LogFatal("Error hashing password: "+err.Error(), "red")
		}

		_, err = r.Table("admins").Insert(map[string]string{
			"Username": "admin",
			"Password": password,
		}).RunWrite(Session)

		if err != nil {
			LogFatal("Error creating default admin user: "+err.Error(), "red")
		}

		Log("Default admin user created successfully", "green")

		Log("WARNING: Default admin user created with username 'admin' and password 'admin'", "red")
		Log("Please change the password immediately", "red")

	}

	Log("Database and table setup complete", "green")

	if Session.IsConnected() {
		Log("Connected to DB", "green")
	} else {
		LogFatal("DB Error, not connected", "red")
	}

	return Session
}
