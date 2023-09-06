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
	"google.golang.org/grpc/test/bufconn"
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

func (m *ConnectionMux) GrpcListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.MatchWithWriters(grpcMatchWriter)
	}
	return
}

func (m *ConnectionMux) WSListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.Match(wsl)
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

func (m *ConnectionMux) TcpListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.Match(tcp)
	}
	return
}

func (m *ConnectionMux) run() error {
	if listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", m.port)); err != nil {
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
		zap.Int32("port", m.port),
		zap.Bool("tls", m.tlsConfig != nil),
	)

	m.mux = cmux.New(m.listener)

	go func() {
		m.mux.Serve()
	}()

	return nil
}

func (m *ConnectionMux) StartServing(_ context.Context) error {
	return nil
}

func (m *ConnectionMux) StopServing(_ context.Context) error {
	return m.listener.Close()
}

type TestConnectionMux struct {
	listener *bufconn.Listener
}

func (m *TestConnectionMux) GrpcListener() (net.Listener, error) {
	return m.listener, nil
}

func (m *TestConnectionMux) HttpListener() (net.Listener, error) {
	return m.listener, nil
}

func (m *TestConnectionMux) StartServing(ctx context.Context) error {
	return nil
}

func (m *TestConnectionMux) StopServing(ctx context.Context) error {
	return nil
}

func (m *TestConnectionMux) Dial() (net.Conn, error) {
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
