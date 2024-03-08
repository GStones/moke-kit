package cmux

import (
	"go.uber.org/zap"
)

func NewConnectionMux(
	logger *zap.Logger,
	port int32,
) (result *ConnectionMux, err error) {
	result = &ConnectionMux{
		logger: logger,
		port:   port,
	}
	err = result.init()
	return
}

func NewTlsConnectionMux(
	logger *zap.Logger,
	port int32,
	tlsCert string,
	tlsKey string,
	clientsCA string,
) (result *ConnectionMux, err error) {
	if config, e := makeTLSConfig(tlsCert, tlsKey, clientsCA); e != nil {
		err = e
	} else {
		result = &ConnectionMux{
			logger:    logger,
			port:      port,
			tlsConfig: config,
		}
		err = result.init()
	}

	return
}
