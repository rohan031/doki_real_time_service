package hub

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetUserInfo() (string, string) {
	args := m.Called()
	return args.String(0), args.String(1)
}

func (m *MockClient) GetConnection() *websocket.Conn {
	return nil
}

func (m *MockClient) WriteToChannel(data *[]byte) {
	m.Called(data)
}

func TestAddClient(t *testing.T) {
	h := CreateHub(nil) // Create hub instance
	mockClient := new(MockClient)
	mockClient.On("GetUserInfo").Return("testUser", "resource1")

	h.addClient("testUser@resource1", mockClient)

	assert.NotNil(t, h.GetIndividualClient("testUser@resource1"))
}

func TestRemoveClient(t *testing.T) {
	h := CreateHub(nil)

	mockClient := new(MockClient)
	mockClient.On("GetUserInfo").Return("testUser", "resource1")

	h.addClient("testUser@resource1", mockClient)
	assert.NotNil(t, h.GetIndividualClient("testUser@resource1"))

	h.removeClient(mockClient)
	assert.Nil(t, h.GetIndividualClient("testUser@resource1"))
}

func TestGetAllConnectedClients(t *testing.T) {
	h := CreateHub(nil)
	mockClient1 := new(MockClient)
	mockClient2 := new(MockClient)
	mockClient1.On("GetUserInfo").Return("testUser", "resource1")
	mockClient2.On("GetUserInfo").Return("testUser", "resource2")

	h.addClient("testUser@resource1", mockClient1)
	h.addClient("testUser@resource2", mockClient2)

	clients := h.GetAllConnectedClients("testUser")
	assert.Len(t, clients, 2)
}

func TestGetIndividualClient(t *testing.T) {
	h := CreateHub(nil)
	mockClient := new(MockClient)

	h.addClient("testuser@resource1", mockClient)

	client := h.GetIndividualClient("testuser@resource1")
	assert.Equal(t, mockClient, client)

	clientNil := h.GetIndividualClient("unknown@resource")
	assert.Nil(t, clientNil)
}