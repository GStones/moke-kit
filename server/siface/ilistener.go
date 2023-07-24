package siface

import "net"

type HasGrpcListener interface {
	GrpcListener() (net.Listener, error)
}

type HasHttpListener interface {
	HttpListener() (net.Listener, error)
}
