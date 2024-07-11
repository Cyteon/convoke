package admins

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

	if login.Password == "" || login.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	session := utils.LoadDB()

	// Find user by email or username
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

	if utils.CheckPasswordHash(login.Password, admin.Password) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})

		return
	}

	admin.Token = utils.GenerateSecureToken(32)

	_, err = rethink.DB("convoke").Table("admins").Get(admin.Username).Update(map[string]string{"Token": admin.Token}).RunWrite(session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Log("Authorized admin: "+admin.Username, "")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Authorized", "token": admin.Token})
}
