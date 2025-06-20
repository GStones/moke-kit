package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestHelper provides common utilities for testing
type TestHelper struct {
	t      *testing.T
	logger *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	
	return &TestHelper{
		t:      t,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Logger returns the test logger
func (h *TestHelper) Logger() *zap.Logger {
	return h.logger
}

// Context returns the test context
func (h *TestHelper) Context() context.Context {
	return h.ctx
}

// T returns the testing.T instance
func (h *TestHelper) T() *testing.T {
	return h.t
}

// Cleanup should be called in defer to clean up resources
func (h *TestHelper) Cleanup() {
	if h.cancel != nil {
		h.cancel()
	}
}

// AssertNoError asserts that error is nil
func (h *TestHelper) AssertNoError(err error, msg ...string) {
	if len(msg) > 0 {
		assert.NoError(h.t, err, msg[0])
	} else {
		assert.NoError(h.t, err)
	}
}

// RequireNoError requires that error is nil
func (h *TestHelper) RequireNoError(err error, msg ...string) {
	if len(msg) > 0 {
		require.NoError(h.t, err, msg[0])
	} else {
		require.NoError(h.t, err)
	}
}

// AssertEqual asserts that expected equals actual
func (h *TestHelper) AssertEqual(expected, actual interface{}, msg ...string) {
	if len(msg) > 0 {
		assert.Equal(h.t, expected, actual, msg[0])
	} else {
		assert.Equal(h.t, expected, actual)
	}
}

// RequireEqual requires that expected equals actual
func (h *TestHelper) RequireEqual(expected, actual interface{}, msg ...string) {
	if len(msg) > 0 {
		require.Equal(h.t, expected, actual, msg[0])
	} else {
		require.Equal(h.t, expected, actual)
	}
}

// AssertTrue asserts that value is true
func (h *TestHelper) AssertTrue(value bool, msg ...string) {
	if len(msg) > 0 {
		assert.True(h.t, value, msg[0])
	} else {
		assert.True(h.t, value)
	}
}

// AssertFalse asserts that value is false
func (h *TestHelper) AssertFalse(value bool, msg ...string) {
	if len(msg) > 0 {
		assert.False(h.t, value, msg[0])
	} else {
		assert.False(h.t, value)
	}
}

// AssertNotNil asserts that value is not nil
func (h *TestHelper) AssertNotNil(value interface{}, msg ...string) {
	if len(msg) > 0 {
		assert.NotNil(h.t, value, msg[0])
	} else {
		assert.NotNil(h.t, value)
	}
}

// AssertNil asserts that value is nil
func (h *TestHelper) AssertNil(value interface{}, msg ...string) {
	if len(msg) > 0 {
		assert.Nil(h.t, value, msg[0])
	} else {
		assert.Nil(h.t, value)
	}
}

// WaitForCondition waits for a condition to be true within timeout
func (h *TestHelper) WaitForCondition(condition func() bool, timeout time.Duration, msg string) bool {
	ctx, cancel := context.WithTimeout(h.ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.t.Fatalf("Timeout waiting for condition: %s", msg)
			return false
		case <-ticker.C:
			if condition() {
				return true
			}
		}
	}
}

// CreateTestTopic creates a test topic with timestamp
func (h *TestHelper) CreateTestTopic(prefix string) string {
	return fmt.Sprintf("%s-test-%d", prefix, time.Now().UnixNano())
}

// Sleep sleeps for the given duration
func (h *TestHelper) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
