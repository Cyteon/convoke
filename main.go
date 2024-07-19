package main

import (
	"convoke/server/api"
	"convoke/server/api/admin"
	"convoke/server/api/player"
	"convoke/server/ui"
	"convoke/server/ws"
	"convoke/server/ws/room"
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
	router.HandleFunc("/ws/room/{slug}", room.HandleConnection)
	router.HandleFunc("/ws/ping", ws.HandlePing)

	utils.Log("Loading api routes", "")
	router.HandleFunc("/api/ping", api.HandlePing).Methods("GET", "POST")

	router.HandleFunc("/api/player/new", player.HandleNew).Methods("POST")
	router.HandleFunc("/api/player/login", player.HandleLogin).Methods("POST")

	router.HandleFunc("/api/admin/login", admin.HandleLogin).Methods("POST")
	router.HandleFunc("/api/admin/verify", admin.HandleVerify).Methods("POST")
	router.HandleFunc("/api/admin/users", admin.HandleUsers).Methods("GET")

	utils.Log("Loading webui routes", "")
	router.HandleFunc("/ui/login", ui.HandleLogin).Methods("GET")
	router.HandleFunc("/ui/admin", ui.HandleAdmin).Methods("GET")
	router.HandleFunc("/ui/admin/users", ui.HandleAdminUsers).Methods("GET")

	utils.Log("Listening on "+config.Websocket.Host+":"+config.Websocket.Port, "green")
	utils.Log("Admin interface at http://"+config.Websocket.Host+":"+config.Websocket.Port+"/ui/admin", "green")

	log.Fatal(http.ListenAndServe(config.Websocket.Host+":"+config.Websocket.Port, router))
}
