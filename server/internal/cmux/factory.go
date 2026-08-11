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
// When mTLS is true, clients must present a certificate signed by clientsCA.
func NewTlsConnectionMux(
	logger *zap.Logger,
	port int32,
	tlsCert string,
	tlsKey string,
	clientsCA string,
	mTLS bool,
) (result *ConnectionMux, err error) {
	config, cleanup, e := makeTLSConfig(logger, tlsCert, tlsKey, clientsCA, mTLS)
	if e != nil {
		return nil, e
	}
	return &ConnectionMux{
		logger:     logger,
		port:       port,
		tlsConfig:  config,
		tlsCleanup: cleanup,
	}, nil
}
