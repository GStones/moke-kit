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

func TestLocalSubscribeReturnsValidSubscription(t *testing.T) {
	mq := NewMessageQueue(zap.NewNop(), 16, false, false)
	var got atomic.Int32
	sub, err := mq.Subscribe(context.Background(), "topic-a", func(msg miface.Message, err error) common.ConsumptionCode {
		got.Add(1)
		return common.ConsumeAck
	})
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}
	if sub == nil || !sub.IsValid() {
		t.Fatal("expected non-nil valid subscription")
	}

	if err := mq.Publish("topic-a", miface.WithBytes([]byte("hello"))); err != nil {
		t.Fatalf("Publish error: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && got.Load() == 0 {
		time.Sleep(10 * time.Millisecond)
	}
	if got.Load() == 0 {
		t.Fatal("expected message to be consumed")
	}

	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("Unsubscribe error: %v", err)
	}
	if sub.IsValid() {
		t.Fatal("expected subscription invalid after unsubscribe")
	}
}
