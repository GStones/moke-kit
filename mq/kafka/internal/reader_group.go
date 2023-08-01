package internal

import (
	"fmt"
	"time"

	k "github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/mq"
)

const (
	deadline = time.Duration(500 * time.Millisecond)
)

type readerGroup struct {
	logger            *zap.Logger
	deliverySemantics mq.DeliverySemantics
	freeReaders       chan *reader
	readerCount       int
	maxReaders        int
	kConfig           k.ReaderConfig
	statsPeriod       StatsPeriod
	deadLetterWriter  *deadLetterWriter
}

func NewReaderGroup(
	logger *zap.Logger,
	deliverySemantics mq.DeliverySemantics,
	groupId string,
	topic string,
	brokers Brokers,
	partitionCount int,
	statsPeriod StatsPeriod,
	deadLetterWriter *deadLetterWriter,
) (*readerGroup, error) {
	config := k.ReaderConfig{
		Brokers:           brokers.Strings(),
		Topic:             topic,
		GroupID:           groupId,
		MinBytes:          0,
		MaxBytes:          10e6, // 10MB
		HeartbeatInterval: time.Duration(3 * time.Second),
		SessionTimeout:    time.Duration(12 * time.Second),
		CommitInterval:    0,
	}

	rg := &readerGroup{
		logger:            logger,
		deliverySemantics: deliverySemantics,
		freeReaders:       make(chan *reader, partitionCount),
		maxReaders:        partitionCount,
		kConfig:           config,
		statsPeriod:       statsPeriod,
	}

	rg.freeReaders <- NewReader(logger, config, rg.freeReaders, statsPeriod, deadLetterWriter)
	rg.readerCount++

	return rg, nil
}

func (rg *readerGroup) Subscribe(
	handler func(msg mq.Message, err error) mq.ConsumptionCode,
	decoder mq.Decoder,
	vPtrFactory mq.ValuePtrFactory,
) (*subscription, error) {
	// We try to grab an available consumer within the deadline

	select {
	case r := <-rg.freeReaders:
		sub := NewSubscription(r)
		// Readers are responsible for returning themselves to the pool upon error or subscription cancel
		r.ReceiveLoop(handler, decoder, vPtrFactory)
		return sub, nil
	case <-time.After(deadline):
		// If no reader was available in time, we try to create a new one
		if rg.readerCount < rg.maxReaders {
			newReader := NewReader(rg.logger, rg.kConfig, rg.freeReaders, rg.statsPeriod, rg.deadLetterWriter)
			rg.readerCount++
			newReader.ReceiveLoop(handler, decoder, vPtrFactory)
			return NewSubscription(newReader), nil
		} else {
			return nil, fmt.Errorf("no readers for designated topic were available within deadline")
		}
	}
}
