package main

import (
	"convoke/server"
	"convoke/utils"
	"log"
	"net/http"
)

func main() {
	utils.Log("Starting convoke", "cyan")

	config := utils.LoadConfig("config.yaml")

	session := utils.LoadDB()
	defer session.Close()

	utils.Log("Loading routes", "")
	http.HandleFunc("/ping", server.HandlePing)

	utils.Log("Listening on "+config.Websocket.Host+":"+config.Websocket.Port, "green")

	log.Fatal(http.ListenAndServe(config.Websocket.Host+":"+config.Websocket.Port, nil))
}
