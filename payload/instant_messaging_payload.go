package payload

import (
	"time"
)

const (
	chatMessageType      = payloadType("chat_message")
	groupChatMessageType = payloadType("group_chat_message")
	typingStatusType     = payloadType("typing_status")
	editMessageType      = payloadType("edit_message")
	deleteMessageType    = payloadType("delete_message")
)

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

func (message *chatMessage) SendPayload(data *[]byte, h *hub, senderResource string) {
	recipient := message.To
	sender := message.From

	// this prevents sending messages twice when user sends self messages
	if recipient != sender {
		recipientConnectedClients := (*h).GetAllConnectedClients(recipient)
		for _, conn := range recipientConnectedClients {
			(*conn).WriteToChannel(data)
		}
	}

	senderConnectedClients := (*h).GetAllConnectedClients(sender)
	for res, conn := range senderConnectedClients {
		if res != senderResource {
			(*conn).WriteToChannel(data)
		}
	}

}

// typingStatus is payload for "typing_status"
type typingStatus struct {
	Type payloadType `json:"type" validate:"required"`
	From string      `json:"from" validate:"required"`
	To   string      `json:"to" validate:"required"`
}

func (status *typingStatus) SendPayload(data *[]byte, h *hub, _ string) {
	recipient := status.To
	if recipient == status.From {
		return
	}

	recipientConnectedClients := (*h).GetAllConnectedClients(recipient)
	for _, conn := range recipientConnectedClients {
		(*conn).WriteToChannel(data)
	}
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

func (message *editMessage) SendPayload(data *[]byte, h *hub, senderResource string) {
	recipient := message.To
	sender := message.From

	if recipient != sender {
		recipientConnectedClients := (*h).GetAllConnectedClients(recipient)
		for _, conn := range recipientConnectedClients {
			(*conn).WriteToChannel(data)
		}
	}

	senderConnectedClients := (*h).GetAllConnectedClients(sender)
	for res, conn := range senderConnectedClients {
		if res != senderResource {
			(*conn).WriteToChannel(data)
		}
	}
}

// deleteMessage is payload for "delete_message"
type deleteMessage struct {
	Type     payloadType `json:"type" validate:"required"`
	From     string      `json:"from" validate:"required"`
	To       string      `json:"to" validate:"required"`
	Id       []string    `json:"id" validate:"required"`
	Everyone bool        `json:"everyone,string"`
}

func (message *deleteMessage) SendPayload(data *[]byte, h *hub, senderResource string) {
	recipient := message.To
	sender := message.From

	if message.Everyone && recipient != sender {
		recipientConnectedClients := (*h).GetAllConnectedClients(recipient)
		for _, conn := range recipientConnectedClients {
			(*conn).WriteToChannel(data)
		}
	}

	senderConnectedClients := (*h).GetAllConnectedClients(sender)
	for res, conn := range senderConnectedClients {
		if res != senderResource {
			(*conn).WriteToChannel(data)
		}
	}
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