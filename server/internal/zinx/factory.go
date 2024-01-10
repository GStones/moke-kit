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
	timeout int32,
	rateLimit int32,
	serverCert string,
	serverKey string,
) (result siface.IZinxServer, err error) {
	deploy := utility.ParseDeployments(deployment)
	zconf.GlobalObject.Name = name
	zconf.GlobalObject.Version = version
	zconf.GlobalObject.LogIsolationLevel = 1
	zconf.GlobalObject.WsPort = int(zinxWsPost)
	zconf.GlobalObject.TCPPort = int(zinxTcpPort)
	zconf.GlobalObject.HeartbeatMax = int(timeout)
	zconf.GlobalObject.CertFile = serverCert
	zconf.GlobalObject.PrivateKeyFile = serverKey

	if zinxTcpPort != 0 && zinxWsPost != 0 {
		zconf.GlobalObject.Mode = ""
	} else if zinxWsPost != 0 {
		zconf.GlobalObject.Mode = WsServerMod
	} else if zinxTcpPort != 0 {
		zconf.GlobalObject.Mode = TcpServerMod
	} else {
		return nil, errors.New("please set wsPort or tcpPort")
	}
	l := logger.With(zap.String("service", name))
	s := znet.NewServer()
	if deploy == utility.DeploymentsProd {
		s.AddInterceptor(interceptors.NewRecoverInterceptor(l))
	}
	s.AddInterceptor(interceptors.NewRateLimitInterceptor(l, rateLimit))
	result = &ZinxServer{
		logger: logger,
		server: s,
	}
	return
}
