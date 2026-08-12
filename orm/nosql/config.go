package nosql

import (
	"errors"
	"time"

	"github.com/gstones/moke-kit/mq/miface"
)

// WriteBackConfig configures write-back workers.
type WriteBackConfig struct {
	// Enabled toggles write-back.
	Enabled bool `json:"enabled" yaml:"enabled" envconfig:"WRITEBACK_ENABLED" default:"false"`
	// Delay is the publish delay before write-back.
	Delay time.Duration `json:"delay" yaml:"delay" envconfig:"WRITEBACK_DELAY" default:"500ms"`
	// BatchSize is reserved for future batching.
	BatchSize int `json:"batch_size" yaml:"batch_size" envconfig:"WRITEBACK_BATCH_SIZE" default:"100"`
	// MaxRetries is reserved for future retry policy.
	MaxRetries int `json:"max_retries" yaml:"max_retries" envconfig:"WRITEBACK_MAX_RETRIES" default:"3"`
	// RetryDelay is reserved for future retry policy.
	RetryDelay time.Duration `json:"retry_delay" yaml:"retry_delay" envconfig:"WRITEBACK_RETRY_DELAY" default:"1s"`
	// WorkerCount is the number of write-back subscribers.
	WorkerCount int `json:"worker_count" yaml:"worker_count" envconfig:"WRITEBACK_WORKER_COUNT" default:"1"`
	// QueueSize is reserved for future buffering.
	QueueSize int `json:"queue_size" yaml:"queue_size" envconfig:"WRITEBACK_QUEUE_SIZE" default:"1000"`
}

// Validate validates configuration.
func (c WriteBackConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
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
	return nil
}

// ToWriteBackOptions converts config into document write-back options.
func (c WriteBackConfig) ToWriteBackOptions(mq miface.MessageQueue) WriteBackOptions {
	return WriteBackOptions{
		Enabled: c.Enabled,
		Delay:   c.Delay,
		MQ:      mq,
	}
}

// DefaultWriteBackConfig returns the default write-back configuration.
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
