package srpc

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type hangService interface {
	isHang()
}

type hangServer struct {
	started chan struct{}
	release chan struct{}
}

func (*hangServer) isHang() {}

func TestGrpcServerStopServingIdle(t *testing.T) {
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

func TestGrpcServerStopServingForceWhenRPCBlocks(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	hang := &hangServer{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	s := grpc.NewServer()
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "test.Hang",
		HandlerType: (*hangService)(nil),
		Methods: []grpc.MethodDesc{{
			MethodName: "Wait",
			Handler: func(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
				if err := dec(new(emptypb.Empty)); err != nil {
					return nil, err
				}
				h := srv.(*hangServer)
				close(h.started)
				select {
				case <-h.release:
					return &emptypb.Empty{}, nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			},
		}},
	}, hang)

	gs := &GrpcServer{
		logger:   zaptest.NewLogger(t),
		server:   s,
		listener: ln,
	}
	if err := gs.StartServing(context.Background()); err != nil {
		t.Fatal(err)
	}

	conn, err := grpc.NewClient(
		ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	invokeDone := make(chan error, 1)
	go func() {
		invokeDone <- conn.Invoke(context.Background(), "/test.Hang/Wait", &emptypb.Empty{}, &emptypb.Empty{})
	}()

	select {
	case <-hang.started:
	case <-time.After(2 * time.Second):
		t.Fatal("RPC did not start")
	}

	// GracefulStop cannot finish while Hang/Wait is in flight; force Stop via short deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	err = gs.StopServing(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("StopServing err=%v want deadline exceeded", err)
	}

	close(hang.release)
	select {
	case <-invokeDone:
	case <-time.After(2 * time.Second):
		t.Fatal("client Invoke did not return after force stop")
	}
}
