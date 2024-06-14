package srpc

import (
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gstones/moke-kit/server/middlewares"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

// NewGrpcServer creates a new grpc server.
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

// NewGatewayServer creates a new gateway server.
func NewGatewayServer(
	logger *zap.Logger,
	listener net.Listener,
) (result *GatewayServer, err error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(matcher),
		runtime.WithOutgoingHeaderMatcher(matcher),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: allowCORS(withLogger(mux)),
	}
	result = &GatewayServer{
		logger:   logger,
		server:   server,
		mux:      mux,
		opts:     opts,
		listener: listener,
	}
	return
}
