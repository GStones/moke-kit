package nosql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWriteBackConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  WriteBackConfig
		wantErr bool
	}{
		{
			name:    "valid disabled config",
			config:  DefaultWriteBackConfig(),
			wantErr: false,
		},
		{
			name: "valid enabled config",
			config: WriteBackConfig{
				Enabled:     true,
				Delay:       time.Second,
				BatchSize:   100,
				MaxRetries:  3,
				RetryDelay:  time.Second,
				WorkerCount: 2,
				QueueSize:   1000,
			},
			wantErr: false,
		},
		{
			name: "invalid: negative delay",
			config: WriteBackConfig{
				Enabled:     true,
				Delay:       -time.Second,
				BatchSize:   100,
				MaxRetries:  3,
				RetryDelay:  time.Second,
				WorkerCount: 1,
				QueueSize:   1000,
			},
			wantErr: true,
		},
		{
			name: "invalid: zero batch size",
			config: WriteBackConfig{
				Enabled:     true,
				Delay:       time.Second,
				BatchSize:   0,
				MaxRetries:  3,
				RetryDelay:  time.Second,
				WorkerCount: 1,
				QueueSize:   1000,
			},
			wantErr: true,
		},
		{
			name: "invalid: negative max retries",
			config: WriteBackConfig{
				Enabled:     true,
				Delay:       time.Second,
				BatchSize:   100,
				MaxRetries:  -1,
				RetryDelay:  time.Second,
				WorkerCount: 1,
				QueueSize:   1000,
			},
			wantErr: true,
		},
		{
			name: "invalid: zero worker count",
			config: WriteBackConfig{
				Enabled:     true,
				Delay:       time.Second,
				BatchSize:   100,
				MaxRetries:  3,
				RetryDelay:  time.Second,
				WorkerCount: 0,
				QueueSize:   1000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteBackConfig_ToWriteBackOptions(t *testing.T) {
	config := WriteBackConfig{
		Enabled: true,
		Delay:   time.Second,
	}
	
	mockMQ := &MockMessageQueue{}
	options := config.ToWriteBackOptions(mockMQ)
	
	assert.Equal(t, config.Enabled, options.Enabled)
	assert.Equal(t, config.Delay, options.Delay)
	assert.Equal(t, mockMQ, options.MQ)
}

func TestNewWriteBackManager(t *testing.T) {
	config := DefaultWriteBackConfig()
	config.Enabled = true
	
	mockMQ := &MockMessageQueue{}
	mockProvider := &MockDocumentProvider{}
	logger := zap.NewNop()
	
	manager, err := NewWriteBackManager(config, mockMQ, mockProvider, logger)
	
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, config, manager.config)
}

func TestNewWriteBackManager_InvalidConfig(t *testing.T) {
	config := WriteBackConfig{
		Enabled:     true,
		Delay:       -time.Second, // 無效配置
		BatchSize:   100,
		MaxRetries:  3,
		RetryDelay:  time.Second,
		WorkerCount: 1,
		QueueSize:   1000,
	}
	
	mockMQ := &MockMessageQueue{}
	mockProvider := &MockDocumentProvider{}
	logger := zap.NewNop()
	
	manager, err := NewWriteBackManager(config, mockMQ, mockProvider, logger)
	
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func TestWriteBackManager_DisabledStart(t *testing.T) {
	config := DefaultWriteBackConfig() // Enabled = false
	
	mockMQ := &MockMessageQueue{}
	mockProvider := &MockDocumentProvider{}
	logger := zap.NewNop()
	
	manager, err := NewWriteBackManager(config, mockMQ, mockProvider, logger)
	assert.NoError(t, err)
	
	err = manager.Start()
	assert.NoError(t, err)
	
	metrics := manager.GetMetrics()
	assert.Equal(t, 0, metrics.WorkerCount)
	
	err = manager.Stop()
	assert.NoError(t, err)
}

func TestWriteBackManager_IsHealthy(t *testing.T) {
	config := DefaultWriteBackConfig()
	
	mockMQ := &MockMessageQueue{}
	mockProvider := &MockDocumentProvider{}
	logger := zap.NewNop()
	
	manager, err := NewWriteBackManager(config, mockMQ, mockProvider, logger)
	assert.NoError(t, err)
	
	// 未啟用時應該是健康的
	assert.True(t, manager.IsHealthy())
}
