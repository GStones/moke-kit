package zinx

import (
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/znet"
	"go.uber.org/zap"
	"moke-kit/server/internal/common"

	"moke-kit/server/siface"
)

func NewZinxServer(
	logger *zap.Logger,
	mod common.ServerMod,
	zinxTcpPort int32,
	zinxWsPost int32,
) (result siface.IZinxServer, err error) {
	zconf.GlobalObject.WsPort = int(zinxWsPost)
	zconf.GlobalObject.TCPPort = int(zinxTcpPort)
	if mod.IsAll() {
		zconf.GlobalObject.Mode = ""
	} else if mod.HasWebsocket() {
		zconf.GlobalObject.Mode = string(common.WsServerMod)
	} else if mod.HasTcp() {
		zconf.GlobalObject.Mode = string(common.TcpServerMod)
	}
	sio := znet.NewServer()
	result = &ZinxServer{
		logger: logger,
		server: sio,
	}
	return
}
