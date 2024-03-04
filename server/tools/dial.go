package tools

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gstones/moke-kit/server/middlewares"
)

// Timeout grpc dial timeout
const Timeout = 2 * time.Second

// DialInsecure dial insecure grpc
func DialInsecure(target string) (cConn *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	logger, _ := zap.NewDevelopment()
	opts := middlewares.MakeClientOptions(logger)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialWithSecurity dial grpc with security
func DialWithSecurity(
	target string,
	clientCert, clientKey, serverName, serverCa string,
) (cConn *grpc.ClientConn, err error) {
	logger, _ := zap.NewDevelopment()
	opts := middlewares.MakeClientOptions(logger)
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	if serverName != "" {
		tlsConfig.ServerName = serverName
	}
	if serverCa != "" {
		ca := x509.NewCertPool()
		caBytes, err := os.ReadFile(serverCa)
		if err != nil {
			return nil, err
		}
		if ok := ca.AppendCertsFromPEM(caBytes); !ok {
			return nil, fmt.Errorf("failed to parse %q", serverCa)
		}
		tlsConfig.RootCAs = ca
	}
	opts = append(
		opts,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
