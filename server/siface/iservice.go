package siface

// IGatewayService is the interface for the gateway service, which is used to register with the gateway server.
// Implement this interface to register the service with the gateway server.
//
// type service struct Service{}
//
//	func (s *Service) RegisterWithGatewayServer(server IGatewayServer) error {
//		return pb.RegisterDemoServiceHandlerFromEndpoint(
//			context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption(),
//		)
//	}
type IGatewayService interface {
	RegisterWithGatewayServer(server IGatewayServer) error
}

// IZinxService is the interface for the zinx service, which is used to register with the zinx server.
// Implement this interface to register the service with the zinx server.
//
// type service struct Service{}
//
//	func (s *Service) RegisterWithServer(server IZinxServer) {
//		server.ZinxServer().AddRouter(1,s.HandlerFunc)
//	}
type IZinxService interface {
	RegisterWithServer(server IZinxServer)
}

// IGrpcService is the interface for the grpc service, which is used to register with the grpc server.
// Implement this interface to register the service with the grpc server.
//
// type service struct Service{}
//
//	func (s *Service) RegisterWithGrpcServer(server IGrpcServer) error {
//		return pb.RegisterDemoServiceServer(server.GrpcServer(), s)
//	}

type IGrpcService interface {
	RegisterWithGrpcServer(server IGrpcServer) error
}
