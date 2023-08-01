package srpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"moke-kit/server/siface"
)

type GatewayServer struct {
	logger   *zap.Logger
	server   *http.Server
	mux      *runtime.ServeMux
	listener siface.IHttpListener
	opts     []grpc.DialOption
	endpoint string
}

func (s *GatewayServer) StartServing(_ context.Context) error {
	if listener, err := s.listener.HttpListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving srpc gateway",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
		)
		go func() {
			if err := s.server.Serve(listener); err != nil {
				if !strings.Contains(err.Error(), "Server closed") {
					s.logger.Error(
						"failed to serve srpc gateway",
						zap.String("network", listener.Addr().Network()),
						zap.String("address", listener.Addr().String()),
						zap.Error(err),
					)
				}
			}
		}()
	}
	return nil
}

func (s *GatewayServer) StopServing(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *GatewayServer) GatewayServer() *http.Server {
	return s.server
}

func (s *GatewayServer) GatewayRuntimeMux() *runtime.ServeMux {
	return s.mux
}

func (s *GatewayServer) GatewayOption() []grpc.DialOption {
	return s.opts
}

func (s *GatewayServer) Endpoint() string {
	return s.endpoint
}

func NewGatewayServer(
	logger *zap.Logger,
	listener siface.IHttpListener,
	port int32,
	endpoint string,
) (result *GatewayServer, err error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(Matcher),
		runtime.WithOutgoingHeaderMatcher(Matcher),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: allowCORS(withLogger(mux)),
	}
	result = &GatewayServer{
		logger:   logger,
		server:   server,
		mux:      mux,
		opts:     opts,
		listener: listener,
		endpoint: endpoint,
	}
	return
}
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	headers := []string{"*"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func Matcher(key string) (string, bool) {
	switch key {
	case TokenContextKey:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
