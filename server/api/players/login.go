package players

import (
	"convoke/utils"
	"encoding/json"
	"net/http"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Login struct {
	Email    string
	Username string
	Password string
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Only POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var login Login
	var player Player

	// Decode the JSON
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if login.Password == "" || (login.Email == "" && login.Username == "") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	session := utils.LoadDB()

	// Find user by email or username
	if login.Email != "" {
		cursor, err := rethink.DB("convoke").Table("players").Filter(rethink.Row.Field("Email").Eq(login.Email)).Run(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close()

		if cursor.IsNil() == false {
			cursor.One(&player)
		}
	} else {
		cursor, err := rethink.DB("convoke").Table("players").Filter(rethink.Row.Field("Username").Eq(login.Username)).Run(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close()

		if cursor.IsNil() == false {
			cursor.One(&player)
		}
	}

	if utils.CheckPasswordHash(login.Password, player.Password) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})

		return
	}

	utils.Log("Authorized user: "+login.Username+""+login.Email, "")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Authorized", "token": player.Token})
}
