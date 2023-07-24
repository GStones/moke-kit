package rpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"moke-kit/server/network"
	"moke-kit/server/siface"
	"moke-kit/tracing/tiface"
)

func NewGrpcServer(
	logger *zap.Logger,
	tracer tiface.Tracer,
	listener siface.HasGrpcListener,
	port network.Port,
	version string,
	//authClient auth.AuthClient,
	opts ...grpc.ServerOption,
) (result siface.IGrpcServer, err error) {
	opts = addInterceptorOptions(logger, tracer, version, opts...)
	result = &GrpcServer{
		logger:   logger,
		listener: listener,
		port:     port,
		server:   grpc.NewServer(opts...),
	}
	return result, nil
}

func NewTestGrpcServer(
	logger *zap.Logger,
	port network.Port,
	opts ...grpc.ServerOption,
) *TestGrpcServer {
	return &TestGrpcServer{
		logger:   logger,
		listener: bufconn.Listen(256 * 1024),
		port:     port,
		server:   grpc.NewServer(opts...),
	}
}
