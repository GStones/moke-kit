package middlewares

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

// MakeServerOptions creates a set of grpc server options
// options include: rate limit, auth, logging, recovery, and opentelemetry etc.
func MakeServerOptions(
	logger *zap.Logger,
	authClient siface.IAuthMiddleware,
	deployments utility.Deployments,
	rateLimit int32,
	opts ...grpc.ServerOption,
) []grpc.ServerOption {

	rl := CreateRateLimiter(int(rateLimit))
	ui := []grpc.UnaryServerInterceptor{
		ratelimit.UnaryServerInterceptor(rl),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc(authClient)), selector.MatchFunc(allBut)),
		logging.UnaryServerInterceptor(
			interceptorLogger(logger),
			logging.WithLevels(logging.DefaultServerCodeToLevel),
			logging.WithFieldsFromContext(fieldsFromCtx),
			logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
		),
	}
	si := []grpc.StreamServerInterceptor{
		ratelimit.StreamServerInterceptor(rl),
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc(authClient)), selector.MatchFunc(allBut)),
		logging.StreamServerInterceptor(
			interceptorLogger(logger),
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
		// add OpenTelemetry what is OpenTelemetry? https://www.datadoghq.com/knowledge-center/opentelemetry/
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	if opts != nil {
		interceptorOpts = append(interceptorOpts, opts...)
	}
	return interceptorOpts
}

func MakeClientOptions(
	logger *zap.Logger,
) []grpc.DialOption {
	ui := []grpc.UnaryClientInterceptor{
		timeout.UnaryClientInterceptor(2 * time.Second),
		logging.UnaryClientInterceptor(
			interceptorLogger(logger),
			logging.WithLevels(logging.DefaultServerCodeToLevel),
			logging.WithFieldsFromContext(fieldsFromCtx),
			logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
		),
	}

	si := []grpc.StreamClientInterceptor{
		logging.StreamClientInterceptor(
			interceptorLogger(logger),
			logging.WithLevels(logging.DefaultServerCodeToLevel),
			logging.WithFieldsFromContext(fieldsFromCtx),
		),
	}
	interceptorOpts := []grpc.DialOption{
		grpc.WithChainStreamInterceptor(si...),
		grpc.WithChainUnaryInterceptor(ui...),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	return interceptorOpts
}
