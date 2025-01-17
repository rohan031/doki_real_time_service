package client

import (
	"doki.co.in/doki_real_time_service/hub"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	pongWait             = 60 * time.Second
	pingInterval         = (pongWait * 9) / 10
	incomingPayloadLimit = int64(512)
)

// ClientList contains all the connection that are currently
// connected to the server.
//
// each user has its own map of connected clients
// at a time same user with multiple device can connect
type ClientList map[string]map[string]*Client

type Client struct {
	connection *websocket.Conn
	hub        *hub.Hub

	// user is complete user with resource part
	// e.g. username: rohan_verma__, is connected through [doki] native client
	// than user will be: rohan_verma__@{resource} where resource is unique string
	// to identify the particular client
	user string

	// channel buffering to prevent writing to connection concurrently
	write chan []byte
}

// ReadMessage reads all the incoming messages from the connection
func (c *Client) ReadMessage() {
	defer func() {
		c.hub.RemoveClient(c.user)
	}()

	// adding max payload any client can send through [connection]
	c.connection.SetReadLimit(incomingPayloadLimit)

	// set pong wait
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting pongwait: %v\n", err)
		return
	}
	c.connection.SetPongHandler(func(string) error {
		return c.connection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}

		// add this message to queue to be handled my message archive write service
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
	}
}

func (c *Client) WriteMessage() {
	// ping ticker
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.hub.RemoveClient(c.user)
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

func CreateClient(conn *websocket.Conn, hub *hub.Hub) *Client {
	return &Client{
		connection: conn,
		hub:        hub,
		write:      make(chan []byte),
	}
}
