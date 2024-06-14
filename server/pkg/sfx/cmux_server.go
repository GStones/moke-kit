package sfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/cmux"
	"github.com/gstones/moke-kit/server/siface"
)

// https://github.com/soheilhy/cmux
// ConnectionMux module:listen and serve same port for different protocols

// ConnectionMuxParams module params for injecting ConnectionMux
type ConnectionMuxParams struct {
	fx.In
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
}

// ConnectionMuxResult module result for exporting ConnectionMux
type ConnectionMuxResult struct {
	fx.Out
	ConnectionMux siface.IConnectionMux `name:"ConnectionMux"`
}

func (cmr *ConnectionMuxResult) init(
	l *zap.Logger,
	g SettingsParams,
	s SecuritySettingsParams,
) error {
	if s.TLSEnable {
		mux, err := cmux.NewTlsConnectionMux(l, g.Port, s.ServerCert, s.ServerKey, s.ClientCaCert)
		if err != nil {
			return err
		}
		cmr.ConnectionMux = mux
	} else {
		mux, err := cmux.NewConnectionMux(l, g.Port)
		if err != nil {
			return err
		}
		cmr.ConnectionMux = mux
	}
	return nil
}

// CreateConnectionMux creates a connection mux for the server
func CreateConnectionMux(
	l *zap.Logger,
	g SettingsParams,
	s SecuritySettingsParams,
) (out ConnectionMuxResult, err error) {
	err = out.init(l, g, s)
	return
}

// ConnectionMuxModule module for ConnectionMux
var ConnectionMuxModule = fx.Provide(
	func(l *zap.Logger, g SettingsParams, s SecuritySettingsParams) (ConnectionMuxResult, error) {
		return CreateConnectionMux(l, g, s)
	},
)
