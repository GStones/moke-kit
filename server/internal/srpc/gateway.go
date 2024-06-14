package srpc

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/utility"
)

// GatewayServer is the struct for the gateway server.
type GatewayServer struct {
	logger   *zap.Logger
	server   *http.Server
	mux      *runtime.ServeMux
	listener net.Listener
	opts     []grpc.DialOption
}

// StartServing starts the gateway server.
func (gs *GatewayServer) StartServing(_ context.Context) error {
	gs.logger.Info(
		"grpc gateway start serving",
		zap.String("network", gs.listener.Addr().Network()),
		zap.String("address", gs.listener.Addr().String()),
	)
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

// StopServing stops the gateway server.
func (gs *GatewayServer) StopServing(ctx context.Context) error {
	if err := gs.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// GatewayServer returns the gateway server.
func (gs *GatewayServer) GatewayServer() *http.Server {
	return gs.server
}

// GatewayRuntimeMux returns the gateway runtime mux.
func (gs *GatewayServer) GatewayRuntimeMux() *runtime.ServeMux {
	return gs.mux
}

// GatewayOption returns the gateway option.
func (gs *GatewayServer) GatewayOption() []grpc.DialOption {
	return gs.opts
}

// Endpoint returns the endpoint.
func (gs *GatewayServer) Endpoint() string {
	return gs.server.Addr
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

func matcher(key string) (string, bool) {
	switch key {
	case string(utility.TokenContextKey):
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
