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

func (s *ZinxServer) ZinxServer() ziface.IServer {
	return s.server
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
