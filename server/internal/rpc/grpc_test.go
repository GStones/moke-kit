package rpc

//
//import (
//	"moke-kit/server/internal"
//	"moke-kit/server/network"
//	"moke-kit/server/siface"
//	"testing"
//
//	"context"
//
//	"net"
//
//	"go.uber.org/zap"
//	"google.golang.org/grpc"
//)
//
//type GrpcServerFactory func(logger *zap.Logger, listener siface.HasGrpcListener, port Port, opts ...grpc.ServerOption) (GrpcServer, error)
//
//type GrpcListener struct{}
//
//func (l GrpcListener) GrpcListener() (net.Listener, error) {
//	return internal.newListener(), nil
//}
//
//func testNewGrpcServer(t *testing.T, f GrpcServerFactory) {
//	listener := GrpcListener{}
//	logger := zap.NewNop()
//
//	type TestCase struct {
//		listener siface.HasGrpcListener
//		port     network.Port
//		logger   *zap.Logger
//		validate func(i int, server GrpcServer, err error)
//	}
//
//	validate := func(i int, server GrpcServer, err error) {
//		if server == nil || err != nil {
//			t.Fatal(i, err, server)
//		}
//	}
//
//	cases := []TestCase{
//		{
//			listener: listener,
//			port:     newPort(),
//			logger:   logger,
//			validate: validate,
//		},
//	}
//
//	for i, tc := range cases {
//		server, err := f(tc.logger, tc.listener, tc.port)
//		tc.validate(i, server, err)
//	}
//}
//
//func createNewTestTcpServer(t *testing.T) *GrpcServer {
//	listener := GrpcListener{}
//	port := newPort()
//	logger := zap.NewNop()
//
//	testServer, err := NewTcpGrpcServer(logger, listener, port)
//	if err != nil {
//		t.Fatal("Couldn't create the rpc server:", err)
//	}
//	return testServer
//}
//
//func createNewTestServer(t *testing.T) *TestGrpcServer {
//	port := newPort()
//	logger := zap.NewNop()
//
//	testServer := NewTestGrpcServer(logger, port)
//
//	return testServer
//}
//
//func TestNewTcpGrpcServer(t *testing.T) {
//	testNewGrpcServer(t, func(logger *zap.Logger, listener HasGrpcListener, port Port, opts ...grpc.ServerOption) (GrpcServer, error) {
//		return NewTcpGrpcServer(logger, listener, port, opts...)
//	})
//}
//
//func TestNewTestGrpcServer(t *testing.T) {
//	testNewGrpcServer(t, func(logger *zap.Logger, listener HasGrpcListener, port Port, opts ...grpc.ServerOption) (GrpcServer, error) {
//		return NewTestGrpcServer(logger, port, opts...), nil
//	})
//}
//
//func TestDialTCPServer(t *testing.T) {
//	testServer := createNewTestTcpServer(t)
//	testServer.Dial("")
//}
//
//func TestDialTestServer(t *testing.T) {
//	testServer := createNewTestServer(t)
//	testServer.Dial("")
//}
//
//func TestTcpGrpcServerServing(t *testing.T) {
//	ctx := context.Background()
//	testServer := createNewTestTcpServer(t)
//	testServer.StartServing(ctx)
//	testServer.StopServing(ctx)
//}
//
//func TestTestGrpcServerServing(t *testing.T) {
//	ctx := context.Background()
//	testServer := createNewTestServer(t)
//	testServer.StartServing(ctx)
//	testServer.StopServing(ctx)
//}
