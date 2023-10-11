package srpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

func NewGrpcServer(
	logger *zap.Logger,
	listener net.Listener,
	auth siface.IAuth,
	deployment string,
	opts ...grpc.ServerOption,
) (result siface.IGrpcServer, err error) {
	deploy := utility.ParseDeployments(deployment)
	opts = addInterceptorOptions(logger, auth, deploy, opts...)
	result = &GrpcServer{
		logger:   logger,
		listener: listener,
		server:   grpc.NewServer(opts...),
	}
	return result, nil
}
