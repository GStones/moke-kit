package zinx

import (
	"errors"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/znet"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/zinx/interceptors"
	"github.com/gstones/moke-kit/server/siface"
)

const (
	TcpServerMod string = "tcp"
	WsServerMod  string = "websocket"
)

func NewZinxServer(
	logger *zap.Logger,
	zinxTcpPort int32,
	zinxWsPost int32,
	name string,
	version string,
) (result siface.IZinxServer, err error) {
	zconf.GlobalObject.Name = name
	zconf.GlobalObject.Version = version
	zconf.GlobalObject.LogIsolationLevel = 3
	zconf.GlobalObject.WsPort = int(zinxWsPost)
	zconf.GlobalObject.TCPPort = int(zinxTcpPort)
	if zinxTcpPort != 0 && zinxWsPost != 0 {
		zconf.GlobalObject.Mode = ""
	} else if zinxWsPost != 0 {
		zconf.GlobalObject.Mode = WsServerMod
	} else if zinxTcpPort != 0 {
		zconf.GlobalObject.Mode = TcpServerMod
	} else {
		return nil, errors.New("please set wsPort or tcpPort")
	}
	s := znet.NewServer()
	s.AddInterceptor(interceptors.NewLoggerInterceptor(logger.With(zap.String("service", name))))
	//sio.AddInterceptor(interceptors.NewRecoverInterceptor(logger.With(zap.String("service", name))))
	result = &ZinxServer{
		logger: logger,
		server: s,
	}
	return
}
