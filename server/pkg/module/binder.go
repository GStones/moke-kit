package module

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/server/internal/srpc"
	"github.com/gstones/moke-kit/server/internal/zinx"
	"github.com/gstones/moke-kit/server/pkg/sfx"
	"github.com/gstones/moke-kit/server/siface"
)

type LifecycleHook = func(lc fx.Lifecycle)
type BinderFunc func(*zap.Logger) ([]LifecycleHook, error)

// ServiceBinder bind all register services to the server
// Service types: grpc, zinx, gateway and opentelemetry provider
type ServiceBinder struct {
	fx.In
	mfx.AppParams // app settings params

	sfx.SettingsParams         // server settings
	sfx.SecuritySettingsParams // server security settings

	sfx.ConnectionMuxParams  // connection mux params
	sfx.GrpcServiceParams    //all grpc service injected (grpc)
	sfx.ZinxServiceParams    // all zinx service injected (tcp/udp/websocket)
	sfx.GatewayServiceParams // all gateway service injected (http)
	sfx.AuthMiddlewareParams // grpc rpc auth middleware injected
	sfx.OTelProviderParams   // opentelemetry provider injected
}

// Bind all types service to the specified server
func (sb *ServiceBinder) Bind(l *zap.Logger, lc fx.Lifecycle) error {
	if hooks, err := bind(
		l,
		sb.bindGrpcServices,
		sb.bindGatewayServices,
		sb.bindZinxServices,
		sb.otelProvider,
	); err != nil {
		return err
	} else {
		if len(sb.GrpcServices) > 0 || len(sb.GatewayServices) > 0 {
			connectionMuxHook(lc, sb.ConnectionMux)
		}
		for _, h := range hooks {
			h(lc)
		}
	}
	return nil
}

func (sb *ServiceBinder) otelProvider(logger *zap.Logger) ([]LifecycleHook, error) {
	if sb.MetricProvider == nil && sb.TracerProvider == nil {
		return nil, nil
	}
	logger.Info("register opentelemetry provider")
	return []LifecycleHook{
		func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					_ = sb.TracerProvider.Shutdown(ctx)
					return sb.MetricProvider.Shutdown(ctx)
				},
			})
		},
	}, nil
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
		sb.AuthMiddleware,
		sb.Deployment,
		sb.RateLimit,
	); e != nil {
		err = e
	} else {
		for _, s := range sb.GrpcServices {
			if e := s.RegisterWithGrpcServer(grpcServer); e != nil {
				err = e
			}
		}
		for k, v := range grpcServer.GrpcServer().GetServiceInfo() {
			l.Info("register grpc service", zap.String("service", k), zap.Any("info", v.Metadata))
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
		sb.SettingsParams,
		sb.SecuritySettingsParams,
		sb.AppName,
		sb.Version,
		sb.Deployment,
		sb.RateLimit,
	); err != nil {
		return nil, err
	} else {
		for _, s := range sb.ZinxServices {
			s.RegisterWithServer(zinxServer)
		}
		l.Info("register zinx service", zap.String("service", zinxServer.ZinxServer().ServerName()))
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
			err := s.RegisterWithGatewayServer(gatewayServer)
			if err != nil {
				return nil, err
			}
		}
		l.Info("register gateway service", zap.String("service", gatewayServer.Endpoint()))
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
