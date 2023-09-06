package zinx

import (
	"errors"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/znet"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/common"

	"github.com/gstones/moke-kit/server/siface"
)

func NewZinxServer(
	logger *zap.Logger,
	zinxTcpPort int32,
	zinxWsPost int32,
) (result siface.IZinxServer, err error) {
	zconf.GlobalObject.LogIsolationLevel = 3
	zconf.GlobalObject.WsPort = int(zinxWsPost)
	zconf.GlobalObject.TCPPort = int(zinxTcpPort)
	if zinxTcpPort != 0 && zinxWsPost != 0 {
		zconf.GlobalObject.Mode = ""
	} else if zinxWsPost != 0 {
		zconf.GlobalObject.Mode = string(common.WsServerMod)
	} else if zinxTcpPort != 0 {
		zconf.GlobalObject.Mode = string(common.TcpServerMod)
	} else {
		return nil, errors.New("please set wsPort or tcpPort")
	}
	sio := znet.NewServer()
	result = &ZinxServer{
		logger: logger,
		server: sio,
	}
	return
}
