package nats

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
)

func startTestNATS(t *testing.T) (addr string, shutdown func()) {
	t.Helper()
	opts := &natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	}
	s, err := natsserver.NewServer(opts)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	go s.Start()
	if !s.ReadyForConnections(5 * time.Second) {
		t.Fatal("nats server not ready")
	}
	return s.ClientURL(), func() { s.Shutdown() }
}

func TestSubscribeCancelUnsubscribes(t *testing.T) {
	addr, shutdown := startTestNATS(t)
	defer shutdown()

	mq, err := NewMessageQueue(zap.NewNop(), addr)
	if err != nil {
		t.Fatalf("NewMessageQueue: %v", err)
	}

	var received atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	sub, err := mq.Subscribe(ctx, "life-topic", func(msg miface.Message, err error) common.ConsumptionCode {
		received.Add(1)
		return common.ConsumeAck
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	if err := mq.Publish("life-topic", miface.WithBytes([]byte("one"))); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if received.Load() > 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if received.Load() == 0 {
		t.Fatal("expected at least one message")
	}

	cancel()
	_ = sub.Unsubscribe()

	before := received.Load()
	_ = mq.Publish("life-topic", miface.WithBytes([]byte("two")))
	time.Sleep(100 * time.Millisecond)
	if received.Load() != before {
		t.Fatalf("received after cancel: got %d want %d", received.Load(), before)
	}
	if sub.IsValid() {
		t.Fatal("expected subscription invalid after Unsubscribe")
	}
}

func TestSubscribeEmptyTopic(t *testing.T) {
	addr, shutdown := startTestNATS(t)
	defer shutdown()

	mq, err := NewMessageQueue(zap.NewNop(), addr)
	if err != nil {
		t.Fatalf("NewMessageQueue: %v", err)
	}
	if _, err := mq.Subscribe(context.Background(), "", func(miface.Message, error) common.ConsumptionCode {
		return common.ConsumeAck
	}); err == nil {
		t.Fatal("expected empty topic error")
	}
}
