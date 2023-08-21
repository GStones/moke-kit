package sfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/common"
	"github.com/gstones/moke-kit/server/internal/srpc"
	"github.com/gstones/moke-kit/server/internal/zinx"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/tracing/tfx"
	"github.com/gstones/moke-kit/tracing/tiface"
)

type ServersParams struct {
	fx.In
	GrpcServer    siface.IGrpcServer    `name:"GrpcServer"`
	GatewayServer siface.IGatewayServer `name:"GatewayServer"`
	ZinxServer    siface.IZinxServer    `name:"ZinxServer"`
}

type ServersResult struct {
	fx.Out
	GrpcServer    siface.IGrpcServer    `name:"GrpcServer"`
	GatewayServer siface.IGatewayServer `name:"GatewayServer"`
	ZinxServer    siface.IZinxServer    `name:"ZinxServer"`
}

func (f *ServersResult) Execute(
	l *zap.Logger,
	tr tiface.ITracer,
	s SettingsParams,
	//a GlobalAuthClient,
	mux ConnectionMuxParams,
) (err error) {
	mod := common.ServerMod(s.Mod)
	if mod.HasGrpc() || mod.HasHttp() {
		if grpcServer, err := srpc.NewGrpcServer(
			l,
			tr,
			mux.ConnectionMux,
			//a.AuthClient,
		); err != nil {
			return err
		} else {
			f.GrpcServer = grpcServer
		}
	}

	if mod.HasTcp() || mod.HasWebsocket() {
		if zinxServer, err := zinx.NewZinxServer(
			l,
			mod,
			s.ZinxTcpPort,
			s.ZinxWSPort,
		); err != nil {
			return err
		} else {
			f.ZinxServer = zinxServer
		}
	}

	if mod.HasHttp() {
		if gateway, err := srpc.NewGatewayServer(
			l,
			mux.ConnectionMux,
			s.Port,
			s.GatewayHost,
		); err != nil {
			return err
		} else {
			f.GatewayServer = gateway
		}
	}
	return
}

var ServersModule = fx.Provide(
	func(
		l *zap.Logger,
		t tfx.TracerParams,
		g SettingsParams,
		//a GlobalAuthClient,
		m ConnectionMuxParams,
	) (out ServersResult, err error) {
		err = out.Execute(l, t.Tracer, g, m)
		return
	},
)
