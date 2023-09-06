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

	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
	ZinxTcpPort   int32                 `name:"ZinxTcpPort"`
	ZinxWSPort    int32                 `name:"ZinxWSPort"`

	GrpcServices    []siface.IGrpcService    `group:"GrpcService"`
	GatewayServices []siface.IGatewayService `group:"GatewayService"`
	ZinxServices    []siface.IZinxService    `group:"ZinxService"`
}

func (g *ServiceBinder) Execute(l *zap.Logger, lc fx.Lifecycle) error {
	connectionMuxHook(lc, g.ConnectionMux)
	if hooks, err := bind(
		l,
		g.bindGrpcServices,
		g.bindGatewayServices,
		g.bindZinxServices,
	); err != nil {
		return err
	} else {

		for _, h := range hooks {
			h(lc)
		}
	}
	return nil
}

func (g *ServiceBinder) bindGrpcServices(
	l *zap.Logger,
) (hooks []LifecycleHook, err error) {
	if len(g.GrpcServices) == 0 {
		return nil, nil
	}
	if grpcServer, err := srpc.NewGrpcServer(
		l,
		nil,
		g.ConnectionMux,
	); err != nil {
		return nil, err
	} else {
		for _, s := range g.GrpcServices {
			l.Info("register grpc service", zap.String("service", reflect.TypeOf(s).String()))
			if err := s.RegisterWithGrpcServer(grpcServer); err != nil {
				return nil, err
			}
		}
		hooks = append(hooks, makeServerHook(grpcServer))
	}
	return
}

func (g *ServiceBinder) bindZinxServices(
	l *zap.Logger,
) (hooks []LifecycleHook, err error) {
	if len(g.ZinxServices) == 0 {
		return nil, nil
	}
	if zinxServer, err := zinx.NewZinxServer(
		l,
		g.ZinxTcpPort,
		g.ZinxWSPort,
	); err != nil {
		return nil, err
	} else {
		for _, s := range g.ZinxServices {
			l.Info("register zinx service", zap.String("service", reflect.TypeOf(s).String()))
			s.RegisterWithServer(zinxServer)
		}
		hooks = append(hooks, makeServerHook(zinxServer))
	}

	return
}

func (g *ServiceBinder) bindGatewayServices(
	l *zap.Logger,
) (hooks []LifecycleHook, err error) {
	if len(g.GatewayServices) == 0 {
		return nil, nil
	}

	if gatewayServer, err := srpc.NewGatewayServer(
		l,
		g.ConnectionMux,
	); err != nil {
		return nil, err
	} else {
		for _, s := range g.GatewayServices {
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
