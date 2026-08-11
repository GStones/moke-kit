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
	started := make(chan struct{})
	var once sync.Once
	sub, err := mq.Subscribe(context.Background(), "race-topic", func(msg miface.Message, err error) common.ConsumptionCode {
		once.Do(func() { close(started) })
		received.Add(1)
		return common.ConsumeAck
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	const publishers = 8
	const messages = 40
	var wg sync.WaitGroup
	wg.Add(publishers)

	// Unsubscribe only after the first message is observed, so publish/unsub overlap under -race.
	unsubDone := make(chan struct{})
	go func() {
		defer close(unsubDone)
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Error("timed out waiting for first message before Unsubscribe")
			return
		}
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
	<-unsubDone

	if sub.IsValid() {
		t.Fatal("expected subscription invalid after Unsubscribe")
	}
	if received.Load() == 0 {
		t.Fatal("expected at least one message before Unsubscribe")
	}
}
