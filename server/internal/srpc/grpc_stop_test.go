package srpc

import (
	"context"
	"net"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

func TestGrpcServerStopServingHonorsContext(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	gs := &GrpcServer{
		logger:   zaptest.NewLogger(t),
		server:   grpc.NewServer(),
		listener: ln,
	}
	if err := gs.StartServing(context.Background()); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := gs.StopServing(ctx); err != nil {
		t.Fatalf("StopServing: %v", err)
	}
}

func TestGrpcServerStopServingForceOnCancel(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	gs := &GrpcServer{
		logger:   zaptest.NewLogger(t),
		server:   grpc.NewServer(),
		listener: ln,
	}
	if err := gs.StartServing(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Already-cancelled ctx should force Stop path (may still return ctx.Err()).
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = gs.StopServing(ctx)
}
