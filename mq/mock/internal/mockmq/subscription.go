package mockmq

import (
	"sync"

	"github.com/gstones/platform/services/common/mq"
)

type subscription struct {
	sync.Mutex
	queue    string
	handler  MessageHandler
	subGroup Unsubscriber
}

func NewSubscription(
	queue string,
	handler MessageHandler,
	subGroup Unsubscriber,
) Subscription {
	return &subscription{queue: queue, handler: handler, subGroup: subGroup}
}

func (s *subscription) IsValid() bool {
	s.Lock()
	defer s.Unlock()

	if s.subGroup == nil || s.handler == nil {
		return false
	} else {
		return true
	}
}

func (s *subscription) Unsubscribe() (err error) {
	s.Lock()

	if s.subGroup == nil || s.handler == nil {
		err = mq.ErrInvalidSubscription
	} else {
		err = s.subGroup.RemoveSubscription(s)

		s.handler = nil
		s.subGroup = nil
	}

	s.Unlock()
	return err
}

func (s *subscription) Queue() (queue string) {
	s.Lock()
	queue = s.queue
	s.Unlock()

	return queue
}

func (s *subscription) SendMessage(message Message) (err error) {
	s.Lock()

	var handler MessageHandler

	if s.subGroup == nil || s.handler == nil {
		err = mq.ErrInvalidSubscription
	} else {
		handler = s.handler
	}

	s.Unlock()

	if handler != nil {
		handler(message)
	}

	return err
}
