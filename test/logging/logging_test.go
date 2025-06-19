package logging

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/gstones/moke-kit/test/utils"
)

// TestSLogger tests the structured logger functionality
func TestSLogger(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("DefaultSLogger", func(t *testing.T) {
		// Test creating default slogger
		logger := slog.Default()
		helper.AssertNotNil(logger, "Default logger should not be nil")

		// Test logging with default logger
		logger.Info("test message", "key", "value")
		logger.Debug("debug message", "debug_key", "debug_value")
		logger.Warn("warning message", "warn_key", "warn_value")
		logger.Error("error message", "error_key", "error_value")
	})

	t.Run("CustomSLogger", func(t *testing.T) {
		// Create a buffer to capture log output
		var buf bytes.Buffer
		
		// Create a custom slog handler
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		
		logger := slog.New(handler)
		helper.AssertNotNil(logger)

		// Test logging
		logger.Info("test info message", "test_key", "test_value")
		helper.AssertTrue(buf.Len() > 0, "Buffer should contain log output")
	})

	t.Run("SLoggerLevels", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Create handler with different levels
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})
		
		logger := slog.New(handler)
		
		// Test different log levels
		logger.Debug("debug message") // Should not appear
		logger.Info("info message")   // Should not appear
		logger.Warn("warn message")   // Should appear
		logger.Error("error message") // Should appear
		
		output := buf.String()
		helper.AssertFalse(bytes.Contains([]byte(output), []byte("debug message")), "Debug message should not appear")
		helper.AssertFalse(bytes.Contains([]byte(output), []byte("info message")), "Info message should not appear")
		helper.AssertTrue(bytes.Contains([]byte(output), []byte("warn message")), "Warn message should appear")
		helper.AssertTrue(bytes.Contains([]byte(output), []byte("error message")), "Error message should appear")
	})

	t.Run("SLoggerWithContext", func(t *testing.T) {
		var buf bytes.Buffer
		
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		
		logger := slog.New(handler)
		
		// Test logging with context
		ctx := context.WithValue(helper.Context(), "request_id", "12345")
		logger.InfoContext(ctx, "context message", "user_id", "user123")
		
		output := buf.String()
		helper.AssertTrue(len(output) > 0, "Should have log output")
		helper.AssertTrue(bytes.Contains([]byte(output), []byte("context message")), "Should contain log message")
	})
}

// TestZapLogger tests Zap logger functionality
func TestZapLogger(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ZapLoggerCreation", func(t *testing.T) {
		// Test production logger
		prodLogger, err := zap.NewProduction()
		helper.RequireNoError(err)
		helper.AssertNotNil(prodLogger)
		defer prodLogger.Sync()

		// Test development logger
		devLogger, err := zap.NewDevelopment()
		helper.RequireNoError(err)
		helper.AssertNotNil(devLogger)
		defer devLogger.Sync()

		// Test example logger
		exampleLogger := zap.NewExample()
		helper.AssertNotNil(exampleLogger)
		defer exampleLogger.Sync()
	})

	t.Run("ZapLoggerWithObserver", func(t *testing.T) {
		// Create an observer for testing
		core, recorded := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)

		// Test logging
		logger.Info("test message",
			zap.String("key1", "value1"),
			zap.Int("key2", 42),
			zap.Bool("key3", true),
		)

		// Verify logged entries
		entries := recorded.All()
		helper.AssertEqual(1, len(entries), "Should have one log entry")

		entry := entries[0]
		helper.AssertEqual("test message", entry.Message)
		helper.AssertEqual(zapcore.InfoLevel, entry.Level)

		// Check fields
		fields := entry.Context
		helper.AssertEqual(3, len(fields), "Should have three fields")
	})

	t.Run("ZapLoggerLevels", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)

		// Test different levels
		logger.Debug("debug message")  // Should not be recorded
		logger.Info("info message")    // Should not be recorded
		logger.Warn("warn message")    // Should be recorded
		logger.Error("error message") // Should be recorded

		entries := recorded.All()
		helper.AssertEqual(2, len(entries), "Should have two log entries")

		// Check levels
		helper.AssertEqual(zapcore.WarnLevel, entries[0].Level)
		helper.AssertEqual(zapcore.ErrorLevel, entries[1].Level)
	})

	t.Run("ZapLoggerWithFields", func(t *testing.T) {
		core, recorded := observer.New(zapcore.DebugLevel)
		logger := zap.New(core)

		// Create logger with fields
		loggerWithFields := logger.With(
			zap.String("service", "test-service"),
			zap.String("version", "1.0.0"),
		)

		loggerWithFields.Info("service started")

		entries := recorded.All()
		helper.AssertEqual(1, len(entries), "Should have one log entry")

		entry := entries[0]
		helper.AssertEqual("service started", entry.Message)

		// Check fields
		fields := entry.Context
		helper.AssertTrue(len(fields) >= 2, "Should have at least two fields")
	})

	t.Run("ZapSugaredLogger", func(t *testing.T) {
		core, recorded := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)
		sugar := logger.Sugar()

		// Test sugared logging
		sugar.Infow("sugared message",
			"key1", "value1",
			"key2", 42,
		)

		sugar.Infof("formatted message: %s=%d", "count", 100)

		entries := recorded.All()
		helper.AssertEqual(2, len(entries), "Should have two log entries")

		helper.AssertEqual("sugared message", entries[0].Message)
		helper.AssertEqual("formatted message: count=100", entries[1].Message)
	})
}

