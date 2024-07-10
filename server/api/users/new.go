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
}

func HandleNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	session := utils.LoadDB()

	cursor, err := rethink.DB("convoke").Table("users").Filter(rethink.Row.Field("Username").Eq(user.Username)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	if cursor.IsNil() == false {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	cursor, err = rethink.DB("convoke").Table("users").Filter(rethink.Row.Field("Email").Eq(user.Email)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	if cursor.IsNil() == false {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	hash, err := utils.HashPassword(user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !utils.CheckPasswordHash(user.Password, hash) {
		http.Error(w, "Hashing error", http.StatusInternalServerError)
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
