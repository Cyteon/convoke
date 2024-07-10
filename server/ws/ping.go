package ws

import (
	"net/http"

	"convoke/server"
	"convoke/utils"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	conn, err := server.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		utils.Log("Error upgrading connection, "+err.Error(), "red")
		return
	}
	defer conn.Close()

	utils.Log("Connection opened with "+r.RemoteAddr, "")

	for {
		// We just return what the client sent back to them

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
