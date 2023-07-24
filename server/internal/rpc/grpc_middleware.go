package rpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"moke-kit/tracing/tiface"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
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
	// TODO read fields from context
	return logging.Fields{}
}

func addInterceptorOptions(
	zapLogger *zap.Logger,
	tracer tiface.Tracer,
	version string,
	//authClient cli.AuthClient,
	opts ...grpc.ServerOption,
) []grpc.ServerOption {
	logger := log.NewLogfmtLogger(os.Stderr)
	si := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
		logging.StreamServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(fieldsFromCtx)),
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc()), selector.MatchFunc(allButLogin)),
	}

	ui := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
		logging.UnaryServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(fieldsFromCtx)),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc()), selector.MatchFunc(allButLogin)),
	}

	interceptorOpts := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(si...),
		grpc.ChainUnaryInterceptor(ui...),
	}

	if opts == nil {
		opts = interceptorOpts
	} else {
		opts = append(opts, interceptorOpts...)
	}

	return opts
}
