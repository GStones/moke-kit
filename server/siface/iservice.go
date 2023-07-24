package siface

import (
	"github.com/aceld/zinx/ziface"
	"google.golang.org/grpc"
	"net/http"
)

type IGatewayService interface {
	RegisterWithGatewayServer(server HasGatewayServer) error
}

type IZinxTcpService interface {
	RegisterWithTCPServer(server HasZinxTCPServer)
}

type IGrpcService interface {
	RegisterWithGrpcServer(server HasGrpcServer) error
}

type HasZinxTCPServer interface {
	ZinxTcpServer() ziface.IServer
}

type HasGrpcServer interface {
	GrpcServer() *grpc.Server
}

type HasGatewayServer interface {
	GatewayServer() *http.Server
}
