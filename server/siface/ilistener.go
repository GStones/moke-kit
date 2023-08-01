package siface

import "net"

type IGrpcListener interface {
	GrpcListener() (net.Listener, error)
}

type IHttpListener interface {
	HttpListener() (net.Listener, error)
}
