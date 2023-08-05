package demo

import (
	"context"
	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/pkg/nfx"

	"go.uber.org/fx"
	"go.uber.org/zap"

	pb "moke-kit/demo/api/gen/demo/api"
	"moke-kit/demo/internal/demo/db"
	"moke-kit/demo/pkg/dfx"
	"moke-kit/server/pkg/sfx"
	"moke-kit/server/siface"
)

type Service struct {
	logger   *zap.Logger
	database db.Database
}

func (s *Service) Hi(ctx context.Context, request *pb.HiRequest) (*pb.HiResponse, error) {
	message := request.GetMessage()
	s.logger.Info("Hi", zap.String("message", message), zap.Any("ctx", ctx))

	if data, err := s.database.LoadOrCreateDemo("19000"); err != nil {
		return nil, err
	} else {
		s.logger.Info("LoadOrCreateDemo", zap.Any("data", data))
		if err := data.Update(func() bool {
			data.SetMessage(message)
			return true
		}); err != nil {
			return nil, err
		}
	}

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
	coll diface.ICollection,
) (result *Service, err error) {
	result = &Service{
		logger:   logger,
		database: db.OpenDatabase(logger, coll),
	}

	return
}

var Module = fx.Provide(
	func(
		l *zap.Logger,
		db dfx.DemoDBParams,
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
	) (out sfx.GrpcServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll); err != nil {
			return out, err
		} else {
			out.GrpcService = s
		}
		return
	},
)

var DemoGatewayModule = fx.Provide(
	func(
		l *zap.Logger,
		db dfx.DemoDBParams,
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
	) (out sfx.GatewayServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll); err != nil {
			return out, err
		} else {
			out.GatewayService = s
		}
		return
	},
)
