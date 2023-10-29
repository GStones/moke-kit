package utility

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const Timeout = 1 * time.Second

func DialWithOptions(target string, secure bool) (cConn *grpc.ClientConn, err error) {
	var opts []grpc.DialOption
	if secure {
		//TODO add secure
		//opts = append(opts, grpc.WithTransportCredentials(utils.GetClientCreds()))
	} else {
		grpc.WithResolvers()
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
