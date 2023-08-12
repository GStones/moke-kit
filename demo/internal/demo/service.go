package demo

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/mq/common"

	pb "moke-kit/demo/api/gen/demo/api"
	"moke-kit/demo/internal/demo/db"
	"moke-kit/demo/pkg/dfx"
	"moke-kit/gorm/nosql/diface"
	"moke-kit/gorm/pkg/nfx"
	"moke-kit/mq/pkg/qfx"
	"moke-kit/mq/qiface"
	"moke-kit/server/pkg/sfx"
	"moke-kit/server/siface"
)

type Service struct {
	logger   *zap.Logger
	database db.Database
	mq       qiface.MessageQueue
}

func (s *Service) Watch(request *pb.WatchRequest, server pb.Hello_WatchServer) error {
	topic := request.GetTopic()
	s.logger.Info("Watch", zap.String("topic", topic))

	// mq subscribe
	sub, err := s.mq.Subscribe(
		common.NatsHeader.CreateTopic(topic),
		func(msg qiface.Message, err error) common.ConsumptionCode {
			if err := server.Send(&pb.WatchResponse{
				Message: string(msg.Data()),
			}); err != nil {
				return common.ConsumeNackPersistentFailure
			}
			return common.ConsumeAck
		})
	if err != nil {
		return err
	}
	<-server.Context().Done()
	if err := sub.Unsubscribe(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Hi(ctx context.Context, request *pb.HiRequest) (*pb.HiResponse, error) {
	message := request.GetMessage()
	s.logger.Info("Hi", zap.String("message", message))

	// database create
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

	// mq publish
	if err := s.mq.Publish(
		common.NatsHeader.CreateTopic("demo"), qiface.WithBytes([]byte(message)),
	); err != nil {
		return nil, err
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
	mq qiface.MessageQueue,
) (result *Service, err error) {
	result = &Service{
		logger:   logger,
		database: db.OpenDatabase(logger, coll),
		mq:       mq,
	}

	return
}

var Module = fx.Provide(
	func(
		l *zap.Logger,
		db dfx.DemoDBParams,
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
		mqParams qfx.MessageQueueParams,
	) (out sfx.GrpcServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll, mqParams.MessageQueue); err != nil {
			return out, err
		} else {
			out.GrpcService = s
		}
		return
	},
)

var GatewayModule = fx.Provide(
	func(
		l *zap.Logger,
		db dfx.DemoDBParams,
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
		mqParams qfx.MessageQueueParams,
	) (out sfx.GatewayServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll, mqParams.MessageQueue); err != nil {
			return out, err
		} else {
			out.GatewayService = s
		}
		return
	},
)
