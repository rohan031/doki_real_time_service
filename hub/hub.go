package hub

import (
	"doki.co.in/doki_real_time_service/client"
	"doki.co.in/doki_real_time_service/utils"
	"errors"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gorilla/websocket"
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
	clients      clientList
	jwks         *keyfunc.Keyfunc
	subscription subscription
}

// addClient adds newly connected client to Hub
func (h *Hub) addClient(user string, client client.Client) {
	h.Lock()
	defer h.Unlock()

	username, resource := utils.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return
	}

	if h.clients[username] == nil {
		h.clients[username] = make(resourceList)
	}

	h.clients[username][resource] = client
}

// removeClient closes and removes connection from Hub
func (h *Hub) removeClient(c client.Client) {
	h.Lock()
	defer h.Unlock()

	username, resource := c.GetUserInfo()
	completeUser := utils.CreateUserFromUsernameAndResource(username, resource)
	if username == "" || resource == "" {
		return
	}

	// check the connection we are tyring to remove and the connection that is present are same
	// this can happen if client resource is same but underlying tcp connection is changed
	if conn, ok := h.clients[username][resource]; ok && conn.GetConnection() == c.GetConnection() {
		// close the websocket connection
		if conn.GetConnection() != nil {
			_ = conn.GetConnection().Close()
		}

		mySubscriptions := conn.GetMySubscriptions()
		for subscription := range mySubscriptions {
			h.Unsubscribe(subscription, completeUser)
		}
		// remove resource from username
		delete(h.clients[username], resource)

		// if empty remove the username too
		if len(h.clients[username]) == 0 {
			// send offline status too for this user
			h.sendPresence(false, username)
			delete(h.clients, username)
		}
	}
}

// GetIndividualClient is used to get user connected client in particular resource
// this will be used when server needs to send updates for the post and other user subscriptions
func (h *Hub) GetIndividualClient(user string) client.Client {
	username, resource := utils.GetUsernameAndResourceFromUser(user)
	if username == "" || resource == "" {
		return nil
	}

	if _, ok := h.clients[username]; !ok {
		return nil
	}

	if _, ok := h.clients[username][resource]; !ok {
		return nil
	}

	return h.clients[username][resource]

}

// GetAllConnectedClients will return all the connected clients for a particular user
// this will be used when forwarding user messages
func (h *Hub) GetAllConnectedClients(username string) map[string]client.Client {
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

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resource := r.URL.Query().Get("resource")
	if resource == "" {
		resource = utils.RandomString()
	}

	user := utils.CreateUserFromUsernameAndResource(username, resource)
	newClient := createClient(conn, h, user)

	h.addClient(user, newClient)

	// sending my initial online presence
	h.sendPresence(true, username)

	go newClient.readMessage()
	go newClient.writeMessage()

}

// CreateHub creates a new hub
func CreateHub(jwks *keyfunc.Keyfunc) *Hub {
	return &Hub{
		clients: make(clientList),
		jwks:    jwks,
		subscription: subscription{
			subscriptions: make(nodeSubscription),
		},
	}
}