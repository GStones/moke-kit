package nosql

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
)

// WriteBackManager manages multiple write-back workers.
type WriteBackManager struct {
	config     WriteBackConfig
	workers    []*WriteBackWorker
	mqClient   miface.MessageQueue
	dbProvider diface.IDocumentProvider
	cache      diface.ICache
	logger     *zap.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	metrics    WriteBackManagerMetrics
}

// WriteBackManagerMetrics aggregates manager metrics.
type WriteBackManagerMetrics struct {
	TotalProcessed int64         `json:"total_processed"`
	TotalFailed    int64         `json:"total_failed"`
	WorkerCount    int           `json:"worker_count"`
	AverageLatency time.Duration `json:"average_latency"`
	Uptime         time.Duration `json:"uptime"`
	StartTime      time.Time     `json:"start_time"`
}

// NewWriteBackManager creates a write-back manager.
// cache may be nil; when set, workers fence delete/recreate generations via __epoch.
func NewWriteBackManager(
	config WriteBackConfig,
	mqClient miface.MessageQueue,
	dbProvider diface.IDocumentProvider,
	logger *zap.Logger,
	cache diface.ICache,
) (*WriteBackManager, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WriteBackManager{
		config:     config,
		mqClient:   mqClient,
		dbProvider: dbProvider,
		cache:      cache,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		metrics: WriteBackManagerMetrics{
			StartTime: time.Now(),
		},
	}, nil
}

// WithCache attaches a document cache used by workers for epoch fencing.
func (m *WriteBackManager) WithCache(cache diface.ICache) *WriteBackManager {
	m.cache = cache
	return m
}

// Start starts configured workers.
func (m *WriteBackManager) Start() error {
	if !m.config.Enabled {
		m.logger.Info("WriteBack is disabled, skipping start")
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cache == nil {
		return errors.New("WriteBack manager requires a document cache for epoch fencing")
	}

	// NATS Subscribe currently fans out to every subscriber (no queue group),
	// so multiple workers would duplicate write-back applies.
	workerCount := m.config.WorkerCount
	if workerCount > 1 {
		m.logger.Warn("WriteBack WorkerCount > 1 is unsupported without queue groups; clamping to 1",
			zap.Int("configured", workerCount),
		)
		workerCount = 1
	}

	m.logger.Info("Starting WriteBack manager",
		zap.Int("worker_count", workerCount),
		zap.Duration("delay", m.config.Delay),
		zap.Bool("epoch_fence", true),
	)

	for i := 0; i < workerCount; i++ {
		worker := NewWriteBackWorker(
			m.mqClient,
			m.dbProvider,
			m.logger.With(zap.Int("worker_id", i)),
		).WithCache(m.cache)
		m.workers = append(m.workers, worker)
		if err := worker.Start(); err != nil {
			m.logger.Error("Failed to start worker", zap.Int("worker_id", i), zap.Error(err))
			m.stopWorkers()
			return err
		}
	}

	m.metrics.WorkerCount = len(m.workers)
	m.logger.Info("WriteBack manager started successfully", zap.Int("workers", len(m.workers)))
	return nil
}

// Stop stops the manager and workers.
func (m *WriteBackManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Stopping WriteBack manager")
	m.cancel()
	m.stopWorkers()
	m.logger.Info("WriteBack manager stopped")
	return nil
}

func (m *WriteBackManager) stopWorkers() {
	var wg sync.WaitGroup
	for i, worker := range m.workers {
		wg.Add(1)
		go func(id int, w *WriteBackWorker) {
			defer wg.Done()
			w.Stop()
			m.logger.Debug("Stopped worker", zap.Int("worker_id", id))
		}(i, worker)
	}
	wg.Wait()
	m.workers = nil
}

// GetMetrics returns aggregated metrics.
func (m *WriteBackManager) GetMetrics() WriteBackManagerMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := m.metrics
	metrics.Uptime = time.Since(metrics.StartTime)

	var totalProcessed, totalFailed int64
	var totalLatency time.Duration
	workerCount := 0
	for _, worker := range m.workers {
		workerMetrics := worker.GetMetrics()
		totalProcessed += workerMetrics.ProcessedCount
		totalFailed += workerMetrics.FailedCount
		totalLatency += workerMetrics.AverageLatency
		workerCount++
	}
	if workerCount > 0 {
		metrics.AverageLatency = totalLatency / time.Duration(workerCount)
	}
	metrics.TotalProcessed = totalProcessed
	metrics.TotalFailed = totalFailed
	return metrics
}

// IsHealthy reports whether the manager looks healthy.
func (m *WriteBackManager) IsHealthy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.config.Enabled {
		return true
	}
	return len(m.workers) > 0
}
