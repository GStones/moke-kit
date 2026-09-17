package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/miface"
)

type stubQueue struct {
	subscribe func(ctx context.Context, topic string, handler miface.SubResponseHandler, opts ...miface.SubOption) (miface.Subscription, error)
	publish   func(topic string, opts ...miface.PubOption) error
}

func (s stubQueue) Subscribe(ctx context.Context, topic string, handler miface.SubResponseHandler, opts ...miface.SubOption) (miface.Subscription, error) {
	if s.subscribe != nil {
		return s.subscribe(ctx, topic, handler, opts...)
	}
	return nil, nil
}

func (s stubQueue) Publish(topic string, opts ...miface.PubOption) error {
	if s.publish != nil {
		return s.publish(topic, opts...)
	}
	return nil
}

func TestParseTopic(t *testing.T) {
	tests := []struct {
		topic   string
		want    mqType
		name    string
		wantErr error
	}{
		{topic: "nats://orders", want: nats, name: "orders"},
		{topic: "local://inproc", want: local, name: "inproc"},
		{topic: "kafka://legacy", wantErr: qerrors.ErrMQTypeUnsupported},
		{topic: "nsq://legacy", wantErr: qerrors.ErrMQTypeUnsupported},
		{topic: "orders", wantErr: qerrors.ErrTopicParse},
		{topic: "nats://", wantErr: qerrors.ErrTopicParse},
		{topic: "://missing", wantErr: qerrors.ErrTopicParse},
	}
	for _, tt := range tests {
		got, name, err := parseTopic(tt.topic)
		if tt.wantErr != nil {
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("parseTopic(%q) err = %v, want %v", tt.topic, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("parseTopic(%q) unexpected err %v", tt.topic, err)
		}
		if got != tt.want || name != tt.name {
			t.Fatalf("parseTopic(%q) = (%d, %q), want (%d, %q)", tt.topic, got, name, tt.want, tt.name)
		}
	}
}

func TestMessageQueueRoutesNatsAndLocal(t *testing.T) {
	var natsTopic, localTopic string
	q := NewMessageQueue(
		stubQueue{publish: func(topic string, _ ...miface.PubOption) error {
			natsTopic = topic
			return nil
		}},
		stubQueue{publish: func(topic string, _ ...miface.PubOption) error {
			localTopic = topic
			return nil
		}},
	)

	if err := q.Publish("nats://orders"); err != nil {
		t.Fatalf("nats publish: %v", err)
	}
	if err := q.Publish("local://inproc"); err != nil {
		t.Fatalf("local publish: %v", err)
	}
	if natsTopic != "orders" || localTopic != "inproc" {
		t.Fatalf("routed topics = nats:%q local:%q", natsTopic, localTopic)
	}
}

func TestMessageQueueMissingBackends(t *testing.T) {
	q := NewMessageQueue(nil, nil)
	if err := q.Publish("nats://orders"); !errors.Is(err, qerrors.ErrNoNatsQueue) {
		t.Fatalf("missing nats err = %v", err)
	}
	if err := q.Publish("local://inproc"); !errors.Is(err, qerrors.ErrNoLocalQueue) {
		t.Fatalf("missing local err = %v", err)
	}
	if _, err := q.Subscribe(context.Background(), "kafka://x", nil); !errors.Is(err, qerrors.ErrMQTypeUnsupported) {
		t.Fatalf("unsupported subscribe err = %v", err)
	}
}
