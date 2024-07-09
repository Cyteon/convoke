package server

import (
	"net/http"

	"convoke/utils"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	conn, err := GetUpgrader().Upgrade(w, r, nil)

	if err != nil {
		utils.Log("Error upgrading connection, "+err.Error(), "red")
		return
	}
	defer conn.Close()

	utils.Log("Connection opened with "+r.RemoteAddr, "")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			utils.Log("Error reading message, "+err.Error(), "red")
			break
		}

		err = conn.WriteMessage(messageType, p)

		if err != nil {
			utils.Log("Error writing message, "+err.Error(), "red")
			break
		}

		utils.Log("Recived: "+string(p), "")
	}
}
