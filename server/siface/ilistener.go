package siface

import "net"

type IGrpcListener interface {
	GrpcListener() (net.Listener, error)
}

type IHttpListener interface {
	HttpListener() (net.Listener, error)
}

type IWebSocketListener interface {
	WSListener() (net.Listener, error)
}

type ITcpListener interface {
	TcpListener() (net.Listener, error)
}
