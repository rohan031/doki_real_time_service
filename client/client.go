package client

import (
	"doki.co.in/doki_real_time_service/hub"
	"github.com/gorilla/websocket"
)

// ClientList contains all the connection that are currently
// connected to the server.
// each user has its own map of connected clients
// at a time same user with multiple device can connect
type ClientList map[string]map[string]*Client

type Client struct {
	connection *websocket.Conn
	hub        *hub.Hub

	// channel buffering to prevent writing to connection concurrently
	write chan []byte
}

func CreateClient(conn *websocket.Conn, hub *hub.Hub) *Client {
	return &Client{
		connection: conn,
		hub:        hub,
		write:      make(chan []byte),
	}
}