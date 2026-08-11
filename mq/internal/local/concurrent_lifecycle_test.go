package local

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
)

func TestConcurrentPublishAndUnsubscribe(t *testing.T) {
	t.Parallel()
	mq := NewMessageQueue(zap.NewNop(), 64, false, false)

	var received atomic.Int32
	sub, err := mq.Subscribe(context.Background(), "race-topic", func(msg miface.Message, err error) common.ConsumptionCode {
		received.Add(1)
		return common.ConsumeAck
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	const publishers = 8
	const messages = 20
	var wg sync.WaitGroup
	wg.Add(publishers + 1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		_ = sub.Unsubscribe()
	}()
	for i := 0; i < publishers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < messages; j++ {
				_ = mq.Publish("race-topic", miface.WithBytes([]byte("x")))
			}
		}()
	}
	wg.Wait()

	// Drain briefly; goal is race-detector cleanliness, not exact delivery count.
	time.Sleep(50 * time.Millisecond)
	_ = received.Load()
	if sub.IsValid() {
		t.Fatal("expected subscription invalid after Unsubscribe")
	}
}
