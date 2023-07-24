package fxsvcapp

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/server/internal/cmux"
	"moke-kit/server/network"
	"moke-kit/server/siface"
)

type ConnectionMuxParams struct {
	fx.In
	ConnectionMux siface.IConnectionMux `name:"IConnectionMux"`
}

type ConnectionMuxResult struct {
	fx.Out
	ConnectionMux siface.IConnectionMux `name:"IConnectionMux"`
}

func (f *ConnectionMuxResult) Execute(
	l *zap.Logger,
	g SettingsParams,
	s SecuritySettingsParams,
) (err error) {
	newTcpConnectionMux := func(out *siface.IConnectionMux) {
		if err != nil {
			return
		}
		port := network.Port(g.Port)

		if g.AppTestMode {
			*out, err = cmux.NewTestTcpConnectionMux()
		} else if s.TlsCert != "" && s.TlsKey != "" {
			*out, err = cmux.NewTlsTcpConnectionMux(l, port, s.TlsCert, s.TlsKey)
		} else {
			*out, err = cmux.NewTcpConnectionMux(l, port)
		}
	}
	newTcpConnectionMux(&f.ConnectionMux)
	return
}

var ConnectionMuxModule = fx.Provide(
	func(l *zap.Logger, g SettingsParams, s SecuritySettingsParams) (out ConnectionMuxResult, err error) {
		err = out.Execute(l, g, s)
		return
	},
)
