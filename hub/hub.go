package hub

import (
	"doki.co.in/doki_real_time_service/helper"
	"errors"
	"github.com/MicahParks/keyfunc/v3"
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
	clients clientList
	jwks    *keyfunc.Keyfunc
}

// addClient adds newly connected client to Hub
func (h *Hub) addClient(user string, client *client) {
	h.Lock()
	defer h.Unlock()

	username, resource := helper.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return
	}

	if h.clients[username] == nil {
		h.clients[username] = make(resourceList)

	}
	h.clients[username][resource] = client
}

// removeClient closes and removes connection from Hub
func (h *Hub) removeClient(c *client) {
	h.Lock()
	defer h.Unlock()

	username, resource := helper.GetUsernameAndResourceFromUser(c.user)
	if username == "" || resource == "" {
		return
	}

	if conn, ok := h.clients[username][resource]; ok && c.connection == conn.connection {
		// close the websocket connection
		_ = conn.connection.Close()

		// remove resource from username
		delete(h.clients[username], resource)

		// if empty remove the username too
		if len(h.clients[username]) == 0 {
			delete(h.clients, username)
		}
	}
}

// getIndividualClient is used to get user connected client in particular resource
// this will be used when server needs to send updates for the post and other user subscriptions
func (h *Hub) getIndividualClient(user string) *client {
	username, resource := helper.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return nil
	}

	return h.clients[username][resource]
}

// getAllConnectedClients will return all the connected clients for a particular user
// this will be used when forwarding user messages
func (h *Hub) getAllConnectedClients(username string) map[string]*client {
	return h.clients[username]
}

// ServeWS methods takes the current [http] request
// and upgrade it to [websocket] connection
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	username, err := parseAuthHeader(r, h.jwks)
	if err != nil {
		var authErrorObject *authError
		if errors.As(err, &authErrorObject) {
			http.Error(w, authErrorObject.Error(), authErrorObject.Code)
		}
		return
	}

	log.Println("new connection for websocket")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading incoming http connection to websocket: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resource := r.URL.Query().Get("resource")
	if resource == "" {
		resource = helper.RandomString()
	}

	log.Printf("new connection: %v@%v\n\n", username, resource)

	user := helper.CreateUserFromUsernameAndResource(username, resource)
	newClient := createClient(conn, h, user)

	h.addClient(user, newClient)

	go newClient.readMessage()
	go newClient.writeMessage()

}

// CreateHub creates a new hub
func CreateHub(jwks *keyfunc.Keyfunc) *Hub {
	return &Hub{
		clients: make(clientList),
		jwks:    jwks,
	}
}