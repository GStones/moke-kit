package srpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/siface"
)

func NewGrpcServer(
	logger *zap.Logger,
	listener net.Listener,
	auth siface.IAuth,
	opts ...grpc.ServerOption,
) (result siface.IGrpcServer, err error) {
	opts = addInterceptorOptions(logger, auth, opts...)
	result = &GrpcServer{
		logger:   logger,
		listener: listener,
		server:   grpc.NewServer(opts...),
	}
	return result, nil
}
