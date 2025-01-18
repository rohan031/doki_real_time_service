package hub

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	pongWait             = 60 * time.Second
	pingInterval         = (pongWait * 9) / 10
	incomingPayloadLimit = int64(512)
)

type resourceList map[string]*client

// clientList contains all the connection that are currently
// connected to the server.
//
// each user has its own map of connected clients
// at a time same user with multiple device can connect
type clientList map[string]resourceList

type client struct {
	connection *websocket.Conn
	hub        *Hub

	// user is complete user with resource part
	// e.g. username: rohan_verma__, is connected through [doki] native client
	// than user will be: rohan_verma__@{resource} where resource is unique string
	// to identify the particular client
	user string

	// channel buffering to prevent writing to connection concurrently
	write chan []byte
}

// readMessage reads all the incoming messages from the connection
func (c *client) readMessage(resource string) {
	defer func() {
		c.hub.removeClient(c.user)
	}()

	// adding max payload any client can send through [connection]
	c.connection.SetReadLimit(incomingPayloadLimit)

	// set pong wait
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting pongwait: %v\n", err)
		return
	}
	c.connection.SetPongHandler(func(string) error {
		log.Printf("received pong from client: %v\n", c.user)
		return c.connection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}

		// parse incoming payload
		var payloadType basePayload
		if !unmarshalAndValidate(&payload, &payloadType) {
			// send error to user
			continue
		}

		switch payloadType.Type {
		case chatMessageType:
			var message chatMessage
			if unmarshalAndValidate(&payload, &message) {
				handleChatMessagePayload(c.hub, &message, &payload, resource)
			}
		case groupChatMessageType:
			var message groupChatMessage
			if unmarshalAndValidate(&payload, &message) {
			}
		case typingStatusType:
			var status typingStatus
			if unmarshalAndValidate(&payload, &status) {
				handleTypingStatusPayload(c.hub, &status, &payload)
			}
		case editMessageType:
			var message editMessage
			if unmarshalAndValidate(&payload, &message) {
			}
		case deleteMessageType:
			var message deleteMessage
			if unmarshalAndValidate(&payload, &message) {
			}
		default:
			// unknown payload type
			// send it to user to tell its unknown or something
			log.Println("unknown payload type")
		}

	}
}

func (c *client) writeMessage() {
	// ping ticker
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.hub.removeClient(c.user)
	}()

	for {
		select {
		case message, ok := <-c.write:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("error closing connection: %v\n", err)
				}
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error sending message: %v\n", err)
			} else {
				log.Printf("message send to client: %v\n\n", c.user)
			}

		case <-ticker.C:
			log.Printf("sending ping to client: %v\n", c.user)
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("error sending ping to client: %v\n", err)
				return
			}
		}
	}
}

func createClient(conn *websocket.Conn, hub *Hub, user string) *client {
	return &client{
		connection: conn,
		hub:        hub,
		write:      make(chan []byte),
		user:       user,
	}
}