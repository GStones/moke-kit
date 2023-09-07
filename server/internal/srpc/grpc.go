package srpc

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type GrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener net.Listener
}

func (s *GrpcServer) StartServing(_ context.Context) error {
	s.logger.Info(
		"grpc start serving ",
		zap.String("network", s.listener.Addr().Network()),
		zap.String("address", s.listener.Addr().String()),
	)

	go func() {
		if err := s.server.Serve(s.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Error(
				"failed to serve grpc",
				zap.String("network", s.listener.Addr().Network()),
				zap.String("address", s.listener.Addr().String()),
				zap.Error(err),
			)
		}
	}()
	return nil
}

func (s *GrpcServer) StopServing(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *GrpcServer) GrpcServer() *grpc.Server {
	return s.server
}

type TestGrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener *bufconn.Listener
	port     int32
}

func (s *TestGrpcServer) StartServing(_ context.Context) error {
	s.logger.Info(
		"test grpc start serving",
		zap.String("network", s.listener.Addr().Network()),
		zap.String("address", s.listener.Addr().String()),
		zap.Int32("port", s.port),
	)
	go func() {
		if err := s.server.Serve(s.listener); err != nil &&
			!errors.Is(err, grpc.ErrServerStopped) {
			panic(err)
		}
	}()
	return nil
}

func (s *TestGrpcServer) StopServing(_ context.Context) error {
	s.server.Stop()
	return nil
}

func (s *TestGrpcServer) GrpcServer() *grpc.Server {
	return s.server
}
