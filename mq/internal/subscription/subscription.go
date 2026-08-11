package subscription

import (
	"sync/atomic"
)

// CancelSubscription implements miface.Subscription by cancelling a per-subscribe
// context instead of closing a shared Watermill subscriber.
type CancelSubscription struct {
	cancel func()
	valid  atomic.Bool
	closed atomic.Bool
}

// NewCancelSubscription creates a valid subscription that cancels via cancel.
func NewCancelSubscription(cancel func()) *CancelSubscription {
	s := &CancelSubscription{cancel: cancel}
	s.valid.Store(true)
	return s
}

// IsValid reports whether the subscription is still active.
func (s *CancelSubscription) IsValid() bool {
	return s != nil && s.valid.Load()
}

// Unsubscribe cancels the subscription. It is safe to call more than once.
func (s *CancelSubscription) Unsubscribe() error {
	if s == nil {
		return nil
	}
	if !s.closed.CompareAndSwap(false, true) {
		return nil
	}
	s.valid.Store(false)
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}
