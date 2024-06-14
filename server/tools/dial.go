package tools

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"go.uber.org/atomic"
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
	tlsConfig, err := makeTls(logger, clientCert, clientKey, serverName, serverCa)
	if err != nil {
		return nil, err
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

// makeTls make tls config
// Watch the client certificate and reload it when it changes
func makeTls(logger *zap.Logger, clientCert, clientKey, serverName, serverCa string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	tlsCert := atomic.Value{}
	tlsCert.Store(cert)

	p, _ := path.Split(clientCert)
	if _, err := Watch(logger, p, time.Second*10, func() {
		logger.Info("client reload certificate")
		c, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return
		}
		tlsCert.Store(c)
	}); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			c := tlsCert.Load()
			if c == nil {
				return nil, fmt.Errorf("certificate not loaded")
			}
			cert := c.(tls.Certificate)
			return &cert, nil
		},
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
	return tlsConfig, nil
}
