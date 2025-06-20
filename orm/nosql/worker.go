package nosql

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

// WriteBackWorker 处理延迟回写消息的工作器
type WriteBackWorker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	mqClient   miface.MessageQueue
	dbProvider diface.IDocumentProvider
	logger     *zap.Logger
	wg         sync.WaitGroup
}

// NewWriteBackWorker 创建一个新的回写工作器
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

// WriteBackMetrics 回写工作器的指标统计
type WriteBackMetrics struct {
	ProcessedCount int64         `json:"processed_count"`
	FailedCount    int64         `json:"failed_count"`
	AverageLatency time.Duration `json:"average_latency"`
	LastProcessed  time.Time     `json:"last_processed"`
}

// GetMetrics 获取工作器指标
func (w *WriteBackWorker) GetMetrics() WriteBackMetrics {
	// 这里可以使用原子操作来保证线程安全
	// 暂时简化实现
	return WriteBackMetrics{
		ProcessedCount: 0, // TODO: 实现指标收集
		FailedCount:    0,
		AverageLatency: 0,
		LastProcessed:  time.Now(),
	}
}

// handleError 处理回写错误并返回适当的消费代码
func (w *WriteBackWorker) handleError(err error, payload WriteBackPayload) common.ConsumptionCode {
	w.logger.Error("WriteBack operation failed",
		zap.Error(err),
		zap.String("collection", payload.CollectionName),
		zap.String("key", payload.Key),
	)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		w.logger.Warn("WriteBack timeout, will retry", zap.String("key", payload.Key))
		return common.ConsumeNackTransientFailure
	case strings.Contains(err.Error(), "version mismatch"):
		w.logger.Warn("Version mismatch, dropping message", zap.String("key", payload.Key))
		return common.ConsumeNackPersistentFailure
	case strings.Contains(err.Error(), "connection"):
		w.logger.Warn("Connection error, will retry", zap.String("key", payload.Key))
		return common.ConsumeNackTransientFailure
	default:
		w.logger.Error("Unknown error, will retry", zap.Error(err), zap.String("key", payload.Key))
		return common.ConsumeNackTransientFailure
	}
}

// Start 启动回写工作器
func (w *WriteBackWorker) Start() error {
	handler := func(ctx context.Context, msg miface.Message, err error) common.ConsumptionCode {
		var payload WriteBackPayload
		if err := json.Unmarshal(msg.Data(), &payload); err != nil {
			w.logger.Error("Failed to unmarshal message", zap.Error(err))
			return common.ConsumeNackPersistentFailure
		}
		coll, e := w.dbProvider.OpenDbDriver(payload.CollectionName)
		if e != nil {
			w.logger.Error("Failed to open collection", zap.String("collection", payload.CollectionName), zap.Error(e))
			return common.ConsumeNackPersistentFailure
		}

		startTime := time.Now()
		// 原始JSON数据可以直接用于Set操作
		_, err = coll.Set(
			ctx,
			key.NewKey(payload.Key),
			noptions.WithSource(payload.Data),
			noptions.WithVersion(payload.Version),
		)
		latency := time.Since(startTime)

		if err != nil {
			w.logger.Error("Failed to write back document",
				zap.String("key", payload.Key),
				zap.Any("data", payload.Data),
				zap.Error(err),
				zap.Duration("latency", latency),
			)
			return w.handleError(err, payload)
		}

		w.logger.Debug("Successfully wrote back document",
			zap.String("key", payload.Key),
			zap.String("collection", payload.CollectionName),
			zap.Duration("latency", latency),
		)
		
		
		w.logger.Debug("Successfully wrote back document",
			zap.String("key", payload.Key),
			zap.String("collection", payload.CollectionName),
			zap.Duration("latency", latency),
		)
		
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
		_ = consumer.Unsubscribe()
	}()

	return nil
}

// Stop 停止回写工作器
func (w *WriteBackWorker) Stop() {
	w.cancel()
	// 等待所有工作完成
	c := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(c)
	}()

	// 最多等待5秒
	select {
	case <-c:
		// 正常退出
	case <-time.After(5 * time.Second):
		w.logger.Warn("Timeout waiting for writeback worker to stop")
	}
}
