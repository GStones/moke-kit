package demo

import (
	"context"

	"go.uber.org/zap"

	pb "moke-kit/demo/api/gen/demo/api"
	"moke-kit/server/siface"
)

type Service struct {
	logger *zap.Logger
}

func (s *Service) Hi(_ context.Context, request *pb.HiRequest) (*pb.HiResponse, error) {
	message := request.GetMessage()
	s.logger.Info("Hi", zap.String("message", message))
	return &pb.HiResponse{
		Message: "response:  " + message,
	}, nil
}

func (s *Service) RegisterWithGrpcServer(server siface.IGrpcServer) error {
	pb.RegisterHelloServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server siface.IGatewayServer) error {
	return pb.RegisterHelloHandlerFromEndpoint(
		context.Background(),
		server.GatewayRuntimeMux(),
		server.Endpoint(),
		server.GatewayOption(),
	)
}

func NewService(
	logger *zap.Logger,
) (result *Service, err error) {
	result = &Service{
		logger: logger,
	}
	return
}
