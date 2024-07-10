package users

import (
	"convoke/utils"
	"encoding/json"
	"net/http"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type User struct {
	Username string
	Password string
	Email    string
	Friends  []string
}

func HandleNew(w http.ResponseWriter, r *http.Request) {
	// Only POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User

	// Decode the JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if user.Username == "" || user.Password == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "username already exists"})
		return
	}

	session := utils.LoadDB()

	// Check if the username or email already exists
	cursor, err := rethink.DB("convoke").Table("users").Filter(rethink.Row.Field("Username").Eq(user.Username)).Run(session)
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

	cursor, err = rethink.DB("convoke").Table("users").Filter(rethink.Row.Field("Email").Eq(user.Email)).Run(session)
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
	hash, err := utils.HashPassword(user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure all the hashing stuff is working
	if !utils.CheckPasswordHash(user.Password, hash) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Hashing error"})
		return
	}

	user.Password = hash

	_, err = rethink.DB("convoke").Table("users").Insert(user).RunWrite(session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Log("Created user: "+user.Username, "")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
