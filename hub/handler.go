package hub

// handleChatMessagePayload sends incoming message to all the receivers active connection
// it also sends the message to other resources of the senders
// it will also add the message to queue to add to the scyllaDb
func handleChatMessagePayload(h *Hub, message *chatMessage, payload *[]byte, resource string) {
	sender, recipient := message.From, message.To

	recipientConnectedClients := h.getAllConnectedClients(recipient)
	senderConnectedClients := h.getAllConnectedClients(sender)

	for _, conn := range recipientConnectedClients {
		conn.write <- *payload
	}

	for res, conn := range senderConnectedClients {
		if res != resource {
			conn.write <- *payload
		}
	}
}