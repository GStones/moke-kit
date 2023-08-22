package srpc

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	zapKit "github.com/go-kit/kit/log/zap"
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
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gstones/moke-kit/tracing/tiface"
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

func interceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
			args := append([]any{"msg", msg}, fields...)
			switch lvl {
			case logging.LevelDebug:
				_ = level.Debug(l).Log(args...)
			case logging.LevelInfo:
				_ = level.Info(l).Log(args...)
			case logging.LevelWarn:
				_ = level.Warn(l).Log(args...)
			case logging.LevelError:
				_ = level.Error(l).Log(args...)
			default:
				panic(fmt.Sprintf("unknown level %v", lvl))
			}
		})
}

func allButLogin(_ context.Context, callMeta interceptors.CallMeta) bool {
	return !strings.Contains(callMeta.FullMethod(), "Login")
}

func fieldsFromCtx(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

func addInterceptorOptions(
	zapLogger *zap.Logger,
	tracer tiface.ITracer,
	//authClient cli.AuthClient,
	opts ...grpc.ServerOption,
) []grpc.ServerOption {
	logger := zapKit.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)

	// Setup metrics.
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
		_ = level.Error(logger).Log("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
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
	//TODO: add auth interceptor here
	//https://github.com/grpc-ecosystem/go-grpc-middleware#auth
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

	logger.Log("msg", "grpc server interceptor options", "options", interceptorOpts)

	if opts == nil {
		opts = interceptorOpts
	} else {
		opts = append(opts, interceptorOpts...)
	}
	return opts
}
