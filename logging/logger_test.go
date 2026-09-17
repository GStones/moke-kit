package logging

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	for _, deployment := range []string{"local", "dev", "prod"} {
		logger, err := NewLogger(deployment)
		if err != nil {
			t.Fatalf("NewLogger(%q) error = %v", deployment, err)
		}
		if logger == nil {
			t.Fatalf("NewLogger(%q) returned nil logger", deployment)
		}
		_ = logger.Sync()
	}
}

func TestNewLoggerReturnsZapLogger(t *testing.T) {
	logger, err := NewLogger("local")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := any(logger).(*zap.Logger); !ok {
		t.Fatalf("expected *zap.Logger, got %T", logger)
	}
}
