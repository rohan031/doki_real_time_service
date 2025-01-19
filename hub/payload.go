package hub

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

var validate = validator.New()

// payloadType contains all the possible payload that a client can send
type payloadType string

const (
	chatMessageType      = payloadType("chat_message")
	groupChatMessageType = payloadType("group_chat_message")
	typingStatusType     = payloadType("typing_status")
	editMessageType      = payloadType("edit_message")
	deleteMessageType    = payloadType("delete_message")
)

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

// chatMessage is payload for "chat_message"
type chatMessage struct {
	Type    payloadType `json:"type" validate:"required"`
	From    string      `json:"from" validate:"required"`
	To      string      `json:"to" validate:"required"`
	Id      string      `json:"id" validate:"required"`
	Subject string      `json:"subject" validate:"required"`
	Body    string      `json:"body" validate:"required"`
	SendAt  time.Time   `json:"sendAt" validate:"required"`
}

// typingStatus is payload for "typing_status"
type typingStatus struct {
	Type payloadType `json:"type" validate:"required"`
	From string      `json:"from" validate:"required"`
	To   string      `json:"to" validate:"required"`
}

// editMessage is payload for "edit_message"
type editMessage struct {
	Type     payloadType `json:"type" validate:"required"`
	From     string      `json:"from" validate:"required"`
	To       string      `json:"to" validate:"required"`
	Id       string      `json:"id" validate:"required"`
	Body     string      `json:"body" validate:"required"`
	EditedOn time.Time   `json:"editedOn" validate:"required"`
}

// deleteMessage is payload for "delete_message"
type deleteMessage struct {
	Type     payloadType `json:"type" validate:"required"`
	From     string      `json:"from" validate:"required"`
	To       string      `json:"to" validate:"required"`
	Id       string      `json:"id" validate:"required"`
	Everyone bool        `json:"everyone" validate:"required"`
}

//// groupChatMessage is payload for "group_chat_message"
//type groupChatMessage struct {
//	Type    payloadType `json:"type" validate:"required"`
//	From    string      `json:"from" validate:"required"`
//	To      string      `json:"to" validate:"required"`
//	Id      string      `json:"id" validate:"required"`
//	Subject string      `json:"subject" validate:"required"`
//	Body    string      `json:"body" validate:"required"`
//	SendAt  time.Time   `json:"sendAt" validate:"required"`
//}

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