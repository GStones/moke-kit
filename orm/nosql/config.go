package nosql

import (
	"errors"
	"time"

	"github.com/gstones/moke-kit/mq/miface"
)

// WriteBackConfig 回写配置
type WriteBackConfig struct {
	// Enabled 是否启用回写功能
	Enabled bool `json:"enabled" yaml:"enabled" envconfig:"WRITEBACK_ENABLED" default:"false"`
	
	// Delay 回写延迟时间
	Delay time.Duration `json:"delay" yaml:"delay" envconfig:"WRITEBACK_DELAY" default:"500ms"`
	
	// BatchSize 批处理大小
	BatchSize int `json:"batch_size" yaml:"batch_size" envconfig:"WRITEBACK_BATCH_SIZE" default:"100"`
	
	// MaxRetries 最大重试次数
	MaxRetries int `json:"max_retries" yaml:"max_retries" envconfig:"WRITEBACK_MAX_RETRIES" default:"3"`
	
	// RetryDelay 重试延迟
	RetryDelay time.Duration `json:"retry_delay" yaml:"retry_delay" envconfig:"WRITEBACK_RETRY_DELAY" default:"1s"`
	
	// WorkerCount 工作器数量
	WorkerCount int `json:"worker_count" yaml:"worker_count" envconfig:"WRITEBACK_WORKER_COUNT" default:"1"`
	
	// QueueSize 队列大小
	QueueSize int `json:"queue_size" yaml:"queue_size" envconfig:"WRITEBACK_QUEUE_SIZE" default:"1000"`
}

// Validate 验证配置
func (c WriteBackConfig) Validate() error {
	if c.Enabled {
		if c.Delay < 0 {
			return errors.New("WriteBack delay cannot be negative")
		}
		if c.BatchSize <= 0 {
			return errors.New("WriteBack batch size must be positive")
		}
		if c.MaxRetries < 0 {
			return errors.New("WriteBack max retries cannot be negative")
		}
		if c.RetryDelay < 0 {
			return errors.New("WriteBack retry delay cannot be negative")
		}
		if c.WorkerCount <= 0 {
			return errors.New("WriteBack worker count must be positive")
		}
		if c.QueueSize <= 0 {
			return errors.New("WriteBack queue size must be positive")
		}
	}
	return nil
}

// ToWriteBackOptions 转换为回写选项
func (c WriteBackConfig) ToWriteBackOptions(mq miface.MessageQueue) WriteBackOptions {
	return WriteBackOptions{
		Enabled: c.Enabled,
		Delay:   c.Delay,
		MQ:      mq,
	}
}

// DefaultWriteBackConfig 返回默認配置
func DefaultWriteBackConfig() WriteBackConfig {
	return WriteBackConfig{
		Enabled:     false,
		Delay:       DefaultWriteBackDelay,
		BatchSize:   100,
		MaxRetries:  3,
		RetryDelay:  time.Second,
		WorkerCount: 1,
		QueueSize:   1000,
	}
}
