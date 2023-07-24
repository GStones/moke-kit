package zinx

import (
	"context"
	"github.com/aceld/zinx/ziface"
	"go.uber.org/zap"
	"moke-kit/server/network"
)

type ZinxServer struct {
	logger *zap.Logger
	port   network.Port
	server ziface.IServer
}

func (s *ZinxServer) StartServing(_ context.Context) error {
	go func() {
		s.server.Serve()
	}()
	return nil
}

func (s *ZinxServer) StopServing(_ context.Context) error {
	s.server.Stop()
	return nil
}

func (s *ZinxServer) ZinxTcpServer() ziface.IServer {
	return s.server
}

func (s *ZinxServer) Port() network.Port {
	return s.port
}
