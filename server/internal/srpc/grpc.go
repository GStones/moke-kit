package srpc

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GrpcServer is the struct for the grpc server.
type GrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener net.Listener
}

// StartServing starts the grpc server.
func (gs *GrpcServer) StartServing(_ context.Context) error {
	gs.logger.Info(
		"grpc start serving",
		zap.String("network", gs.listener.Addr().Network()),
		zap.String("address", gs.listener.Addr().String()),
	)

	go func() {
		if err := gs.server.Serve(gs.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			gs.logger.Error(
				"failed to serve grpc",
				zap.String("network", gs.listener.Addr().Network()),
				zap.String("address", gs.listener.Addr().String()),
				zap.Error(err),
			)
		}
	}()
	return nil
}

// StopServing stops the grpc server.
func (gs *GrpcServer) StopServing(_ context.Context) error {
	gs.server.GracefulStop()
	return nil
}

// GrpcServer returns the grpc server.
func (gs *GrpcServer) GrpcServer() *grpc.Server {
	return gs.server
}
