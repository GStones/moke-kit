package siface

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/gstones/zinx/ziface"
)

// IServer is the interface for the server, which is used to start and stop the server.
// Implement this interface to start and stop the server.
type IServer interface {
	StartServing(ctx context.Context) error
	StopServing(ctx context.Context) error
}

// IConnectionMux is the interface for the connection mux, which is used to create the listener for the grpc and http.
// Implement this interface to create the listener for the grpc and http.
type IConnectionMux interface {
	IServer
	GrpcListener() (net.Listener, error)
	HTTPListener() (net.Listener, error)
}

// IGrpcServer is the interface for the grpc server, which is used to create the grpc server.
type IGrpcServer interface {
	IServer
	GrpcServer() *grpc.Server
}

// IGatewayServer is the interface for the gateway server, which is used to create the gateway server.
type IGatewayServer interface {
	IServer
	GatewayRuntimeMux() *runtime.ServeMux
	GatewayOption() []grpc.DialOption
	GatewayServer() *http.Server
	Endpoint() string
}

// IZinxServer is the interface for the zinx server, which is used to create the zinx server.
type IZinxServer interface {
	IServer
	ZinxServer() ziface.IServer
}
