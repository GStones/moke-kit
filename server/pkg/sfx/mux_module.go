package sfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/server/internal/cmux"
	"moke-kit/server/siface"
)

type ConnectionMuxParams struct {
	fx.In
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
}

type ConnectionMuxResult struct {
	fx.Out
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
}

func (f *ConnectionMuxResult) Execute(
	l *zap.Logger,
	g SettingsParams,
	s SecuritySettingsParams,
) (err error) {
	newConnectionMux := func(out *siface.IConnectionMux) {
		if err != nil {
			return
		}
		if s.TlsCert != "" && s.TlsKey != "" {
			*out, err = cmux.NewTlsConnectionMux(l, g.Port, s.TlsCert, s.TlsKey)
		} else {
			*out, err = cmux.NewConnectionMux(l, g.Port)
		}
	}
	newConnectionMux(&f.ConnectionMux)
	return
}

var ConnectionMuxModule = fx.Provide(
	func(l *zap.Logger, g SettingsParams, s SecuritySettingsParams) (out ConnectionMuxResult, err error) {
		err = out.Execute(l, g, s)
		return
	},
)
