package subscription

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestCancelSubscription(t *testing.T) {
	called := false
	sub := NewCancelSubscription(func() { called = true })
	if !sub.IsValid() {
		t.Fatal("expected valid subscription")
	}
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("Unsubscribe error: %v", err)
	}
	if !called {
		t.Fatal("expected cancel to be called")
	}
	if sub.IsValid() {
		t.Fatal("expected invalid after unsubscribe")
	}
	// idempotent
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("second Unsubscribe error: %v", err)
	}
}

func TestCancelSubscriptionConcurrentUnsubscribe(t *testing.T) {
	t.Parallel()
	var cancels atomic.Int32
	sub := NewCancelSubscription(func() { cancels.Add(1) })

	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := sub.Unsubscribe(); err != nil {
				t.Errorf("Unsubscribe: %v", err)
			}
		}()
	}
	wg.Wait()

	if cancels.Load() != 1 {
		t.Fatalf("cancel calls=%d want 1", cancels.Load())
	}
	if sub.IsValid() {
		t.Fatal("expected invalid after concurrent unsubscribe")
	}
}
