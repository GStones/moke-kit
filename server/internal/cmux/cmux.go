package cmux

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"
	"moke-kit/server/network"
	"net"
)

var (
	httpMatcher     cmux.Matcher
	grpcMatchWriter cmux.MatchWriter
)

func init() {
	grpcMatchWriter = cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/rpc")
	httpMatcher = cmux.HTTP1Fast()
}

type ConnectionMux struct {
	logger    *zap.Logger
	listener  net.Listener
	mux       cmux.CMux
	port      network.Port
	tlsConfig *tls.Config
}

func (m *ConnectionMux) GrpcListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.MatchWithWriters(grpcMatchWriter)
	}
	return
}

func (m *ConnectionMux) HttpListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.Match(httpMatcher)
	}

	return
}

func (m *ConnectionMux) StartServing(_ context.Context) error {
	if listener, err := network.NewTcpListener(m.port); err != nil {
		return err
	} else {
		if m.tlsConfig != nil {
			m.listener = tls.NewListener(listener, m.tlsConfig)
		} else {
			m.listener = listener
		}
	}

	m.logger.Info(
		"multiplexing traffic",
		zap.String("network", m.listener.Addr().Network()),
		zap.String("address", m.listener.Addr().String()),
		zap.Int("port", m.port.Value()),
		zap.Bool("tls", m.tlsConfig != nil),
	)

	m.mux = cmux.New(m.listener)

	go func() {
		m.mux.Serve()
	}()

	return nil
}

func (m *ConnectionMux) StopServing(_ context.Context) error {
	return m.listener.Close()
}

func (m *ConnectionMux) Port() network.Port {
	return m.port
}

type TestTcpConnectionMux struct {
	listener *bufconn.Listener
}

func (m *TestTcpConnectionMux) Port() network.Port {
	return network.Port(0)
}

func (m *TestTcpConnectionMux) GrpcListener() (net.Listener, error) {
	return m.listener, nil
}

func (m *TestTcpConnectionMux) HttpListener() (net.Listener, error) {
	return m.listener, nil
}

func (m *TestTcpConnectionMux) StartServing(ctx context.Context) error {
	return nil
}

func (m *TestTcpConnectionMux) StopServing(ctx context.Context) error {
	return nil
}

func (m *TestTcpConnectionMux) Dial() (net.Conn, error) {
	return m.listener.Dial()
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
