package subscription

import "testing"

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
