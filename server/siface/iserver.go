package siface

import (
	"context"
	"github.com/aceld/zinx/ziface"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"moke-kit/server/network"
)

type IServer interface {
	Port() network.Port
	StartServing(ctx context.Context) error
	StopServing(ctx context.Context) error
}

type IConnectionMux interface {
	IServer
	HasGrpcListener
	HasHttpListener
}

type IGrpcServer interface {
	IServer
	GrpcServer() *grpc.Server
	Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error)
}

type IGatewayServer interface {
	IServer
	GatewayRuntimeMux() *runtime.ServeMux
	GatewayOption() []grpc.DialOption
}

type IZinxServer interface {
	IServer
	ZinxTcpServer() ziface.IServer
}
