package admin

import (
	"convoke/utils"
	"encoding/json"
	"net/http"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type TokenLogin struct {
	Token string
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	var login TokenLogin
	// Decode the JSON
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check if the request is valid
	if login.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	session := utils.LoadDB()

	// Find user by username
	cursor, err := rethink.DB("convoke").Table("admins").Filter(rethink.Row.Field("Token").Eq(login.Token)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	var admin Admin

	if cursor.IsNil() == false {
		cursor.One(&admin)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	// This will print every time an admin goes to any page or refresh, so lets not have this spam logs
	//utils.Log("Authorized admin: "+admin.Username, "")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Authorized", "username": admin.Username})
}
