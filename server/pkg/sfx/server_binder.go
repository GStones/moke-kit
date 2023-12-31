package sfx

import (
	"context"
	"reflect"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/srpc"
	"github.com/gstones/moke-kit/server/internal/zinx"
	"github.com/gstones/moke-kit/server/siface"
)

type LifecycleHook = func(lc fx.Lifecycle)
type BinderFunc func(*zap.Logger) ([]LifecycleHook, error)

type ServiceBinder struct {
	fx.In
	AppName    string `name:"AppName"`
	AppId      string `name:"AppId"`
	Deployment string `name:"Deployment"`
	Version    string `name:"Version"`
	RateLimit  int32  `name:"RateLimit"`

	Timeout       int32                 `name:"Timeout"`
	AuthService   siface.IAuth          `name:"AuthService" optional:"true"`
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
	ZinxTcpPort   int32                 `name:"ZinxTcpPort"`
	ZinxWSPort    int32                 `name:"ZinxWSPort"`

	GrpcServices    []siface.IGrpcService    `group:"GrpcService"`
	GatewayServices []siface.IGatewayService `group:"GatewayService"`
	ZinxServices    []siface.IZinxService    `group:"ZinxService"`
}

func (sb *ServiceBinder) Execute(l *zap.Logger, lc fx.Lifecycle) error {
	if hooks, err := bind(
		l,
		sb.bindGrpcServices,
		sb.bindGatewayServices,
		sb.bindZinxServices,
	); err != nil {
		return err
	} else {
		connectionMuxHook(lc, sb.ConnectionMux)
		for _, h := range hooks {
			h(lc)
		}
	}
	return nil
}

func (sb *ServiceBinder) bindGrpcServices(l *zap.Logger) (hooks []LifecycleHook, err error) {
	if len(sb.GrpcServices) == 0 {
		return nil, nil
	}
	if listener, e := sb.ConnectionMux.GrpcListener(); e != nil {
		err = e
	} else if grpcServer, e := srpc.NewGrpcServer(
		l,
		listener,
		sb.AuthService,
		sb.Deployment,
		sb.RateLimit,
	); e != nil {
		err = e
	} else {
		for _, s := range sb.GrpcServices {
			l.Info("register grpc service", zap.String("service", reflect.TypeOf(s).String()))
			if e := s.RegisterWithGrpcServer(grpcServer); e != nil {
				err = e
			}
		}
		hooks = append(hooks, makeServerHook(grpcServer))
	}
	return
}

func (sb *ServiceBinder) bindZinxServices(
	l *zap.Logger,
) (hooks []LifecycleHook, err error) {
	if len(sb.ZinxServices) == 0 {
		return nil, nil
	}
	if zinxServer, err := zinx.NewZinxServer(
		l,
		sb.ZinxTcpPort,
		sb.ZinxWSPort,
		sb.AppName,
		sb.Version,
		sb.Deployment,
		sb.Timeout,
		sb.RateLimit,
	); err != nil {
		return nil, err
	} else {
		for _, s := range sb.ZinxServices {
			l.Info("register zinx service", zap.String("service", reflect.TypeOf(s).String()))
			s.RegisterWithServer(zinxServer)
		}
		hooks = append(hooks, makeServerHook(zinxServer))
	}

	return
}

func (sb *ServiceBinder) bindGatewayServices(
	l *zap.Logger,
) (hooks []LifecycleHook, err error) {
	if len(sb.GatewayServices) == 0 {
		return nil, nil
	}
	if hLis, e := sb.ConnectionMux.HTTPListener(); e != nil {
		err = e
	} else if gatewayServer, err := srpc.NewGatewayServer(
		l,
		hLis,
	); err != nil {
		return nil, err
	} else {
		for _, s := range sb.GatewayServices {
			l.Info("register gateway service", zap.String("service", reflect.TypeOf(s).String()))
			err := s.RegisterWithGatewayServer(gatewayServer)
			if err != nil {
				return nil, err
			}
		}
		hooks = append(hooks, makeServerHook(gatewayServer))
	}

	return
}

func makeServerHook(s siface.IServer) LifecycleHook {
	return func(lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.StartServing(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return s.StopServing(ctx)
			},
		})
	}
}

func connectionMuxHook(lc fx.Lifecycle, m siface.IConnectionMux) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.StartServing(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return m.StopServing(ctx)
		},
	})
}

func bind(
	l *zap.Logger,
	fs ...BinderFunc,
) (hs []LifecycleHook, err error) {
	for _, f := range fs {
		if gs, e := f(l); e != nil {
			err = e
			break
		} else {
			hs = append(hs, gs...)
		}
	}
	return
}
