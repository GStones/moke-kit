package srpc

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gstones/moke-kit/utility"
)

type GatewayServer struct {
	logger   *zap.Logger
	server   *http.Server
	mux      *runtime.ServeMux
	listener net.Listener
	opts     []grpc.DialOption
}

func (gs *GatewayServer) StartServing(_ context.Context) error {
	go func() {
		if err := gs.server.Serve(gs.listener); err != nil {
			if !strings.Contains(err.Error(), "Server closed") {
				gs.logger.Error(
					"failed to serve grpc gateway",
					zap.String("network", gs.listener.Addr().Network()),
					zap.String("address", gs.listener.Addr().String()),
					zap.Error(err),
				)
			}
		}
	}()
	return nil
}

func (gs *GatewayServer) StopServing(ctx context.Context) error {
	if err := gs.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (gs *GatewayServer) GatewayServer() *http.Server {
	return gs.server
}

func (gs *GatewayServer) GatewayRuntimeMux() *runtime.ServeMux {
	return gs.mux
}

func (gs *GatewayServer) GatewayOption() []grpc.DialOption {
	return gs.opts
}

func (gs *GatewayServer) Endpoint() string {
	return gs.server.Addr
}

func NewGatewayServer(
	logger *zap.Logger,
	listener net.Listener,
) (result *GatewayServer, err error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(Matcher),
		runtime.WithOutgoingHeaderMatcher(Matcher),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: allowCORS(withLogger(mux)),
	}
	result = &GatewayServer{
		logger:   logger,
		server:   server,
		mux:      mux,
		opts:     opts,
		listener: listener,
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

func preflightHandler(w http.ResponseWriter, _ *http.Request) {
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
	case utility.TokenContextKey:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
