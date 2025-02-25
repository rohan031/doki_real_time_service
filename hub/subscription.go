package hub

import "sync"

type subscribers map[string]bool

// nodeSubscription contains complete user to send the subscription updates
// node identifier -> complete user
// [complete user] has subscribed to node updates
type nodeSubscription map[string]subscribers

type subscription struct {
	sync.RWMutex
	subscriptions nodeSubscription
}

// Subscribe the complete user to node
// for user node it is presence update
// for poll node it is poll votes
// subscriber is complete user
// userPresence determines if nodeIdentifier is username or not
func (h *Hub) Subscribe(nodeIdentifier, subscriber string, userPresence bool) {
	h.subscription.Lock()
	defer h.subscription.Unlock()

	if h.subscription.subscriptions[nodeIdentifier] == nil {
		h.subscription.subscriptions[nodeIdentifier] = make(subscribers)
	}

	h.subscription.subscriptions[nodeIdentifier][subscriber] = true
	if userPresence {
		h.sendInitialPresence(nodeIdentifier, subscriber)
	}

	// adding subscription to client for cleanup
	conn := h.GetIndividualClient(subscriber)
	if conn != nil {
		conn.AddSubscription(nodeIdentifier)
	}
}

func (h *Hub) Unsubscribe(nodeIdentifier, subscriber string) {
	h.subscription.Lock()
	defer h.subscription.Unlock()

	_, ok := h.subscription.subscriptions[nodeIdentifier]
	if !ok {
		return
	}

	delete(h.subscription.subscriptions[nodeIdentifier], subscriber)
	if len(h.subscription.subscriptions[nodeIdentifier]) == 0 {
		delete(h.subscription.subscriptions, nodeIdentifier)
	}
}