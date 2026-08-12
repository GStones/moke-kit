package nosql

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/mock"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

func TestWriteBackOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    WriteBackOptions
		wantErr bool
	}{
		{
			name: "valid disabled options",
			opts: WriteBackOptions{
				Enabled: false,
				Delay:   time.Second,
				MQ:      nil,
			},
			wantErr: false,
		},
		{
			name: "valid enabled options",
			opts: WriteBackOptions{
				Enabled: true,
				Delay:   time.Second,
				MQ:      &capturingMQ{},
			},
			wantErr: false,
		},
		{
			name: "invalid: enabled but no MQ",
			opts: WriteBackOptions{
				Enabled: true,
				Delay:   time.Second,
				MQ:      nil,
			},
			wantErr: true,
		},
		{
			name: "invalid: negative delay",
			opts: WriteBackOptions{
				Enabled: false,
				Delay:   -time.Second,
				MQ:      nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteBackWorker_GetMetrics(t *testing.T) {
	worker := NewWriteBackWorker(
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
	)
	metrics := worker.GetMetrics()
	assert.GreaterOrEqual(t, metrics.ProcessedCount, int64(0))
	assert.GreaterOrEqual(t, metrics.FailedCount, int64(0))
}

func TestWriteBackPayload_JSON(t *testing.T) {
	payload := WriteBackPayload{
		CollectionName: "test_collection",
		Key:            "test_key",
		Data:           json.RawMessage(`{"name":"John","age":30}`),
		Version:        noptions.Version(1),
	}

	data, err := json.Marshal(payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded WriteBackPayload
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, payload.CollectionName, decoded.CollectionName)
	assert.Equal(t, payload.Key, decoded.Key)
	assert.Equal(t, payload.Version, decoded.Version)
}
