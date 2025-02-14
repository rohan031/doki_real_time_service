package hub

import (
	"doki.co.in/doki_real_time_service/payload"
	"doki.co.in/doki_real_time_service/utils"
	"log"
	"maps"
	"sync"
)

type presenceList map[string]bool

// presenceSubscription contains complete user to send the user status
// username -> complete user
// [complete user] has subscribed to username
type presenceSubscription map[string]presenceList

type userPresence struct {
	sync.RWMutex
	subscriptions presenceSubscription
}

// sendPresence sends user presence updates to all the subscribed users
func (h *Hub) sendPresence(online bool, username string) {
	// find user in subscription and send status change
	usersSubscribed, ok := h.presenceSubscription.subscriptions[username]
	if !ok {
		return
	}

	for completeUser := range maps.Keys(usersSubscribed) {
		conn := h.GetIndividualClient(completeUser)
		if conn == nil {
			h.UnsubscribeUserPresence(username, completeUser)
			continue
		}

		user, resource := utils.GetUsernameAndResourceFromUser(completeUser)
		presencePayload := payload.CreatePresencePayload(username, user, online)

		log.Printf("\nSending user presence: %v, %v\n", username, completeUser)
		data := utils.PayloadToJson(presencePayload)
		if data != nil {
			presencePayload.SendPayload(data, h, resource)
		}
	}
}

// sendInitialPresence sends the given user presence on initial subscription
func (h *Hub) sendInitialPresence(userPresence string, completeUser string) {
	_, ok := h.clients[userPresence]
	conn := h.GetIndividualClient(completeUser)

	if conn == nil {
		h.UnsubscribeUserPresence(userPresence, completeUser)
		return
	}

	log.Printf("\nSending initial presence: %v, %v\n", userPresence, completeUser)
	username, resource := utils.GetUsernameAndResourceFromUser(completeUser)
	presencePayload := payload.CreatePresencePayload(userPresence, username, ok)

	data := utils.PayloadToJson(presencePayload)
	if data != nil {
		presencePayload.SendPayload(data, h, resource)
	}
}

// SubscribeUserPresence subscribes the given complete user to the user presence updates
func (h *Hub) SubscribeUserPresence(userToSubscribe string, completeUser string) {
	h.presenceSubscription.Lock()
	defer h.presenceSubscription.Unlock()

	log.Printf("\nSubscribing to user presence: %v, %v\n", userToSubscribe, completeUser)
	if h.presenceSubscription.subscriptions[userToSubscribe] == nil {
		h.presenceSubscription.subscriptions[userToSubscribe] = make(presenceList)
	}

	h.presenceSubscription.subscriptions[userToSubscribe][completeUser] = true
	h.sendInitialPresence(userToSubscribe, completeUser)
}

// UnsubscribeUserPresence unsubscribes the given complete user to the user presence updates
func (h *Hub) UnsubscribeUserPresence(userToUnsubscribe string, completeUser string) {

	h.presenceSubscription.Lock()
	defer h.presenceSubscription.Unlock()

	log.Printf("\nUn-Subscribing to user presence: %v, %v\n", userToUnsubscribe, completeUser)
	// find slice
	_, ok := h.presenceSubscription.subscriptions[userToUnsubscribe]
	if !ok {
		return
	}

	delete(h.presenceSubscription.subscriptions[userToUnsubscribe], completeUser)

	if len(h.presenceSubscription.subscriptions[userToUnsubscribe]) == 0 {
		delete(h.presenceSubscription.subscriptions, userToUnsubscribe)
	}
}