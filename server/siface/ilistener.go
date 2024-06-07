package siface

import "net"

type IGrpcListener interface {
	GrpcListener() (net.Listener, error)
}

type IHttpListener interface {
	HTTPListener() (net.Listener, error)
}

type IWebSocketListener interface {
	WSListener() (net.Listener, error)
}
