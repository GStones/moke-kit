package rpc

import (
	"context"
	"moke-kit/server/network"
	"moke-kit/server/siface"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	timeoutDuration = 10 * time.Second
)

type GrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener siface.HasGrpcListener
	port     network.Port
}

func (s *GrpcServer) StartServing(_ context.Context) error {
	if listener, err := s.listener.GrpcListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"grpc start serving ",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
			zap.Int("port", s.port.Value()),
		)

		go func() {
			if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
				panic(err)
			}
		}()
	}

	return nil
}

func (s *GrpcServer) StopServing(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *GrpcServer) Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	return grpc.DialContext(ctx, target, opts...)
}

func (s *GrpcServer) GrpcServer() *grpc.Server {
	return s.server
}

func (s *GrpcServer) Port() network.Port {
	return s.port
}

type TestGrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener *bufconn.Listener
	port     network.Port
}

func (s *TestGrpcServer) StartServing(_ context.Context) error {
	s.logger.Info(
		"test grpc start serving",
		zap.String("network", s.listener.Addr().Network()),
		zap.String("address", s.listener.Addr().String()),
		zap.Int("port", s.port.Value()),
	)
	go func() {
		if err := s.server.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()
	return nil
}

func (s *TestGrpcServer) StopServing(_ context.Context) error {
	s.server.Stop()
	return nil
}

func (s *TestGrpcServer) Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDialer(func(name string, duration time.Duration) (net.Conn, error) {
			return s.listener.Dial()
		}),
	)

	return grpc.DialContext(context.Background(), target, opts...)
}

func (s *TestGrpcServer) GrpcServer() *grpc.Server {
	return s.server
}

func (s *TestGrpcServer) Port() network.Port {
	return s.port
}
