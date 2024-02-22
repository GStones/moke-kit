package middlewares

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/utility"
)

// interceptorLogger adapts zap logger to interceptor logger.
func interceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]
			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func fieldsFromCtx(ctx context.Context) logging.Fields {
	fields := logging.Fields{}
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		fields = append(fields, "traceID", span.TraceID().String())
	}
	if v, ok := ctx.Value(utility.UIDContextKey).(string); ok {
		fields = append(fields, utility.UIDContextKey.String(), v)
	}
	if v, ok := ctx.Value(utility.WithOutTag).(bool); ok {
		fields = append(fields, utility.WithOutTag.String(), v)
	}
	return fields
}
