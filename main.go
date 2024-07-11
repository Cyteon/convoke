package main

import (
	"convoke/server/api"
	"convoke/server/api/admins"
	"convoke/server/api/players"
	"convoke/server/ui"
	"convoke/server/ws"
	"convoke/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	utils.Log("Starting convoke", "cyan")

	config := utils.LoadConfig("config.yaml")

	session := utils.SetupDB()
	defer session.Close()

	router := mux.NewRouter()

	utils.Log("Loading websocket routes", "")
	router.HandleFunc("/ws/{slug}", ws.HandleConnection)
	router.HandleFunc("/ws/ping", ws.HandlePing)

	utils.Log("Loading api routes", "")
	router.HandleFunc("/api/ping", api.HandlePing).Methods("GET", "POST")

	router.HandleFunc("/api/players/new", players.HandleNew).Methods("POST")
	router.HandleFunc("/api/players/login", players.HandleLogin).Methods("POST")

	router.HandleFunc("/api/admins/login", admins.HandleLogin).Methods("POST")

	utils.Log("Loading webui routes", "")
	router.HandleFunc("/ui/login", ui.HandleLogin).Methods("GET")
	router.HandleFunc("/ui/admin", ui.HandleAdmin).Methods("GET")

	utils.Log("Listening on "+config.Websocket.Host+":"+config.Websocket.Port, "green")

	log.Fatal(http.ListenAndServe(config.Websocket.Host+":"+config.Websocket.Port, router))
}
