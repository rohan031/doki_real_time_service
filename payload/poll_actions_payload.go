package payload

import "doki.co.in/doki_real_time_service/utils"

const (
	pollsSubscriptionType = payloadType("poll_subscription")
	pollsVotesUpdateType  = payloadType("poll_votes_update")
)

type pollsSubscription struct {
	Type   payloadType `json:"type" validate:"required"`
	From   string      `json:"from" validate:"required"`
	PollId string      `json:"pollId" validate:"required"`
}

func (payload *pollsSubscription) SendPayload(_ *[]byte, h hub, senderResource string) {
	completeUser := utils.CreateUserFromUsernameAndResource(payload.From, senderResource)

	h.Subscribe(payload.PollId, completeUser, false)
}

type pollsVotesUpdate struct {
	Type   payloadType `json:"type" validate:"required"`
	From   string      `json:"from" validate:"required"`
	PollId string      `json:"pollId" validate:"required"`
	Votes  []int       `json:"votes" validate:"required"`
}

func (payload *pollsVotesUpdate) SendPayload(data *[]byte, h hub, senderResource string) {
	// get subscribers and send
	subscribers := h.GetSubscribers(payload.PollId)
	for subscriber := range subscribers {
		conn := h.GetIndividualClient(subscriber)
		if conn == nil {
			h.Unsubscribe(payload.PollId, subscriber)
			continue
		}

		user, res := conn.GetUserInfo()
		if user != payload.From || res != senderResource {
			// send votes update to user
			conn.WriteToChannel(data)
		}
	}
}