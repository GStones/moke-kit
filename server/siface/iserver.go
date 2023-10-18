package siface

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/gstones/zinx/ziface"
	"google.golang.org/grpc"
)

type IServer interface {
	StartServing(ctx context.Context) error
	StopServing(ctx context.Context) error
}

type IConnectionMux interface {
	IServer
	IGrpcListener
	IHttpListener
}

type IGrpcServer interface {
	IServer
	GrpcServer() *grpc.Server
}

type IGatewayServer interface {
	IServer
	GatewayRuntimeMux() *runtime.ServeMux
	GatewayOption() []grpc.DialOption
	GatewayServer() *http.Server
	Endpoint() string
}

type IZinxServer interface {
	IServer
	ZinxServer() ziface.IServer
}
