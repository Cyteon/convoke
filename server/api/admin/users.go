package admin

import (
	"convoke/server/api/player"
	"convoke/utils"
	"encoding/json"
	"net/http"
	"strings"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	var token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	// Check if the request is valid
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	session := utils.LoadDB()

	// Find user by token
	cursor, err := rethink.DB("convoke").Table("admins").Filter(rethink.Row.Field("Token").Eq(token)).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	var admin Admin
	if cursor.IsNil() {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized", "status": "401"})
		return
	}
	err = cursor.One(&admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cursor.Close()

	cursor, err = rethink.DB("convoke").Table("players").Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close()

	var users []player.Player
	err = cursor.All(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var usersJSON []map[string]interface{}
	for _, user := range users {
		usersJSON = append(usersJSON, map[string]interface{}{
			"ID":       user.ID,
			"Username": user.Username,
			"Email":    user.Email,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"users": usersJSON})
}
