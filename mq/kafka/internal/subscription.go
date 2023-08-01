package internal

import (
	"sync"

	"github.com/gstones/platform/services/common/mq"
)

type subscription struct {
	reader *reader
	mutex  sync.Mutex
}

func NewSubscription(reader *reader) *subscription {
	return &subscription{reader: reader}
}

func (s *subscription) IsValid() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.reader == nil {
		return false
	} else {
		return true
	}
}

func (s *subscription) Unsubscribe() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.reader == nil {
		return mq.ErrInvalidSubscription
	} else {
		s.reader.ReturnToPool()
		s.reader = nil

		return nil
	}
}
