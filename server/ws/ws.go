package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"convoke/server"
	"convoke/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type AuthMessage struct {
	Token string `json:"token"`
}

type Player struct {
	Username string
	Password string
	Email    string
	Friends  []string
	Token    string
}

type Room struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mutex     sync.Mutex
}

type Message struct {
	Content []byte
	Sender  string
}

var rooms = make(map[string]*Room)
var roomsMutex sync.Mutex

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	conn, err := server.Upgrader.Upgrade(w, r, nil)

	authorized := false

	if err != nil {
		utils.Log("Error upgrading connection, "+err.Error(), "red")
		return
	}
	defer conn.Close()

	utils.Log("Connection opened with "+r.RemoteAddr+" to room "+slug+", awaiting authorization", "")

	var room *Room

	for {
		if authorized == false {
			_, p, err := conn.ReadMessage()
			if err != nil {
				utils.Log("Error reading message, "+err.Error(), "red")
				break
			}

			var authMessage AuthMessage

			err = json.Unmarshal(p, &authMessage)
			if err != nil {
				utils.Log("Error unmarshaling message, "+err.Error(), "red")
				break
			}

			if authMessage.Token == "" {
				utils.Log("No token provided", "red")
				break
			}

			// Check if the token is valid

			session := utils.LoadDB()

			cursor, err := rethink.DB("convoke").Table("players").Filter(rethink.Row.Field("Token").Eq(authMessage.Token)).Run(session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer cursor.Close()

			var player Player

			if cursor.IsNil() == false {
				cursor.One(&player)

				authorized = true

				room = getOrCreateRoom(slug)
				room.addClient(conn)
				defer room.removeClient(conn)

				utils.Log("Authorized "+player.Username+" to room "+slug, "")
			}
		} else {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			sender := conn.RemoteAddr().String()
			utils.Log("Received message from "+sender+": "+string(message), "")

			room.broadcast <- Message{Content: message, Sender: sender}
		}
	}
}

func getOrCreateRoom(slug string) *Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if room, exists := rooms[slug]; exists {
		return room
	}

	room := &Room{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
	rooms[slug] = room

	go room.run()

	return room
}

func (r *Room) run() {
	for message := range r.broadcast {
		r.mutex.Lock()
		for client := range r.clients {
			if client.RemoteAddr().String() == message.Sender {
				continue
			}

			err := client.WriteMessage(websocket.TextMessage, message.Content)
			if err != nil {
				utils.Log("Error writing message: "+err.Error(), "yellow")
				client.Close()
				delete(r.clients, client)
			}
		}
		r.mutex.Unlock()
	}
}

func (r *Room) addClient(client *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.clients[client] = true
}

func (r *Room) removeClient(client *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.clients, client)
}
