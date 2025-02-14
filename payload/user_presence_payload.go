package payload

import (
	"doki.co.in/doki_real_time_service/utils"
)

const (
	userPresenceSubscriptionType = payloadType("user_presence_subscription")
	userPresenceInfoType         = payloadType("user_presence_info")
)

type userPresenceSubscription struct {
	Type      payloadType `json:"type" validate:"required"`
	From      string      `json:"from" validate:"required"`
	User      string      `json:"user" validate:"required"`
	Subscribe bool        `json:"subscribe"`
}

func (payload *userPresenceSubscription) SendPayload(data *[]byte, h hub, senderResource string) {
	completeUser := utils.CreateUserFromUsernameAndResource(payload.From, senderResource)

	if payload.Subscribe {
		h.SubscribeUserPresence(payload.User, completeUser)
	} else {
		h.UnsubscribeUserPresence(payload.User, completeUser)
	}
}

// only server sends this
type userPresenceInfoPayload struct {
	Type   payloadType `json:"type"`
	To     string      `json:"to"`
	User   string      `json:"user"`
	Online bool        `json:"online"`
}

func (payload *userPresenceInfoPayload) SendPayload(data *[]byte, h hub, userResource string) {
	completeUser := utils.CreateUserFromUsernameAndResource(payload.To, userResource)

	conn := h.GetIndividualClient(completeUser)
	if conn != nil {
		conn.WriteToChannel(data)
	}
}