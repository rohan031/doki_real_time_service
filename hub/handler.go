package hub

// handleChatMessagePayload sends incoming message to all the receivers active connection
// it also sends the message to other resources of the senders
// it will also add the message to queue to add to the scyllaDb
func handleChatMessagePayload(h *Hub, message *chatMessage, payload *[]byte, username, resource string) {
	recipient := message.To

	if recipient != username {
		recipientConnectedClients := h.getAllConnectedClients(recipient)
		for _, conn := range recipientConnectedClients {
			conn.write <- *payload
		}
	}

	senderConnectedClients := h.getAllConnectedClients(username)
	for res, conn := range senderConnectedClients {
		if res != resource {
			conn.write <- *payload
		}
	}
}

// handleTypingStatusPayload sends typing status to all the connected recipient
func handleTypingStatusPayload(h *Hub, status *typingStatus, payload *[]byte) {
	recipient := status.To
	if recipient == status.From {
		return
	}

	recipientConnectedClients := h.getAllConnectedClients(recipient)
	for _, conn := range recipientConnectedClients {
		conn.write <- *payload
	}
}

func handleDeleteMessagePayload(h *Hub, message *deleteMessage, payload *[]byte, username, resource string) {
	senderConnectedClients := h.getAllConnectedClients(username)
	recipient := message.To

	if message.Everyone && recipient != username {
		recipientConnectedClients := h.getAllConnectedClients(recipient)
		for _, conn := range recipientConnectedClients {
			conn.write <- *payload
		}
	}

	for res, conn := range senderConnectedClients {
		if res != resource {
			conn.write <- *payload
		}
	}
}

func handleEditMessagePayload(h *Hub, message *editMessage, payload *[]byte, username, resource string) {
	recipient := message.To

	if recipient != username {
		recipientConnectedClients := h.getAllConnectedClients(recipient)
		for _, conn := range recipientConnectedClients {
			conn.write <- *payload
		}
	}

	senderConnectedClients := h.getAllConnectedClients(username)
	for res, conn := range senderConnectedClients {
		if res != resource {
			conn.write <- *payload
		}
	}
}