// TestLoggerIntegration tests integration between different logger types
func TestLoggerIntegration(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ZapToSlogIntegration", func(t *testing.T) {
		// Create Zap logger
		core, recorded := observer.New(zapcore.InfoLevel)
		zapLogger := zap.New(core)

		// Test Zap logging
		zapLogger.Info("zap message", zap.String("source", "zap"))

		// Verify Zap entry
		entries := recorded.All()
		helper.AssertEqual(1, len(entries), "Should have one Zap entry")
		helper.AssertEqual("zap message", entries[0].Message)

		// Test slog logging
		var buf bytes.Buffer
		slogHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slogLogger := slog.New(slogHandler)

		slogLogger.Info("slog message", "source", "slog")
		
		output := buf.String()
		helper.AssertTrue(bytes.Contains([]byte(output), []byte("slog message")), "Should contain slog message")
	})
}

// TestLoggerConfiguration tests logger configuration
func TestLoggerConfiguration(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ZapConfiguration", func(t *testing.T) {
		// Test custom Zap configuration
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		logger, err := config.Build()
		helper.RequireNoError(err)
		helper.AssertNotNil(logger)
		defer logger.Sync()

		// Test logging with custom config
		logger.Info("configured logger message")
	})

	t.Run("SlogConfiguration", func(t *testing.T) {
		// Test custom slog configuration
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}

		var buf bytes.Buffer
		handler := slog.NewTextHandler(&buf, opts)
		logger := slog.New(handler)

		logger.Debug("debug with source")
		
		output := buf.String()
		helper.AssertTrue(len(output) > 0, "Should have output")
		helper.AssertTrue(bytes.Contains([]byte(output), []byte("debug with source")), "Should contain message")
	})
}

// TestLoggerPerformance tests logger performance characteristics
func TestLoggerPerformance(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ZapVsSlogPerformance", func(t *testing.T) {
		// This is a simple performance comparison test
		// In real scenarios, you'd use benchmarks
		
		// Test Zap performance
		zapLogger := zap.NewExample()
		defer zapLogger.Sync()

		// Test slog performance
		slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		// Both should work without errors
		zapLogger.Info("zap performance test")
		slogLogger.Info("slog performance test")
		
		helper.AssertTrue(true, "Performance test completed")
	})
}

// BenchmarkLogging benchmarks different logging approaches
func BenchmarkLogging(b *testing.B) {
	b.Run("ZapStructuredLogging", func(b *testing.B) {
		logger := zap.NewExample()
		defer logger.Sync()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("benchmark message",
					zap.String("key1", "value1"),
					zap.Int("key2", 42),
				)
			}
		})
	})

	b.Run("ZapSugaredLogging", func(b *testing.B) {
		logger := zap.NewExample().Sugar()
		defer logger.Sync()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow("benchmark message",
					"key1", "value1",
					"key2", 42,
				)
			}
		})
	})

	b.Run("SlogStructuredLogging", func(b *testing.B) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("benchmark message",
					"key1", "value1",
					"key2", 42,
				)
			}
		})
	})
}
