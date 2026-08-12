package nosql

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

// WriteBackWorker consumes delayed write-back messages and persists them.
type WriteBackWorker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	mqClient   miface.MessageQueue
	dbProvider diface.IDocumentProvider
	logger     *zap.Logger
	wg         sync.WaitGroup

	processedCount atomic.Int64
	failedCount    atomic.Int64
	totalLatency   atomic.Int64
	lastProcessed  atomic.Value
}

// NewWriteBackWorker creates a write-back worker.
func NewWriteBackWorker(
	mqClient miface.MessageQueue,
	dbProvider diface.IDocumentProvider,
	logger *zap.Logger,
) *WriteBackWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &WriteBackWorker{
		ctx:        ctx,
		cancel:     cancel,
		mqClient:   mqClient,
		dbProvider: dbProvider,
		logger:     logger,
	}
}

// WriteBackMetrics exposes worker counters.
type WriteBackMetrics struct {
	ProcessedCount int64         `json:"processed_count"`
	FailedCount    int64         `json:"failed_count"`
	AverageLatency time.Duration `json:"average_latency"`
	LastProcessed  time.Time     `json:"last_processed"`
}

// GetMetrics returns worker metrics.
func (w *WriteBackWorker) GetMetrics() WriteBackMetrics {
	processed := w.processedCount.Load()
	failed := w.failedCount.Load()
	totalLatency := w.totalLatency.Load()

	var avgLatency time.Duration
	if processed > 0 {
		avgLatency = time.Duration(totalLatency / processed)
	}

	var lastProcessed time.Time
	if val := w.lastProcessed.Load(); val != nil {
		lastProcessed = val.(time.Time)
	}

	return WriteBackMetrics{
		ProcessedCount: processed,
		FailedCount:    failed,
		AverageLatency: avgLatency,
		LastProcessed:  lastProcessed,
	}
}

func (w *WriteBackWorker) handleError(err error, payload WriteBackPayload) common.ConsumptionCode {
	w.failedCount.Add(1)
	w.logger.Error("WriteBack operation failed",
		zap.Error(err),
		zap.String("collection", payload.CollectionName),
		zap.String("key", payload.Key),
	)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return common.ConsumeNackTransientFailure
	case strings.Contains(err.Error(), "version mismatch") ||
		strings.Contains(err.Error(), "ErrVersionNotMatch"):
		return common.ConsumeNackPersistentFailure
	case strings.Contains(err.Error(), "connection"):
		return common.ConsumeNackTransientFailure
	default:
		return common.ConsumeNackTransientFailure
	}
}

// Start subscribes to the write-back topic.
func (w *WriteBackWorker) Start() error {
	handler := func(msg miface.Message, err error) common.ConsumptionCode {
		if err != nil {
			w.logger.Error("WriteBack subscription error", zap.Error(err))
			return common.ConsumeNackTransientFailure
		}

		var payload WriteBackPayload
		if err := json.Unmarshal(msg.Data(), &payload); err != nil {
			w.logger.Error("Failed to unmarshal message", zap.Error(err))
			return common.ConsumeNackPersistentFailure
		}

		coll, e := w.dbProvider.OpenDbDriver(payload.CollectionName)
		if e != nil {
			w.logger.Error("Failed to open collection",
				zap.String("collection", payload.CollectionName),
				zap.Error(e),
			)
			return common.ConsumeNackPersistentFailure
		}

		var src any
		if err := json.Unmarshal(payload.Data, &src); err != nil {
			w.logger.Error("Failed to decode write-back data", zap.Error(err))
			return common.ConsumeNackPersistentFailure
		}

		dbCtx, dbCancel := context.WithTimeout(w.ctx, 30*time.Second)
		defer dbCancel()

		startTime := time.Now()
		_, setErr := coll.Set(
			dbCtx,
			key.NewKey(payload.Key),
			noptions.WithSource(src),
			noptions.WithVersion(payload.Version),
		)
		latency := time.Since(startTime)
		if setErr != nil {
			w.logger.Error("Failed to write back document",
				zap.String("key", payload.Key),
				zap.Error(setErr),
				zap.Duration("latency", latency),
			)
			return w.handleError(setErr, payload)
		}

		w.processedCount.Add(1)
		w.totalLatency.Add(int64(latency))
		w.lastProcessed.Store(time.Now())
		return common.ConsumeAck
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		consumer, err := w.mqClient.Subscribe(w.ctx, WriteBackTopic, handler)
		if err != nil {
			w.logger.Error("Failed to subscribe to writeback topic", zap.Error(err))
			return
		}
		<-w.ctx.Done()
		if consumer != nil {
			_ = consumer.Unsubscribe()
		}
	}()

	return nil
}

// Stop stops the worker.
func (w *WriteBackWorker) Stop() {
	w.cancel()
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		w.logger.Warn("Timeout waiting for writeback worker to stop")
	}
}
