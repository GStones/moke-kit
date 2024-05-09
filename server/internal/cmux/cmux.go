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

	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/tools"
)

var (
	httpMatcher     cmux.Matcher
	grpcMatchWriter cmux.MatchWriter
	wsl             cmux.Matcher
	tcp             cmux.Matcher
)

func init() {
	grpcMatchWriter = cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc")
	wsl = cmux.HTTP1HeaderField("Upgrade", "websocket")
	httpMatcher = cmux.HTTP1Fast()
	tcp = cmux.Any()
}

type ConnectionMux struct {
	logger    *zap.Logger
	listener  net.Listener
	mux       cmux.CMux
	port      int32
	tlsConfig *tls.Config
}

func (cm *ConnectionMux) GrpcListener() (listener net.Listener, err error) {
	if cm.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = cm.mux.MatchWithWriters(grpcMatchWriter)
	}
	return
}

func (cm *ConnectionMux) WSListener() (listener net.Listener, err error) {
	if cm.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = cm.mux.Match(wsl)
	}
	return
}

func (cm *ConnectionMux) HTTPListener() (listener net.Listener, err error) {
	if cm.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = cm.mux.Match(httpMatcher)
	}

	return
}

func (cm *ConnectionMux) TCPListener() (listener net.Listener, err error) {
	if cm.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = cm.mux.Match(tcp)
	}
	return
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

func (cm *ConnectionMux) StopServing(_ context.Context) error {
	cm.mux.Close()
	return nil
}

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
