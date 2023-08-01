package dfx

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "moke-kit/demo/api/gen/demo/api"
)

type DemoClient struct {
	client pb.HelloClient
}

func (dc *DemoClient) Hi(ctx context.Context, message string) (string, error) {
	if resp, err := dc.client.Hi(ctx, &pb.HiRequest{
		Message: message,
	}); err != nil {
		return "", err
	} else {
		return resp.Message, nil
	}
}

func newClient(target string, secure bool) (cConn *grpc.ClientConn, err error) {
	var opts []grpc.DialOption
	if secure {
		//TODO add secure
		//opts = append(opts, grpc.WithTransportCredentials(utils.GetClientCreds()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewDemoClient(
	target string,
	secure bool,
) (result *DemoClient, err error) {
	cConn, err := newClient(target, secure)
	if err != nil {
		return nil, err
	}
	result = &DemoClient{
		client: pb.NewHelloClient(cConn),
	}
	return
}
