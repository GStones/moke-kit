package fxsvcapp

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/server/internal/rpc"
	"moke-kit/server/internal/zinx"
	"moke-kit/server/network"
	fxsvcapp "moke-kit/tracing/pkg"
	"moke-kit/tracing/tiface"

	"moke-kit/server/siface"
)

type ServersParams struct {
	fx.In
	TcpGrpcServer    siface.IGrpcServer    `name:"GrpcServer"`
	TcpGatewayServer siface.IGatewayServer `name:"GatewayServer"`
	ZinxServer       siface.IZinxServer    `name:"ZinxServer"`
}

type ServersResult struct {
	fx.Out
	TcpGrpcServer    siface.IGrpcServer    `name:"GrpcServer"`
	TcpGatewayServer siface.IGatewayServer `name:"GatewayServer"`
	ZinxServer       siface.IZinxServer    `name:"ZinxServer"`
}

func (f *ServersResult) Execute(
	l *zap.Logger,
	tr tiface.Tracer,
	s SettingsParams,
	//a GlobalAuthClient,
	mux ConnectionMuxParams,
) (err error) {
	if grpcServer, err := rpc.NewGrpcServer(
		l,
		tr,
		mux.ConnectionMux,
		network.Port(s.Port),
		s.Version,
		//a.AuthClient,
	); err != nil {
		return err
	} else {
		f.TcpGrpcServer = grpcServer
	}

	if zinxServer, err := zinx.NewZinxTcpServer(
		l,
		network.Port(s.Port),
	); err != nil {
		return err
	} else {
		f.ZinxServer = zinxServer
	}

	if gateway, err := rpc.NewTcpGatewayServer(
		l,
		mux.ConnectionMux,
		network.Port(s.Port),
	); err != nil {
		return err
	} else {
		f.TcpGatewayServer = gateway
	}
	return
}

var ServersModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.TracerParams,
		g SettingsParams,
		//a GlobalAuthClient,
		s SecuritySettingsParams,
		m ConnectionMuxParams,
	) (out ServersResult, err error) {
		err = out.Execute(l, t.Tracer, g, m)
		return
	},
)
