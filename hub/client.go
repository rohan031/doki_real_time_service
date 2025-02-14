package hub

import (
	"doki.co.in/doki_real_time_service/client"
	"doki.co.in/doki_real_time_service/payload"
	"doki.co.in/doki_real_time_service/utils"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	pongWait             = 30 * time.Second
	pingInterval         = (pongWait * 9) / 10
	incomingPayloadLimit = int64(16384)
)

type rawClient interface {
	client.Client
	readMessage()
	writeMessage()
}

type resourceList map[string]client.Client

// clientList contains all the connection that are currently
// connected to the server.
//
// each user has its own map of connected Clients
// at a time same user with multiple device can connect
type clientList map[string]resourceList

type clientImpl struct {
	connection *websocket.Conn
	hub        *Hub

	// user is complete user with resource part
	// e.g. username: rohan_verma__, is connected through [doki] native Client
	// than user will be: rohan_verma__@{resource} where resource is unique string
	// to identify the particular Client
	user string

	// channel buffering to prevent writing to connection concurrently
	write chan []byte
}

func (c *clientImpl) GetConnection() *websocket.Conn {
	return c.connection
}

func (c *clientImpl) GetUserInfo() (string, string) {
	return utils.GetUsernameAndResourceFromUser(c.user)
}

func (c *clientImpl) WriteToChannel(data *[]byte) {
	c.write <- *data
}

// readMessage reads all the incoming messages from the connection
func (c *clientImpl) readMessage() {
	defer func() {
		c.hub.removeClient(c)
	}()

	username, resource := c.GetUserInfo()

	// adding max payload any Client can send through [connection]
	c.connection.SetReadLimit(incomingPayloadLimit)

	// set pong wait
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting pongwait: %v\n", err)
		return
	}
	c.connection.SetPongHandler(func(string) error {
		log.Printf("received pong from Client: %v\n", c.user)
		return c.connection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, data, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			log.Printf("error reading message: %v\n", err)
			return
		}

		incomingPayload, err := payload.CreatePayload(&data, username)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// sending payload to relevant recipients
		incomingPayload.SendPayload(&data, c.hub, resource)
	}
}

func (c *clientImpl) writeMessage() {
	// ping ticker
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.hub.removeClient(c)
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
				log.Printf("message send to Client: %v\n\n", c.user)
			}

		case <-ticker.C:
			log.Printf("sending ping to Client: %v\n", c.user)
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("error sending ping to Client: %v\n", err)
				return
			}
		}
	}
}

func createClient(conn *websocket.Conn, hub *Hub, user string) rawClient {
	return &clientImpl{
		connection: conn,
		hub:        hub,
		write:      make(chan []byte),
		user:       user,
	}
}