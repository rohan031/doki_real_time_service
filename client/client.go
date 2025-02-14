package client

import "github.com/gorilla/websocket"

// Client is the global interface used by payload and hub to store and get clients
type Client interface {
	GetConnection() *websocket.Conn
	GetUserInfo() (string, string)
	WriteToChannel(*[]byte)
	GetMySubscriptions() map[string]bool
	AddSubscription(string)
}