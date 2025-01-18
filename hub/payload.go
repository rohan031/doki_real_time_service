package hub

import "time"

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
	Type payloadType `json:"type"`
}

// chatMessage is payload for "chat_message"
type chatMessage struct {
	Type    payloadType `json:"type"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Id      string      `json:"id"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
	SendAt  time.Time   `json:"sendAt"`
}

// typingStatus is payload for "typing_status"
type typingStatus struct {
	Type payloadType `json:"type"`
	From string      `json:"from"`
	To   string      `json:"to"`
}

// editMessage is payload for "edit_message"
type editMessage struct {
	Type     payloadType `json:"type"`
	From     string      `json:"from"`
	To       string      `json:"to"`
	Id       string      `json:"id"`
	Body     string      `json:"body"`
	EditedOn time.Time   `json:"editedOn"`
}

// deleteMessage is payload for "delete_message"
type deleteMessage struct {
	Type     payloadType `json:"type"`
	From     string      `json:"from"`
	To       string      `json:"to"`
	Id       string      `json:"id"`
	Everyone bool        `json:"everyone"`
}

// groupChatMessage is payload for "group_chat_message"
type groupChatMessage struct {
	Type    payloadType `json:"type"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Id      string      `json:"id"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
	SendAt  time.Time   `json:"sendAt"`
}