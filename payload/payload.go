package payload

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"log"
)

var validate = validator.New()

// payloadType contains all the possible payload that a client can send
type payloadType string

type Payload interface {
	SendPayload(*[]byte, *hub, string)
}

// this is payload client interface to handle sending payload
type client interface {
	GetConnection() *websocket.Conn
	GetUserInfo() (string, string)
	WriteToChannel(*[]byte)
}

// this is payload hub interface required to send data
type hub interface {
	GetAllConnectedClients(string) map[string]*client
	GetIndividualClient(string) *client
}

type InvalidPayload struct {
	reason string
}

func (p *InvalidPayload) Error() string {
	return p.reason
}

// basePayload is used to identify what's the actual payload that user has sent
//
// possible payload types that a client can send are:
// "chat_message": 1 to 1 individual chat message
// "group_chat_message": messages that are send to group chats
// "typing_status"
// "edit_message"
// "delete_message"
type basePayload struct {
	Type payloadType `json:"type" validate:"required"`
	From string      `json:"from" validate:"required"`
}

// unmarshalAndValidate first unmarshal payload json and validates it
func unmarshalAndValidate[T any](payload *[]byte, target *T) bool {
	if err := json.Unmarshal(*payload, target); err != nil {
		log.Printf("error unmarshalling payload: %v\n", err)
		return false
	}

	if err := validate.Struct(target); err != nil {
		log.Println("missing required field in payload.")
		return false
	}

	return true
}

// CreatePayload is factory method to create different payloads based on type
func CreatePayload(data *[]byte, from string) (Payload, error) {
	var base basePayload
	if !unmarshalAndValidate(data, &base) {
		return nil, &InvalidPayload{
			reason: "Invalid payload received.",
		}
	}

	if base.From != from {
		return nil, &InvalidPayload{
			reason: "Client username and payload from mismatch.",
		}
	}

	switch base.Type {
	case chatMessageType:
		var message chatMessage
		if unmarshalAndValidate(data, &message) {
			return &message, nil
		}

	case typingStatusType:
		var status typingStatus
		if unmarshalAndValidate(data, &status) {
			return &status, nil
		}
	case editMessageType:
		var message editMessage
		if unmarshalAndValidate(data, &message) {
			return &message, nil
		}
	case deleteMessageType:
		var message deleteMessage
		if unmarshalAndValidate(data, &message) {
			return &message, nil
		}
	default:
		return nil, &InvalidPayload{
			reason: "unknown payload type",
		}
	}

	return nil, &InvalidPayload{
		reason: "Invalid payload received",
	}
}