package srpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gstones/moke-kit/server/middlewares"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

// GatewaySecurity holds TLS settings used when the gateway dials gRPC.
type GatewaySecurity struct {
	TLSEnable    bool
	MTLSEnable   bool
	ClientCert   string
	ClientKey    string
	ServerCaCert string
	ServerName   string
}

// NewGrpcServer creates a new grpc server.
func NewGrpcServer(
	logger *zap.Logger,
	listener net.Listener,
	auth siface.IAuthMiddleware,
	deployment string,
	rateLimit int32,
	opts ...grpc.ServerOption,
) (result siface.IGrpcServer, err error) {
	deploy := utility.ParseDeployments(deployment)
	opts = middlewares.MakeServerOptions(logger, auth, deploy, rateLimit, opts...)
	result = &GrpcServer{
		logger:   logger,
		listener: listener,
		server:   grpc.NewServer(opts...),
	}
	return result, nil
}

// NewGatewayServer creates a new gateway server.
func NewGatewayServer(
	logger *zap.Logger,
	listener net.Listener,
	sec GatewaySecurity,
	corsAllowOrigins string,
) (result *GatewayServer, err error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(matcher),
		runtime.WithOutgoingHeaderMatcher(matcher),
	)
	dialOpts, err := gatewayDialOptions(sec)
	if err != nil {
		return nil, err
	}
	origins := parseCORSAllowOrigins(corsAllowOrigins)
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: allowCORS(withLogger(mux), origins),
	}
	result = &GatewayServer{
		logger:           logger,
		server:           server,
		mux:              mux,
		opts:             dialOpts,
		listener:         listener,
		corsAllowOrigins: origins,
	}
	return
}

func parseCORSAllowOrigins(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func gatewayDialOptions(sec GatewaySecurity) ([]grpc.DialOption, error) {
	if !sec.TLSEnable && !sec.MTLSEnable {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if sec.ServerName != "" {
		tlsConfig.ServerName = sec.ServerName
	}
	if sec.ServerCaCert != "" {
		caBytes, err := os.ReadFile(sec.ServerCaCert)
		if err != nil {
			return nil, err
		}
		ca := x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM(caBytes); !ok {
			return nil, fmt.Errorf("failed to parse server CA %q", sec.ServerCaCert)
		}
		tlsConfig.RootCAs = ca
	}
	if sec.MTLSEnable {
		if sec.ClientCert == "" || sec.ClientKey == "" {
			return nil, fmt.Errorf("client certificate and key are required for mTLS gateway dial")
		}
		cert, err := tls.LoadX509KeyPair(sec.ClientCert, sec.ClientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))}, nil
}
