package logger

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

// ZapLoggerAdapter zap logger adapter
type ZapLoggerAdapter struct {
	logger *zap.Logger
}

// NewZapLoggerAdapter new zap logger adapter
func NewZapLoggerAdapter(logger *zap.Logger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{logger: logger}
}

// Error  zap logger adapter error
func (z *ZapLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	f := make([]zap.Field, 0, len(fields)/2)
	for s, i := range fields {
		f = append(f, zap.Any(s, i))
	}
	z.logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, append(f, zap.Error(err))...)
}

// Info zap logger adapter info
func (z *ZapLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	f := make([]zap.Field, 0, len(fields)/2)
	for s, i := range fields {
		f = append(f, zap.Any(s, i))
	}
	z.logger.WithOptions(zap.AddCallerSkip(1)).With(f...).Info(msg)
}

// Debug zap logger adapter debug
func (z *ZapLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	f := make([]zap.Field, 0, len(fields)/2)
	for s, i := range fields {
		f = append(f, zap.Any(s, i))
	}
	z.logger.WithOptions(zap.AddCallerSkip(1)).With(f...).Debug(msg)
}

// Trace zap logger
func (z *ZapLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	f := make([]zap.Field, 0, len(fields)/2)
	for s, i := range fields {
		f = append(f, zap.Any(s, i))
	}
	z.logger.WithOptions(zap.AddCallerSkip(1)).With(f...).Debug(msg)
}

// With zap logger adapter with
func (z *ZapLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	f := make([]zap.Field, 0, len(fields)/2)
	for s, i := range fields {
		f = append(f, zap.Any(s, i))
	}
	return &ZapLoggerAdapter{logger: z.logger.With(f...)}
}
