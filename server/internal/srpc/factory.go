package srpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/tracing/tiface"
)

func NewGrpcServer(
	logger *zap.Logger,
	tracer tiface.ITracer,
	listener siface.IGrpcListener,
	//authClient auth.AuthClient,
	opts ...grpc.ServerOption,
) (result siface.IGrpcServer, err error) {
	opts = addInterceptorOptions(logger, tracer, opts...)
	result = &GrpcServer{
		logger:   logger,
		listener: listener,
		server:   grpc.NewServer(opts...),
	}
	return result, nil
}

func NewTestGrpcServer(
	logger *zap.Logger,
	port int32,
	opts ...grpc.ServerOption,
) *TestGrpcServer {
	return &TestGrpcServer{
		logger:   logger,
		listener: bufconn.Listen(256 * 1024),
		port:     port,
		server:   grpc.NewServer(opts...),
	}
}
