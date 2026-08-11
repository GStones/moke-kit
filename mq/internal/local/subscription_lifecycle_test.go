package local

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
)

func TestUnsubscribeDoesNotBreakOtherSubscriptions(t *testing.T) {
	mq := NewMessageQueue(zap.NewNop(), 16, false, false)

	var aCount, bCount atomic.Int32
	subA, err := mq.Subscribe(context.Background(), "topic-a", func(msg miface.Message, err error) common.ConsumptionCode {
		aCount.Add(1)
		return common.ConsumeAck
	})
	if err != nil {
		t.Fatalf("subscribe A: %v", err)
	}
	subB, err := mq.Subscribe(context.Background(), "topic-b", func(msg miface.Message, err error) common.ConsumptionCode {
		bCount.Add(1)
		return common.ConsumeAck
	})
	if err != nil {
		t.Fatalf("subscribe B: %v", err)
	}

	if err := subA.Unsubscribe(); err != nil {
		t.Fatalf("unsubscribe A: %v", err)
	}
	if subA.IsValid() {
		t.Fatal("A should be invalid")
	}
	if !subB.IsValid() {
		t.Fatal("B should remain valid")
	}

	if err := mq.Publish("topic-b", miface.WithBytes([]byte("b"))); err != nil {
		t.Fatalf("publish B: %v", err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && bCount.Load() == 0 {
		time.Sleep(10 * time.Millisecond)
	}
	if bCount.Load() == 0 {
		t.Fatal("expected topic-b message after unsubscribing A")
	}
	if aCount.Load() != 0 {
		t.Fatalf("topic-a should not receive messages, got %d", aCount.Load())
	}
}
