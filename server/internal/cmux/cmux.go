package cmux

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/tools"
)

var (
	httpMatcher     cmux.Matcher
	grpcMatchWriter cmux.MatchWriter
	wsl             cmux.Matcher
)

func init() {
	grpcMatchWriter = cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc")
	wsl = cmux.HTTP1HeaderField("Upgrade", "websocket")
	httpMatcher = cmux.HTTP1Fast()
}

// ConnectionMux is the struct for the connection mux.
// https://github.com/soheilhy/cmux
type ConnectionMux struct {
	logger    *zap.Logger
	listener  net.Listener
	mux       cmux.CMux
	port      int32
	tlsConfig *tls.Config
}

// GrpcListener returns the grpc listener from the connection mux.
func (cm *ConnectionMux) GrpcListener() (net.Listener, error) {
	if cm.mux == nil {
		if err := cm.init(); err != nil {
			return nil, err
		}
	}
	return cm.mux.MatchWithWriters(grpcMatchWriter), nil
}

// WSListener returns the websocket listener from the connection mux.
func (cm *ConnectionMux) WSListener() (net.Listener, error) {
	if cm.mux == nil {
		if err := cm.init(); err != nil {
			return nil, err
		}
	}
	return cm.mux.Match(wsl), nil
}

// HTTPListener returns the http listener from the connection mux.
func (cm *ConnectionMux) HTTPListener() (listener net.Listener, err error) {
	if cm.mux == nil {
		if err := cm.init(); err != nil {
			return nil, err
		}
	}
	return cm.mux.Match(httpMatcher), nil
}

func (cm *ConnectionMux) init() error {
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cm.port)); err != nil {
		return err
	} else {
		if cm.tlsConfig != nil {
			listener = tls.NewListener(listener, cm.tlsConfig)
		}
		cm.mux = cmux.New(listener)
	}
	return nil
}

// StartServing starts the connection mux.
func (cm *ConnectionMux) StartServing(_ context.Context) error {
	go func() {
		if err := cm.mux.Serve(); err != nil {
			cm.logger.Error(
				"failed to serve",
				zap.Any("listener", cm.listener),
				zap.Error(err),
			)
		} else {
			cm.logger.Info(
				"multiplexing traffic",
				zap.String("network", cm.listener.Addr().Network()),
				zap.String("address", cm.listener.Addr().String()),
				zap.Int32("port", cm.port),
				zap.Bool("tls", cm.tlsConfig != nil),
			)
		}
	}()
	return nil
}

// StopServing stops the connection mux.
func (cm *ConnectionMux) StopServing(_ context.Context) error {
	cm.mux.Close()
	return nil
}

// NewConnectionMux creates a new connection mux.
// Watch the tls certificate and reload it when it changes.
func makeTLSConfig(logger *zap.Logger, tlsCert, tlsKey string, clientCa string) (*tls.Config, error) {
	if cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey); err != nil {
		return nil, err
	} else if caBytes, err := os.ReadFile(clientCa); err != nil {
		return nil, err
	} else {
		ca := x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM(caBytes); !ok {
			return nil, fmt.Errorf("failed to parse %v ", clientCa)
		}

		tlsCertValue := atomic.Value{}
		tlsCertValue.Store(cert)
		p, _ := path.Split(tlsCert)
		if _, err := tools.Watch(logger, p, time.Second*10, func() {
			logger.Info("service reloading x509 key pair")
			c, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
			if err != nil {
				logger.Error("service failed to load x509 key pair", zap.Error(err))
				return
			}
			tlsCertValue.Store(c)
		}); err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				c := tlsCertValue.Load()
				if c == nil {
					return nil, fmt.Errorf("certificate not loaded")
				}
				res := c.(tls.Certificate)
				return &res, nil
			},
			ClientCAs: ca,
		}
		return tlsConfig, nil
	}

}
