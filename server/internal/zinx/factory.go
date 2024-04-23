package zinx

import (
	"errors"

	"github.com/gstones/zinx/zconf"
	"github.com/gstones/zinx/znet"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/zinx/interceptors"
	"github.com/gstones/moke-kit/server/pkg/sfx"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

const (
	TcpServerMod string = "tcp"
	WsServerMod  string = "websocket"
)

func NewZinxServer(
	logger *zap.Logger,
	serverSetting sfx.SettingsParams,
	securitySetting sfx.SecuritySettingsParams,
	name string,
	version string,
	deployment string,
	rateLimit int32,
) (result siface.IZinxServer, err error) {
	deploy := utility.ParseDeployments(deployment)
	zconf.GlobalObject.Name = name
	zconf.GlobalObject.Version = version
	zconf.GlobalObject.LogIsolationLevel = 2
	zconf.GlobalObject.WsPort = int(serverSetting.ZinxWSPort)
	zconf.GlobalObject.TCPPort = int(serverSetting.ZinxTcpPort)
	zconf.GlobalObject.HeartbeatMax = int(serverSetting.Timeout)
	zconf.GlobalObject.WorkerPoolSize = serverSetting.WorkerPoolSize
	zconf.GlobalObject.MaxPacketSize = serverSetting.MaxPacketSize
	zconf.GlobalObject.MaxWorkerTaskLen = serverSetting.MaxWorkerTaskLen
	zconf.GlobalObject.MaxMsgChanLen = serverSetting.MaxMsgChanLen

	if securitySetting.TCPTlsEnable {
		zconf.GlobalObject.CertFile = securitySetting.ServerCert
		zconf.GlobalObject.PrivateKeyFile = securitySetting.ServerKey
	}

	if serverSetting.ZinxTcpPort != 0 && serverSetting.ZinxWSPort != 0 {
		zconf.GlobalObject.Mode = ""
	} else if serverSetting.ZinxWSPort != 0 {
		zconf.GlobalObject.Mode = WsServerMod
	} else if serverSetting.ZinxTcpPort != 0 {
		zconf.GlobalObject.Mode = TcpServerMod
	} else {
		return nil, errors.New("please set wsPort or tcpPort")
	}
	l := logger.With(zap.String("service", name), zap.Any("deployment", deploy))
	s := znet.NewServer()
	s.AddInterceptor(interceptors.NewRateLimitInterceptor(l, rateLimit))
	result = &ZinxServer{
		logger: logger,
		server: s,
	}
	return
}
