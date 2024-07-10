package users

import (
	"convoke/utils"
	"encoding/json"
	"net/http"
)

type User struct {
	Username string
	Password string
	Email    string
	PassHash string
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

	hash, err := utils.HashPassword(user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.PassHash = hash

	if !utils.CheckPasswordHash(user.Password, user.PassHash) {
		http.Error(w, "Hashing error", http.StatusInternalServerError)
		return
	}

	utils.Log("Creating user: "+user.Username+" with hash: "+user.PassHash, "")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
