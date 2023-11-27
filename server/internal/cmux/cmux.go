package cmux

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
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
				zap.String("network", cm.listener.Addr().Network()),
				zap.String("address", cm.listener.Addr().String()),
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

func makeTlsConfig(tlsCert, tlsKey string) (config *tls.Config, err error) {
	sCert, _ := base64.StdEncoding.DecodeString(tlsCert)
	sKey, _ := base64.StdEncoding.DecodeString(tlsKey)

	if certificate, e := tls.X509KeyPair(sCert, sKey); e != nil {
		err = e
	} else {
		config = &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
	}

	return
}
