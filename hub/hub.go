package hub

import (
	"doki.co.in/doki_real_time_service/client"
	"doki.co.in/doki_real_time_service/payload"
	"doki.co.in/doki_real_time_service/utils"
	"errors"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"slices"
	"sync"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub handles all the client connection and related methods
type Hub struct {
	sync.RWMutex
	clients              clientList
	jwks                 *keyfunc.Keyfunc
	presenceSubscription presenceSubscription
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
	if username == "" || resource == "" {
		return
	}

	// check the connection we are tyring to remove and the connection that is present are same
	// this can happen if client resource is same but underlying tcp connection is changed
	if conn, ok := h.clients[username][resource]; ok && c.GetConnection() == conn.GetConnection() {
		// close the websocket connection
		if conn.GetConnection() != nil {
			_ = conn.GetConnection().Close()
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

// sendPresence sends user presence updates to all the subscribed users
func (h *Hub) sendPresence(online bool, username string) {
	// find user in subscription and send status change
	usersSubscribed, ok := h.presenceSubscription[username]
	if !ok {
		return
	}

	for _, completeUser := range usersSubscribed {
		conn := h.GetIndividualClient(completeUser)
		if conn == nil {
			h.UnsubscribeUserPresence(username, completeUser)
			continue
		}

		user, resource := utils.GetUsernameAndResourceFromUser(completeUser)
		presencePayload := payload.CreatePresencePayload(username, user, online)

		log.Printf("\nSending user presence: %v, %v\n", username, completeUser)
		data := utils.PayloadToJson(presencePayload)
		if data != nil {
			presencePayload.SendPayload(data, h, resource)
		}
	}
}

// sendInitialPresence sends the given user presence on initial subscription
func (h *Hub) sendInitialPresence(userPresence string, completeUser string) {
	_, ok := h.clients[userPresence]
	conn := h.GetIndividualClient(completeUser)

	if conn == nil {
		h.UnsubscribeUserPresence(userPresence, completeUser)
		return
	}

	log.Printf("\nSending initial presence: %v, %v\n", userPresence, completeUser)
	username, resource := utils.GetUsernameAndResourceFromUser(completeUser)
	presencePayload := payload.CreatePresencePayload(userPresence, username, ok)

	data := utils.PayloadToJson(presencePayload)
	if data != nil {
		presencePayload.SendPayload(data, h, resource)
	}
}

// SubscribeUserPresence subscribes the given complete user to the user presence updates
func (h *Hub) SubscribeUserPresence(userToSubscribe string, completeUser string) {
	h.Lock()
	defer h.Unlock()
	log.Printf("\nSubscribing to user presence: %v, %v\n", userToSubscribe, completeUser)
	h.presenceSubscription[userToSubscribe] = append(h.presenceSubscription[userToSubscribe], completeUser)
	h.sendInitialPresence(userToSubscribe, completeUser)
}

// UnsubscribeUserPresence unsubscribes the given complete user to the user presence updates
func (h *Hub) UnsubscribeUserPresence(userToUnsubscribe string, completeUser string) {
	h.Lock()
	defer h.Unlock()

	log.Printf("\nUn-Subscribing to user presence: %v, %v\n", userToUnsubscribe, completeUser)
	// find slice
	_, ok := h.presenceSubscription[userToUnsubscribe]
	if !ok {
		return
	}

	h.presenceSubscription[userToUnsubscribe] = slices.DeleteFunc(h.presenceSubscription[userToUnsubscribe], func(user string) bool {
		return completeUser == user
	})

	if len(h.presenceSubscription[userToUnsubscribe]) < 0 {
		delete(h.presenceSubscription, userToUnsubscribe)
	}
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
		resource = utils.RandomString()
	}

	log.Printf("new connection: %v@%v\n\n", username, resource)

	user := utils.CreateUserFromUsernameAndResource(username, resource)
	newClient := createClient(conn, h, user)

	h.addClient(user, newClient)
	h.sendPresence(true, username)

	go newClient.readMessage()
	go newClient.writeMessage()

}

// CreateHub creates a new hub
func CreateHub(jwks *keyfunc.Keyfunc) *Hub {
	return &Hub{
		clients:              make(clientList),
		jwks:                 jwks,
		presenceSubscription: make(presenceSubscription),
	}
}