package zinx

import (
	"github.com/aceld/zinx/znet"
	"go.uber.org/zap"
	"moke-kit/server/network"
)

func NewZinxTcpServer(
	logger *zap.Logger,
	port network.Port,
) (result *ZinxServer, err error) {
	sio := znet.NewServer()
	result = &ZinxServer{
		logger: logger,
		port:   port,
		server: sio,
	}
	return
}
