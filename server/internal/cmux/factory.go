package cmux

import (
	"go.uber.org/zap"
)

// NewConnectionMux creates a new connection mux.
func NewConnectionMux(
	logger *zap.Logger,
	port int32,
) (result *ConnectionMux, err error) {
	result = &ConnectionMux{
		logger: logger,
		port:   port,
	}
	return
}

// NewTlsConnectionMux creates a new connection mux with TLS.
func NewTlsConnectionMux(
	logger *zap.Logger,
	port int32,
	tlsCert string,
	tlsKey string,
	clientsCA string,
) (result *ConnectionMux, err error) {
	if config, e := makeTLSConfig(logger, tlsCert, tlsKey, clientsCA); e != nil {
		err = e
	} else {
		result = &ConnectionMux{
			logger:    logger,
			port:      port,
			tlsConfig: config,
		}
	}

	return
}
