package payload

import (
	"doki.co.in/doki_real_time_service/client"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
)

var validate = validator.New()

// payloadType contains all the possible payload that a client can send
type payloadType string

type Payload interface {
	// SendPayload expects raw payload, hub, and senders resource
	SendPayload(*[]byte, hub, string)
}

// this is payload hub interface required to send data
type hub interface {
	GetAllConnectedClients(string) map[string]client.Client

	GetIndividualClient(string) client.Client
}

type InvalidPayload struct {
	reason string
}

func (p *InvalidPayload) Error() string {
	return p.reason
}

// basePayload is used to identify what's the actual payload that user has sent
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

func InitPayload() {
	// instant messaging payloads
	payloadMap[chatMessageType] = func() Payload { return &chatMessage{} }
	payloadMap[typingStatusType] = func() Payload { return &typingStatus{} }
	payloadMap[editMessageType] = func() Payload { return &editMessage{} }
	payloadMap[deleteMessageType] = func() Payload { return &deleteMessage{} }

	// user to user action payload
	payloadMap[userSendFriendRequestType] = func() Payload { return &userSendFriendRequest{} }
	payloadMap[userAcceptedFriendRequestType] = func() Payload { return &userAcceptFriendRequest{} }
	payloadMap[userRemovesFriendRelationType] = func() Payload { return &userRemovesFriendRelation{} }

	// user profile self action payload
	payloadMap[userUpdateProfileType] = func() Payload { return &userUpdateProfile{} }
	payloadMap[userCreateRootNodeType] = func() Payload { return &userCreateRootNode{} }
}