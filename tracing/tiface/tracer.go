package tiface

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type ITracer interface {
	Runner

	StreamClientInterceptor() grpc.StreamClientInterceptor
	StreamServerInterceptor() grpc.StreamServerInterceptor
	UnaryClientInterceptor() grpc.UnaryClientInterceptor
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	ServeMux() HttpServeMux

	SpanContext(context.Context) (SpanContext, bool)
}

type HttpServeMux interface {
	http.Handler

	Handle(string, http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type SpanContext interface {
	TraceID() uint64
	SpanID() uint64
}
