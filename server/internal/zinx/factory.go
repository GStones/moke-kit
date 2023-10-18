package zinx

import (
	"errors"

	"github.com/gstones/zinx/zconf"
	"github.com/gstones/zinx/znet"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/zinx/interceptors"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
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
	deployment string,
) (result siface.IZinxServer, err error) {
	deploy := utility.ParseDeployments(deployment)
	zconf.GlobalObject.Name = name
	zconf.GlobalObject.Version = version
	zconf.GlobalObject.LogIsolationLevel = 2
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
	if deploy == utility.DeploymentsProd {
		s.AddInterceptor(interceptors.NewRecoverInterceptor(logger.With(zap.String("service", name))))
	} else {
		s.AddInterceptor(interceptors.NewLoggerInterceptor(logger.With(zap.String("service", name))))
	}
	result = &ZinxServer{
		logger: logger,
		server: s,
	}
	return
}
