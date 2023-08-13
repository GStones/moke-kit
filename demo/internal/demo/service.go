package demo

import (
	"context"
	"github.com/aceld/zinx/ziface"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	pb "moke-kit/demo/api/gen/demo/api"
	"moke-kit/demo/internal/demo/db_nosql"
	"moke-kit/demo/internal/demo/handlers"
	"moke-kit/demo/pkg/dfx"
	"moke-kit/mq/logic"
	"moke-kit/mq/pkg/qfx"
	"moke-kit/orm/nosql/diface"
	"moke-kit/orm/pkg/nfx"
	"moke-kit/server/pkg/sfx"
	"moke-kit/server/siface"
)

type Service struct {
	logger      *zap.Logger
	demoHandler *handlers.Demo
}

// ---------------- grpc ----------------

func (s *Service) Watch(request *pb.WatchRequest, server pb.Demo_WatchServer) error {
	topic := request.GetTopic()
	s.logger.Info("Watch", zap.String("topic", topic))

	if err := s.demoHandler.Watch(
		server.Context(),
		topic,
		func(message string) error {
			if err := server.Send(&pb.WatchResponse{
				Message: message,
			}); err != nil {
				return err
			}
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (s *Service) Hi(ctx context.Context, request *pb.HiRequest) (*pb.HiResponse, error) {
	message := request.GetMessage()
	s.logger.Info("Hi", zap.String("message", message))

	if err := s.demoHandler.Hi(request.GetUid(), request.GetMessage()); err != nil {
		return nil, err
	}
	return &pb.HiResponse{
		Message: "response:  " + message,
	}, nil

}
func (s *Service) RegisterWithGrpcServer(server siface.IGrpcServer) error {
	pb.RegisterDemoServer(server.GrpcServer(), s)
	return nil
}

// ---------------- gateway ----------------

func (s *Service) RegisterWithGatewayServer(server siface.IGatewayServer) error {
	return pb.RegisterDemoHandlerFromEndpoint(
		context.Background(),
		server.GatewayRuntimeMux(),
		server.Endpoint(),
		server.GatewayOption(),
	)
}

//---------------- zinx ----------------

func (s *Service) PreHandle(request ziface.IRequest) {

}

func (s *Service) Handle(request ziface.IRequest) {
	switch request.GetMsgID() {
	case 1:
		req := &pb.HiRequest{}
		if err := proto.Unmarshal(request.GetData(), req); err != nil {
			s.logger.Error("unmarshal request data error", zap.Error(err))
		} else {
			if err := s.demoHandler.Hi(req.GetUid(), req.GetMessage()); err != nil {
				s.logger.Error("Hi error", zap.Error(err))
			}
		}
	case 2:
		req := &pb.WatchRequest{}
		if err := proto.Unmarshal(request.GetData(), req); err != nil {
			s.logger.Error("unmarshal request data error", zap.Error(err))
		} else {
			if err := s.demoHandler.Watch(
				request.GetConnection().Context(),
				req.GetTopic(),
				func(message string) error {
					resp := &pb.WatchResponse{
						Message: message,
					}
					if data, err := proto.Marshal(resp); err != nil {
						return err
					} else if err := request.GetConnection().SendMsg(2, data); err != nil {
						return err
					}
					return nil
				}); err != nil {
				s.logger.Error("Watch error", zap.Error(err))
			}
		}
	}
}

func (s *Service) PostHandle(request ziface.IRequest) {

}

func (s *Service) RegisterWithServer(server siface.IZinxServer) {
	server.ZinxServer().AddRouter(1, s)
	server.ZinxServer().AddRouter(2, s)
}

func NewService(
	logger *zap.Logger,
	coll diface.ICollection,
	mq logic.MessageQueue,
	gdb *gorm.DB,
) (result *Service, err error) {
	handler := handlers.NewDemo(
		logger,
		db_nosql.OpenDatabase(logger, coll),
		mq,
		gdb,
	)

	result = &Service{
		logger:      logger,
		demoHandler: handler,
	}
	return
}

var GrpcModule = fx.Provide(
	func(
		l *zap.Logger,
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
		mqParams qfx.MessageQueueParams,
		gParams nfx.GormParams,

	) (out sfx.GrpcServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll, mqParams.MessageQueue, gParams.GormDB); err != nil {
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
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
		mqParams qfx.MessageQueueParams,
		gParams nfx.GormParams,
	) (out sfx.GatewayServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll, mqParams.MessageQueue, gParams.GormDB); err != nil {
			return out, err
		} else {
			out.GatewayService = s
		}
		return
	},
)

var ZinxModule = fx.Provide(
	func(
		l *zap.Logger,
		dProvider nfx.DocumentStoreParams,
		setting dfx.SettingsParams,
		mqParams qfx.MessageQueueParams,
		gParams nfx.GormParams,
	) (out sfx.ZinxServiceResult, err error) {
		if coll, err := dProvider.DriverProvider.OpenDbDriver(setting.DbName); err != nil {
			return out, err
		} else if s, err := NewService(l, coll, mqParams.MessageQueue, gParams.GormDB); err != nil {
			return out, err
		} else {
			out.ZinxService = s
		}
		return
	},
)
