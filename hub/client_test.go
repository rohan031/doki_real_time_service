package hub

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConnection(t *testing.T) {
	mockConn := &websocket.Conn{}
	client := &clientImpl{connection: mockConn}

	assert.Equal(t, mockConn, client.GetConnection())
}

func TestGetUserInfo(t *testing.T) {
	client := &clientImpl{user: "testuser@resource1"}

	username, resource := client.GetUserInfo()
	assert.Equal(t, "testuser", username)
	assert.Equal(t, "resource1", resource)
}

func TestWriteToChannel(t *testing.T) {
	client := &clientImpl{write: make(chan []byte, 1)}
	data := []byte("test message")

	client.WriteToChannel(&data)

	select {
	case msg := <-client.write:
		assert.Equal(t, data, msg)
	default:
		t.Fatal("message was not written to the channel")
	}
}