package hub

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Hub handles all the client connection and related methods

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
}

// ServeWS methods takes the current [http] request
// and upgrade it to [websocket] connection
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection for websocket")

	_, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading incoming http connection to websocket: %v\n", err)
		return
	}
}

// CreateHub creates a new hub
func CreateHub() *Hub {
	return &Hub{}
}