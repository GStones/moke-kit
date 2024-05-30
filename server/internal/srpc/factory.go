package srpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/middlewares"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

func NewGrpcServer(
	logger *zap.Logger,
	listener net.Listener,
	auth siface.IAuthMiddleware,
	deployment string,
	rateLimit int32,
	opts ...grpc.ServerOption,
) (result siface.IGrpcServer, err error) {
	deploy := utility.ParseDeployments(deployment)
	opts = middlewares.MakeServerOptions(logger, auth, deploy, rateLimit, opts...)
	result = &GrpcServer{
		logger:   logger,
		listener: listener,
		server:   grpc.NewServer(opts...),
	}
	return result, nil
}
