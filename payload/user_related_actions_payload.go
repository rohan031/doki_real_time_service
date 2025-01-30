package payload

import "time"

const (
	userSendFriendRequestType     = payloadType("user_send_friend_request")
	userAcceptedFriendRequestType = payloadType("user_accepted_friend_request")
	userRemovesFriendRelationType = payloadType("user_removes_friend_relation")
)

type userSendFriendRequest struct {
	Type        payloadType `json:"type" validate:"required"`
	From        string      `json:"from" validate:"required"`
	To          string      `json:"to" validate:"required"`
	RequestedBy string      `json:"requestedBy" validate:"required"`
	AddedOn     time.Time   `json:"addedOn" validate:"required"`
}

func (req *userSendFriendRequest) SendPayload(data *[]byte, h hub, senderResource string) {
	userSendingRequest := req.From
	userToSendRequest := req.To

	if userSendingRequest == userToSendRequest {
		// invalid state
		return
	}

	userSendingRequestConnectedClients := h.GetAllConnectedClients(userSendingRequest)
	for res, conn := range userSendingRequestConnectedClients {
		if res != senderResource {
			conn.WriteToChannel(data)
		}
	}

	userToSendRequestConnectedClients := h.GetAllConnectedClients(userToSendRequest)
	for _, conn := range userToSendRequestConnectedClients {
		conn.WriteToChannel(data)
	}
}

type userAcceptFriendRequest struct {
	Type        payloadType `json:"type" validate:"required"`
	From        string      `json:"from" validate:"required"`
	To          string      `json:"to" validate:"required"`
	RequestedBy string      `json:"requestedBy" validate:"required"`
	AddedOn     time.Time   `json:"addedOn" validate:"required"`
}

func (req *userAcceptFriendRequest) SendPayload(data *[]byte, h hub, senderResource string) {
	userAcceptingRequest := req.From
	userToAcceptRequest := req.To

	if userAcceptingRequest == userToAcceptRequest {
		// invalid state
		return
	}

	userAcceptingRequestConnectedClients := h.GetAllConnectedClients(userAcceptingRequest)
	for res, conn := range userAcceptingRequestConnectedClients {
		if res != senderResource {
			conn.WriteToChannel(data)
		}
	}

	userToAcceptRequestConnectedClients := h.GetAllConnectedClients(userToAcceptRequest)
	for _, conn := range userToAcceptRequestConnectedClients {
		conn.WriteToChannel(data)
	}
}

type userRemovesFriendRelation struct {
	Type payloadType `json:"type" validate:"required"`
	From string      `json:"from" validate:"required"`
	To   string      `json:"to" validate:"required"`
}

func (req *userRemovesFriendRelation) SendPayload(data *[]byte, h hub, senderResource string) {
	from := req.From
	to := req.To

	if from == to {
		// invalid state
		return
	}

	fromConnectedClients := h.GetAllConnectedClients(from)
	for res, conn := range fromConnectedClients {
		if res != senderResource {
			conn.WriteToChannel(data)
		}
	}

	toConnectedClients := h.GetAllConnectedClients(to)
	for _, conn := range toConnectedClients {
		conn.WriteToChannel(data)
	}
}