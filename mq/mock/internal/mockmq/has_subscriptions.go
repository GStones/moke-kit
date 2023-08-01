package mockmq

import (
	"sync"

	"github.com/gstones/platform/services/common/mq"
)

type HasSubscriptions struct {
	sync.RWMutex
	subs []Subscription
}

func (h *HasSubscriptions) NewSubscription(queue string, handler MessageHandler) Subscription {
	h.Lock()

	sub := NewSubscription(queue, handler, h)
	h.subs = append(h.subs, sub)

	h.Unlock()

	return sub
}

func (h *HasSubscriptions) RemoveSubscription(sub *subscription) (err error) {
	h.Lock()

	index := -1
	for i, s := range h.subs {
		if s == sub {
			index = i
			break
		}
	}

	if index >= 0 {
		h.subs[index] = h.subs[len(h.subs)-1]
		h.subs[len(h.subs)-1] = nil
		h.subs = h.subs[:len(h.subs)-1]
	} else {
		err = mq.ErrSubNotFound
	}

	h.Unlock()

	return err
}

func (h *HasSubscriptions) Subscriptions() []Subscription {
	h.RLock()
	subscriptions := h.subs
	h.RUnlock()

	return subscriptions
}

func (h *HasSubscriptions) SendMessage(topic string, data []byte) {
	for _, sub := range h.Subscriptions() {
		sub.SendMessage(Message{Topic: topic, Data: data})
	}
}
