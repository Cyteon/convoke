package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetUpgrader() *websocket.Upgrader {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	return &upgrader
}
