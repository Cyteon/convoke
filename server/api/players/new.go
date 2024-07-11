package players

import (
	"convoke/utils"
	"encoding/json"
	"net/http"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Player struct {
	Username string
	Password string
	Email    string
	Friends  []string
	Token    string
}

func HandleNew(w http.ResponseWriter, r *http.Request) {
	// Only POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var player Player

	// Decode the JSON
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if player.Username == "" || player.Password == "" || player.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	session := utils.LoadDB()

	// Check if the username or email already exists
	cursor, err := rethink.DB("convoke").Table("players").Filter(rethink.Row.Field("Username").Eq(player.Username)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	if cursor.IsNil() == false {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "username already exists"})
		return
	}

	cursor, err = rethink.DB("convoke").Table("players").Filter(rethink.Row.Field("Email").Eq(player.Email)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	if cursor.IsNil() == false {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email already exists"})
		return
	}

	// We have to hash the password for safe storing
	hash, err := utils.HashPassword(player.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure all the hashing stuff is working
	if !utils.CheckPasswordHash(player.Password, hash) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Hashing error"})
		return
	}

	player.Password = hash

	var token string

	for {
		token = utils.GenerateSecureToken(32)

		cursor, err = rethink.DB("convoke").Table("players").Filter(rethink.Row.Field("Token").Eq(token)).Run(session)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if cursor.IsNil() {
			break
		}
	}

	player.Token = token

	_, err = rethink.DB("convoke").Table("players").Insert(player).RunWrite(session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Log("Created user: "+player.Username, "")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully", "token": token})
}
