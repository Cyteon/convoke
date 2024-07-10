package api

import (
	"convoke/utils"
	"net/http"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	utils.Log("Ping", "")
	w.Write([]byte("OK"))
}
