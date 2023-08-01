package internal

import (
	"context"

	k "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type writer struct {
	logger  *StatsLogger
	kWriter *k.Writer
}

func NewWriter(logger *zap.Logger, config k.WriterConfig, statsPeriod StatsPeriod) *writer {
	kWriter := k.NewWriter(config)

	statsLogger := NewStatsLogger(
		logger,
		statsPeriod,
		&WriterStatsProvider{writer: kWriter},
	)

	statsLogger.Start()

	return &writer{logger: statsLogger, kWriter: kWriter}
}

func (w *writer) Publish(topic string, data []byte) error {
	return w.kWriter.WriteMessages(context.Background(), k.Message{Value: data})
}
