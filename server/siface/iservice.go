package siface

type IGatewayService interface {
	RegisterWithGatewayServer(server IGatewayServer) error
}

type IZinxService interface {
	RegisterWithServer(server IZinxServer)
}

type IGrpcService interface {
	RegisterWithGrpcServer(server IGrpcServer) error
}
