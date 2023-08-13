package cmux

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"
)

func NewTestConnectionMux() (result *TestConnectionMux, err error) {
	result = &TestConnectionMux{
		listener: bufconn.Listen(256 * 1024),
	}
	return
}

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

func NewTlsConnectionMux(
	logger *zap.Logger,
	port int32,
	tlsCert string,
	tlsKey string,
) (result *ConnectionMux, err error) {
	if config, e := makeTlsConfig(tlsCert, tlsKey); e != nil {
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
