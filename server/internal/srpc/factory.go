package srpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/tracing/tiface"
)

func NewGrpcServer(
	logger *zap.Logger,
	tracer tiface.ITracer,
	listener net.Listener,
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
