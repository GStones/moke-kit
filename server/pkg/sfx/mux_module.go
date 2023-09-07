package sfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/cmux"
	"github.com/gstones/moke-kit/server/siface"
)

type ConnectionMuxParams struct {
	fx.In
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
}

type ConnectionMuxResult struct {
	fx.Out
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
}

func (cmr *ConnectionMuxResult) Execute(
	l *zap.Logger,
	g SettingsParams,
	s SecuritySettingsParams,
) (err error) {
	if s.TlsCert != "" && s.TlsKey != "" {
		cmr.ConnectionMux, err = cmux.NewTlsConnectionMux(l, g.Port, s.TlsCert, s.TlsKey)
	} else {
		cmr.ConnectionMux, err = cmux.NewConnectionMux(l, g.Port)
	}
	return
}

var ConnectionMuxModule = fx.Provide(
	func(l *zap.Logger, g SettingsParams, s SecuritySettingsParams) (out ConnectionMuxResult, err error) {
		err = out.Execute(l, g, s)
		return
	},
)
