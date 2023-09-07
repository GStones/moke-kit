package srpc

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
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
)

const (
	TokenContextKey = "bearer"
)

func authFunc() auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if token, err := auth.AuthFromMD(ctx, TokenContextKey); err != nil {
			return nil, err
			//TODO change auth client validate token  perform proper Oauth/OIDC verification!
		} else if token != "test" {
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		} else {
			return ctx, nil
		}
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

func allButLogin(_ context.Context, callMeta interceptors.CallMeta) bool {
	return !strings.Contains(callMeta.FullMethod(), "Auth")
}

func fieldsFromCtx(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

func addInterceptorOptions(
	logger *zap.Logger,
	//authClient cli.AuthClient,
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

	//TODO: add rate limit interceptor here
	//https: //github.com/grpc-ecosystem/go-grpc-middleware#server
	//TODO: add external auth interceptor here
	//https://github.com/grpc/proposal/blob/master/A43-grpc-authorization-api.md
	ui := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
		srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
		logging.UnaryServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(fieldsFromCtx)),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc()), selector.MatchFunc(allButLogin)),
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
	}
	si := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
		srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
		logging.StreamServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(fieldsFromCtx)),
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc()), selector.MatchFunc(allButLogin)),
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
	}

	interceptorOpts := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(si...),
		grpc.ChainUnaryInterceptor(ui...),
	}
	logger.Info("grpc server interceptor options", zap.Any("options", interceptorOpts))
	if opts == nil {
		opts = interceptorOpts
	} else {
		opts = append(opts, interceptorOpts...)
	}
	return opts
}
