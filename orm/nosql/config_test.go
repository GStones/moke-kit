package nosql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/mock"
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

	mq := &capturingMQ{}
	options := config.ToWriteBackOptions(mq)

	assert.Equal(t, config.Enabled, options.Enabled)
	assert.Equal(t, config.Delay, options.Delay)
	assert.Equal(t, mq, options.MQ)
}

func TestNewWriteBackManager(t *testing.T) {
	config := DefaultWriteBackConfig()
	config.Enabled = true

	manager, err := NewWriteBackManager(
		config,
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
		newMemoryHashCache(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, config, manager.config)
	assert.NotNil(t, manager.cache)
}

func TestNewWriteBackManager_InvalidConfig(t *testing.T) {
	config := WriteBackConfig{
		Enabled:     true,
		Delay:       -time.Second,
		BatchSize:   100,
		MaxRetries:  3,
		RetryDelay:  time.Second,
		WorkerCount: 1,
		QueueSize:   1000,
	}

	manager, err := NewWriteBackManager(
		config,
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
		nil,
	)
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func TestWriteBackManager_DisabledStart(t *testing.T) {
	config := DefaultWriteBackConfig()

	manager, err := NewWriteBackManager(
		config,
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
		nil,
	)
	assert.NoError(t, err)

	assert.NoError(t, manager.Start())
	assert.Equal(t, 0, manager.GetMetrics().WorkerCount)
	assert.NoError(t, manager.Stop())
}

func TestWriteBackManager_EnabledRequiresCache(t *testing.T) {
	config := DefaultWriteBackConfig()
	config.Enabled = true

	manager, err := NewWriteBackManager(
		config,
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
		nil,
	)
	assert.NoError(t, err)
	assert.Error(t, manager.Start())
}

func TestWriteBackManager_IsHealthy(t *testing.T) {
	config := DefaultWriteBackConfig()

	manager, err := NewWriteBackManager(
		config,
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
		nil,
	)
	assert.NoError(t, err)
	assert.True(t, manager.IsHealthy())
}
