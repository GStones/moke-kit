package nosql

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
)

// WriteBackManager 管理多個回寫工作器
type WriteBackManager struct {
	config     WriteBackConfig
	workers    []*WriteBackWorker
	mqClient   miface.MessageQueue
	dbProvider diface.IDocumentProvider
	logger     *zap.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
	metrics    WriteBackManagerMetrics
}

// WriteBackManagerMetrics 管理器指標
type WriteBackManagerMetrics struct {
	TotalProcessed int64         `json:"total_processed"`
	TotalFailed    int64         `json:"total_failed"`
	WorkerCount    int           `json:"worker_count"`
	AverageLatency time.Duration `json:"average_latency"`
	Uptime         time.Duration `json:"uptime"`
	StartTime      time.Time     `json:"start_time"`
}

// NewWriteBackManager 创建新的回写管理器
func NewWriteBackManager(
	config WriteBackConfig,
	mqClient miface.MessageQueue,
	dbProvider diface.IDocumentProvider,
	logger *zap.Logger,
) (*WriteBackManager, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	manager := &WriteBackManager{
		config:     config,
		mqClient:   mqClient,
		dbProvider: dbProvider,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		metrics: WriteBackManagerMetrics{
			StartTime: time.Now(),
		},
	}

	return manager, nil
}

// Start 啟動管理器
func (m *WriteBackManager) Start() error {
	if !m.config.Enabled {
		m.logger.Info("WriteBack is disabled, skipping start")
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Starting WriteBack manager",
		zap.Int("worker_count", m.config.WorkerCount),
		zap.Duration("delay", m.config.Delay),
	)

	// 创建工作器
	for i := 0; i < m.config.WorkerCount; i++ {
		worker := NewWriteBackWorker(m.mqClient, m.dbProvider, m.logger.With(zap.Int("worker_id", i)))
		m.workers = append(m.workers, worker)
		
		if err := worker.Start(); err != nil {
			m.logger.Error("Failed to start worker", zap.Int("worker_id", i), zap.Error(err))
			// 停止已啟動的工作器
			m.stopWorkers()
			return err
		}
	}

	m.metrics.WorkerCount = len(m.workers)
	m.logger.Info("WriteBack manager started successfully", zap.Int("workers", len(m.workers)))
	
	return nil
}

// Stop 停止管理器
func (m *WriteBackManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Stopping WriteBack manager")
	m.cancel()
	
	m.stopWorkers()
	
	m.logger.Info("WriteBack manager stopped")
	return nil
}

// stopWorkers 停止所有工作器
func (m *WriteBackManager) stopWorkers() {
	for i, worker := range m.workers {
		worker.Stop()
		m.logger.Debug("Stopped worker", zap.Int("worker_id", i))
	}
	m.workers = nil
}

// GetMetrics 獲取管理器指標
func (m *WriteBackManager) GetMetrics() WriteBackManagerMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := m.metrics
	metrics.Uptime = time.Since(metrics.StartTime)
	
	// 聚合工作器指標
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

// IsHealthy 检查管理器健康状态
func (m *WriteBackManager) IsHealthy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.config.Enabled {
		return true
	}

	// 检查是否有工作器运行
	if len(m.workers) == 0 {
		return false
	}

	// 可以添加更多健康检查逻辑
	return true
}
