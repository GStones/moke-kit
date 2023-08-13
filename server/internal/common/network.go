package common

import "strings"

type ServerMod string

const (
	GrpcServerMod ServerMod = "grpc"
	TcpServerMod  ServerMod = "tcp"
	WsServerMod   ServerMod = "websocket"
	HttpServerMod ServerMod = "http"
	AllServerMod  ServerMod = "all"
)

func (sm ServerMod) HasGrpc() bool {
	return strings.Contains(string(sm), string(GrpcServerMod)) || strings.Contains(string(sm), string(AllServerMod))
}

func (sm ServerMod) HasTcp() bool {
	return strings.Contains(string(sm), string(TcpServerMod)) || strings.Contains(string(sm), string(AllServerMod))
}

func (sm ServerMod) HasWebsocket() bool {
	return strings.Contains(string(sm), string(WsServerMod)) || strings.Contains(string(sm), string(AllServerMod))
}

func (sm ServerMod) HasHttp() bool {
	return strings.Contains(string(sm), string(HttpServerMod)) || strings.Contains(string(sm), string(AllServerMod))
}

func (sm ServerMod) IsAll() bool {
	return strings.Contains(string(sm), string(AllServerMod))
}
