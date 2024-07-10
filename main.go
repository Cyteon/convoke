package main

import (
	"convoke/server/api"
	"convoke/server/ws"
	"convoke/utils"
	"log"
	"net/http"
)

func main() {
	utils.Log("Starting convoke", "cyan")

	config := utils.LoadConfig("config.yaml")

	session := utils.LoadDB()
	defer session.Close()

	utils.Log("Loading websocket routes", "")
	http.HandleFunc("/ws/ping", ws.HandlePing)

	utils.Log("Loading api routes", "")
	http.HandleFunc("/api/ping", api.HandlePing)

	utils.Log("Listening on "+config.Websocket.Host+":"+config.Websocket.Port, "green")

	log.Fatal(http.ListenAndServe(config.Websocket.Host+":"+config.Websocket.Port, nil))
}
