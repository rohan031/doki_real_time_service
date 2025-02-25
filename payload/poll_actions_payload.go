package payload

const (
	pollsSubscriptionType = payloadType("poll_subscription")
	pollsVotesUpdateType  = payloadType("poll_votes_update")
)

type pollsSubscription struct {
	Type   payloadType `json:"type" validate:"required"`
	From   string      `json:"from" validate:"required"`
	PollId string      `json:"pollId" validate:"required"`
}

type pollsVotesUpdate struct {
	Type   payloadType `json:"type" validate:"required"`
	From   string      `json:"from" validate:"required"`
	PollId string      `json:"pollId" validate:"required"`
	Votes  []int       `json:"votes" validate:"required"`
}