package hub

import "time"

// PayloadType contains all the possible payload that a client can send
type PayloadType string

var (
	ChatMessageType      = PayloadType("chat_message")
	GroupChatMessageType = PayloadType("group_chat_message")
	TypingStatusType     = PayloadType("typing_status")
	EditMessageType      = PayloadType("edit_message")
	DeleteMessageType    = PayloadType("delete_message")
)

// BasePayload is used to identify what's the actual payload that user has sent
//
// possible payload types that a client can send are:
// "chat_message": 1 to 1 individual chat message
// "group_chat_message": messages that are send to group chats
// "typing_status"
// "edit_message"
// "delete_message"
type BasePayload struct {
	Type PayloadType `json:"type"`
}

// ChatMessage is payload for "chat_message"
type ChatMessage struct {
	Type    PayloadType `json:"type"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Id      string      `json:"id"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
	SendAt  time.Time   `json:"sendAt"`
}

// TypingStatus is payload for "typing_status"
type TypingStatus struct {
	Type PayloadType `json:"type"`
	From string      `json:"from"`
	To   string      `json:"to"`
}

// EditMessage is payload for "edit_message"
type EditMessage struct {
	Type     PayloadType `json:"type"`
	From     string      `json:"from"`
	To       string      `json:"to"`
	Id       string      `json:"id"`
	Body     string      `json:"body"`
	EditedOn time.Time   `json:"editedOn"`
}

// DeleteMessage is payload for "delete_message"
type DeleteMessage struct {
	Type     PayloadType `json:"type"`
	From     string      `json:"from"`
	To       string      `json:"to"`
	Id       string      `json:"id"`
	Everyone bool        `json:"everyone"`
}

// GroupChatMessage is payload for "group_chat_message"
type GroupChatMessage struct {
	Type    PayloadType `json:"type"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Id      string      `json:"id"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
	SendAt  time.Time   `json:"sendAt"`
}