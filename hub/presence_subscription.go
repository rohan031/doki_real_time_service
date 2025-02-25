package hub

import (
	"doki.co.in/doki_real_time_service/payload"
	"doki.co.in/doki_real_time_service/utils"
	"maps"
)

// sendPresence sends user presence updates to all the subscribed users
func (h *Hub) sendPresence(online bool, username string) {
	// find user in subscription and send status change
	usersSubscribed, ok := h.subscription.subscriptions[username]
	if !ok {
		return
	}

	for completeUser := range maps.Keys(usersSubscribed) {
		conn := h.GetIndividualClient(completeUser)
		if conn == nil {
			h.Unsubscribe(username, completeUser)
			continue
		}

		user, resource := utils.GetUsernameAndResourceFromUser(completeUser)
		presencePayload := payload.CreatePresencePayload(username, user, online)

		data := utils.PayloadToJson(presencePayload)
		if data != nil {
			presencePayload.SendPayload(data, h, resource)
		}
	}
}

// sendInitialPresence sends the given user presence on initial subscription
func (h *Hub) sendInitialPresence(userPresence string, completeUser string) {
	_, ok := h.clients[userPresence]

	username, resource := utils.GetUsernameAndResourceFromUser(completeUser)
	presencePayload := payload.CreatePresencePayload(userPresence, username, ok)

	data := utils.PayloadToJson(presencePayload)
	if data != nil {
		presencePayload.SendPayload(data, h, resource)
	}
}