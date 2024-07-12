package admin

import (
	"convoke/utils"
	"encoding/json"
	"net/http"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Login struct {
	Username string
	Password string
}

type Admin struct {
	ID       string
	Username string
	Password string
	Token    string
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var login Login
	// Decode the JSON
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check if the request is valid
	if login.Password == "" || login.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	session := utils.LoadDB()

	// Find user by username
	cursor, err := rethink.DB("convoke").Table("admins").Filter(rethink.Row.Field("Username").Eq(login.Username)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	var admin Admin

	if cursor.IsNil() == false {
		cursor.One(&admin)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not found"})
		return
	}

	// Check the password
	if utils.CheckPasswordHash(login.Password, admin.Password) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})

		return
	}

	token := utils.GenerateSecureToken(32)

	// Update the token so you can only be logged in on one device
	_, err = rethink.DB("convoke").Table("admins").Get(admin.ID).Update(map[string]string{"Token": token}).RunWrite(session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Log("Admin: "+admin.Username+" logged in", "")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Authorized", "token": token})
}
