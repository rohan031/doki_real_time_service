package hub

import (
	"doki.co.in/doki_real_time_service/helper"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub handles all the client connection and related methods
type Hub struct {
	sync.RWMutex
	clients ClientList
}

// AddClient adds newly connected client to Hub
func (h *Hub) AddClient(user string, client *Client) {
	h.Lock()
	defer h.Unlock()

	username, resource := helper.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return
	}

	h.clients[username][resource] = client
}

// RemoveClient closes and removes connection from Hub
func (h *Hub) RemoveClient(user string) {
	h.Lock()
	defer h.Unlock()

	username, resource := helper.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return
	}

	if conn, ok := h.clients[username][resource]; ok {
		// close the websocket connection
		_ = conn.Connection.Close()

		// remove resource from username
		delete(h.clients[username], resource)

		// if empty remove the username too
		if len(h.clients[username]) == 0 {
			delete(h.clients, username)
		}
	}
}

// GetIndividualClient is used to get user connected client in particular resource
// this will be used when server needs to send updates for the post and other user subscriptions
func (h *Hub) GetIndividualClient(user string) *Client {
	username, resource := helper.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return nil
	}

	return h.clients[username][resource]
}

// GetAllConnectedClients will return all the connected clients for a particular user
// this will be used when forwarding user messages
func (h *Hub) GetAllConnectedClients(username string) map[string]*Client {
	return h.clients[username]
}

// ServeWS methods takes the current [http] request
// and upgrade it to [websocket] connection
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection for websocket")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading incoming http connection to websocket: %v\n", err)
		return
	}

	resource := helper.RandomString()
	username := "rohan_verma__"

	user := helper.CreateUserFromUsernameAndResource(username, resource)
	newClient := CreateClient(conn, h, user)

	if h.clients[username] == nil {
		h.clients[username] = make(map[string]*Client)

	}
	h.clients[username][resource] = newClient

	go newClient.ReadMessage()
	go newClient.WriteMessage()

}

// CreateHub creates a new hub
func CreateHub() *Hub {
	return &Hub{
		clients: make(ClientList),
	}
}