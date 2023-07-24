package cmux

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"
	"moke-kit/server/network"
)

func NewTestTcpConnectionMux() (result *TestTcpConnectionMux, err error) {
	result = &TestTcpConnectionMux{
		listener: bufconn.Listen(256 * 1024),
	}
	return
}

func NewTcpConnectionMux(logger *zap.Logger, port network.Port) (result *ConnectionMux, err error) {
	result = &ConnectionMux{
		logger: logger,
		port:   port,
	}

	return
}

func NewTlsTcpConnectionMux(
	logger *zap.Logger,
	port network.Port,
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
