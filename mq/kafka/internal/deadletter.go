package internal

import (
	"context"
	"time"

	k "github.com/segmentio/kafka-go"

	"github.com/gstones/platform/services/common/mq"
)

const (
	deadLetterTopicName = "deadletter"
)

type deadLetterWriter struct {
	kWriter *k.Writer
}

func NewDeadLetterWriter(
	brokers Brokers,
) (*deadLetterWriter, error) {

	topicConfig := k.TopicConfig{
		Topic:              mq.NamespaceTopic(deadLetterTopicName),
		NumPartitions:      int(deadLetterNumPartitions),
		ReplicationFactor:  int(deadLetterReplicationFactor),
		ReplicaAssignments: nil,
		ConfigEntries:      nil,
	}

	writerConfig := k.WriterConfig{
		Brokers:  brokers.Strings(),
		Topic:    mq.NamespaceTopic(deadLetterTopicName),
		Balancer: &k.LeastBytes{},
	}

	broker := brokers.Random()

	if kConn, err := k.Dial(broker.Network.String(), broker.String()); err != nil {
		return nil, err
	} else if err = kConn.CreateTopics(topicConfig); err != nil {
		return nil, err
	} else {
		return &deadLetterWriter{kWriter: k.NewWriter(writerConfig)}, nil
	}
}

func (w *deadLetterWriter) Write(data []byte) {
	go w.write(data)
}

func (w *deadLetterWriter) write(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	return w.kWriter.WriteMessages(ctx, k.Message{Value: data})
}
