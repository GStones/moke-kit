package cmux

import (
	"context"
	"log"
	"moke-kit/server/network"
	"testing"

	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"
)

func newListener() net.Listener {
	return bufconn.Listen(256 * 1024)
}

func testConnectionMuxFactory(t *testing.T) (mux *ConnectionMux, err error) {
	type testCase struct {
		portNumber int
	}

	cases := []testCase{
		{
			portNumber: 0, // @TODO: PBAUMAN: This is passed in to avoid port conflicts in test. Looking for a better solution
		},
	}

	for _, tc := range cases {
		logger := zap.NewNop()
		port := network.Port(tc.portNumber)

		if tcpConnectionMux, err := NewConnectionMux(logger, port); err != nil {
			t.Error(err)
		} else {
			return tcpConnectionMux, err
		}
	}
	return
}

func TestNewTcpConnectionMux(t *testing.T) {
	logger := zap.NewNop()
	port := network.Port(9083)

	if _, err := NewConnectionMux(logger, port); err != nil {
		t.Fatal("Could not create a new TCP connection multiplexer:", err)
	}
}

func TestTcpConnectionMux_Port(t *testing.T) {
	if mux, err := testConnectionMuxFactory(t); err != nil {
		t.Fatal("Could not create a new TCP connection multiplexer:", err)
	} else {
		if port := mux.Port(); err != nil {
			t.Error("Couldn't return the TCP connection multiplexer port:", err)
		} else if testing.Verbose() {
			log.Println(port)
		}
	}
}

func TestTcpConnectionMux_GrpcConnectionMux(t *testing.T) {
	if mux, err := testConnectionMuxFactory(t); err != nil {
		t.Fatal("Could not create a new TCP connection multiplexer:", err)
	} else {
		mux.GrpcListener()
	}
}

func TestTcpConnectionMux_HttpConnectionMux(t *testing.T) {
	mux, _ := testConnectionMuxFactory(t)
	mux.HttpListener()
}

func TestTcpConnectionMux_StartAndStopListening(t *testing.T) {
	ctx := context.Background()
	mux, _ := testConnectionMuxFactory(t)

	if err := mux.StartServing(ctx); err != nil {
		t.Fatal("Couldn't begin serving with the TCP connection multiplexer:", err)

	}

	if err := mux.StopServing(ctx); err != nil {
		t.Fatal("Couldn't stop serving with the TCP connection multiplexer:", err)
	}
}
