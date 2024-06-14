package zinx

import (
	"context"

	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"
)

// ZinxServer is the struct for the zinx server.
// https://github.com/aceld/zinx
type ZinxServer struct {
	logger *zap.Logger
	server ziface.IServer
}

// ZinxServer returns the zinx server.
func (zs *ZinxServer) ZinxServer() ziface.IServer {
	return zs.server
}

// StartServing starts the zinx server.
func (zs *ZinxServer) StartServing(_ context.Context) error {
	go func() {
		zs.server.Serve()
	}()
	return nil
}

// StopServing stops the zinx server.
func (zs *ZinxServer) StopServing(_ context.Context) error {
	zs.server.Stop()
	return nil
}
