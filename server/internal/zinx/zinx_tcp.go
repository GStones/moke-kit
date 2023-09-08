package zinx

import (
	"context"

	"github.com/aceld/zinx/ziface"
	"go.uber.org/zap"
)

type ZinxServer struct {
	logger *zap.Logger
	server ziface.IServer
}

func (zs *ZinxServer) ZinxServer() ziface.IServer {
	return zs.server
}

func (zs *ZinxServer) StartServing(_ context.Context) error {
	go func() {
		zs.server.Serve()
	}()
	return nil
}

func (zs *ZinxServer) StopServing(_ context.Context) error {
	zs.server.Stop()
	return nil
}
