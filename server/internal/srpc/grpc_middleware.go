package srpc

import (
	"context"
	"fmt"
	"runtime/debug"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gstones/moke-kit/server/internal/common"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

func authFunc(authClient siface.IAuth) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if token, err := auth.AuthFromMD(ctx, string(utility.TokenContextKey)); err != nil {
			return nil, err
		} else if authClient != nil {
			if uid, err := authClient.Auth(token); err != nil {
				return nil, err
			} else {
				return context.WithValue(ctx, utility.UIDContextKey, uid), nil
			}
		}
		return ctx, nil
	}
}

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

func allBut(_ context.Context, _ interceptors.CallMeta) bool {
	return true
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

func addInterceptorOptions(
	logger *zap.Logger,
	authClient siface.IAuth,
	deployments utility.Deployments,
	rateLimit int32,
	opts ...grpc.ServerOption,
) []grpc.ServerOption {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	// Setup metric for panic recoveries.
	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		logger.Error("recovered from panic", zap.Any("panic", p), zap.String("stack", string(debug.Stack())))
		return status.Errorf(codes.Internal, "%s", p)
	}
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}
	rl := common.CreateRateLimiter(int(rateLimit))
	ui := []grpc.UnaryServerInterceptor{
		ratelimit.UnaryServerInterceptor(rl),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc(authClient)), selector.MatchFunc(allBut)),
		srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
		logging.UnaryServerInterceptor(
			interceptorLogger(logger),
			logging.WithDisableLoggingFields(
				logging.ComponentFieldKey,
				logging.MethodTypeFieldKey,
				logging.SystemTag[0],
				"custom-field-should-be-ignored",
			),
			logging.WithLevels(logging.DefaultServerCodeToLevel),
			logging.WithFieldsFromContext(fieldsFromCtx),
			logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
		),
	}
	si := []grpc.StreamServerInterceptor{
		ratelimit.StreamServerInterceptor(rl),
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc(authClient)), selector.MatchFunc(allBut)),
		srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
		logging.StreamServerInterceptor(
			interceptorLogger(logger),
			logging.WithDisableLoggingFields(
				logging.ComponentFieldKey,
				logging.MethodTypeFieldKey,
				logging.SystemTag[0],
				"custom-field-should-be-ignored",
			),
			logging.WithLevels(logging.DefaultServerCodeToLevel),
			logging.WithFieldsFromContext(fieldsFromCtx),
		),
	}

	if deployments.IsProd() {
		ui = append(ui, recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)))
		si = append(si, recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)))
	}

	interceptorOpts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(0),
		grpc.ChainStreamInterceptor(si...),
		grpc.ChainUnaryInterceptor(ui...),
	}
	logger.Info("grpc server interceptor options", zap.Any("options", interceptorOpts))
	if opts == nil {
		opts = interceptorOpts
	} else {
		opts = append(opts, interceptorOpts...)
		// add OpenTelemetry what is OpenTelemetry? https://www.datadoghq.com/knowledge-center/opentelemetry/
		opt := grpc.StatsHandler(otelgrpc.NewServerHandler())
		opts = append(opts, opt)
	}

	return opts
}
