package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/test/utils"
	"github.com/gstones/zinx/ziface"
)

// MockServer is a mock implementation of IServer
type MockServer struct {
	mock.Mock
}

func (m *MockServer) StartServing(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockServer) StopServing(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockConnectionMux is a mock implementation of IConnectionMux
type MockConnectionMux struct {
	MockServer
}

func (m *MockConnectionMux) GrpcListener() (net.Listener, error) {
	args := m.Called()
	return args.Get(0).(net.Listener), args.Error(1)
}

func (m *MockConnectionMux) HTTPListener() (net.Listener, error) {
	args := m.Called()
	return args.Get(0).(net.Listener), args.Error(1)
}

// MockGrpcServer is a mock implementation of IGrpcServer
type MockGrpcServer struct {
	MockServer
	server *grpc.Server
}

func (m *MockGrpcServer) GrpcServer() *grpc.Server {
	args := m.Called()
	return args.Get(0).(*grpc.Server)
}

// MockGatewayServer is a mock implementation of IGatewayServer
type MockGatewayServer struct {
	MockServer
}

func (m *MockGatewayServer) GatewayRuntimeMux() *runtime.ServeMux {
	args := m.Called()
	return args.Get(0).(*runtime.ServeMux)
}

func (m *MockGatewayServer) GatewayOption() []grpc.DialOption {
	args := m.Called()
	return args.Get(0).([]grpc.DialOption)
}

func (m *MockGatewayServer) GatewayServer() *http.Server {
	args := m.Called()
	return args.Get(0).(*http.Server)
}

func (m *MockGatewayServer) Endpoint() string {
	args := m.Called()
	return args.String(0)
}

// MockZinxServer is a mock implementation of IZinxServer
type MockZinxServer struct {
	MockServer
}

func (m *MockZinxServer) ZinxServer() ziface.IServer {
	args := m.Called()
	return args.Get(0).(ziface.IServer)
}

// MockListener is a mock implementation of net.Listener
type MockListener struct {
	mock.Mock
}

func (m *MockListener) Accept() (net.Conn, error) {
	args := m.Called()
	return args.Get(0).(net.Conn), args.Error(1)
}

func (m *MockListener) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockListener) Addr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

// TestServerInterfaces tests the server interfaces
func TestServerInterfaces(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("IServer", func(t *testing.T) {
		mockServer := &MockServer{}

		// Test StartServing
		mockServer.On("StartServing", mock.Anything).Return(nil)
		err := mockServer.StartServing(helper.Context())
		helper.AssertNoError(err)

		// Test StopServing
		mockServer.On("StopServing", mock.Anything).Return(nil)
		err = mockServer.StopServing(helper.Context())
		helper.AssertNoError(err)

		mockServer.AssertExpectations(t)
	})

	t.Run("IConnectionMux", func(t *testing.T) {
		mockMux := &MockConnectionMux{}
		mockListener := &MockListener{}

		// Test GrpcListener
		mockMux.On("GrpcListener").Return(mockListener, nil)
		listener, err := mockMux.GrpcListener()
		helper.RequireNoError(err)
		helper.AssertNotNil(listener)

		// Test HTTPListener
		mockMux.On("HTTPListener").Return(mockListener, nil)
		listener, err = mockMux.HTTPListener()
		helper.RequireNoError(err)
		helper.AssertNotNil(listener)

		mockMux.AssertExpectations(t)
	})

	t.Run("IGrpcServer", func(t *testing.T) {
		mockGrpcServer := &MockGrpcServer{}
		grpcServer := grpc.NewServer()

		// Test GrpcServer
		mockGrpcServer.On("GrpcServer").Return(grpcServer)
		server := mockGrpcServer.GrpcServer()
		helper.AssertNotNil(server)
		helper.AssertEqual(grpcServer, server)

		mockGrpcServer.AssertExpectations(t)
	})

	t.Run("IGatewayServer", func(t *testing.T) {
		mockGateway := &MockGatewayServer{}

		// Setup mock returns
		mux := runtime.NewServeMux()
		httpServer := &http.Server{}
		dialOpts := []grpc.DialOption{}
		endpoint := "localhost:8080"

		mockGateway.On("GatewayRuntimeMux").Return(mux)
		mockGateway.On("GatewayServer").Return(httpServer)
		mockGateway.On("GatewayOption").Return(dialOpts)
		mockGateway.On("Endpoint").Return(endpoint)

		// Test methods
		resultMux := mockGateway.GatewayRuntimeMux()
		helper.AssertEqual(mux, resultMux)

		resultServer := mockGateway.GatewayServer()
		helper.AssertEqual(httpServer, resultServer)

		resultOpts := mockGateway.GatewayOption()
		helper.AssertEqual(dialOpts, resultOpts)

		resultEndpoint := mockGateway.Endpoint()
		helper.AssertEqual(endpoint, resultEndpoint)

		mockGateway.AssertExpectations(t)
	})
}

// TestServerLifecycle tests server lifecycle management
func TestServerLifecycle(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ServerStartStop", func(t *testing.T) {
		mockServer := &MockServer{}

		// Test successful start
		mockServer.On("StartServing", mock.Anything).Return(nil)
		err := mockServer.StartServing(helper.Context())
		helper.AssertNoError(err)

		// Test successful stop
		mockServer.On("StopServing", mock.Anything).Return(nil)
		err = mockServer.StopServing(helper.Context())
		helper.AssertNoError(err)

		mockServer.AssertExpectations(t)
	})

	t.Run("ServerStartStopWithTimeout", func(t *testing.T) {
		mockServer := &MockServer{}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(helper.Context(), 5*time.Second)
		defer cancel()

		// Test start with timeout
		mockServer.On("StartServing", mock.Anything).Return(nil)
		err := mockServer.StartServing(ctx)
		helper.AssertNoError(err)

		// Test stop with timeout
		mockServer.On("StopServing", mock.Anything).Return(nil)
		err = mockServer.StopServing(ctx)
		helper.AssertNoError(err)

		mockServer.AssertExpectations(t)
	})

	t.Run("ConcurrentServers", func(t *testing.T) {
		numServers := 5
		servers := make([]*MockServer, numServers)
		done := make(chan bool, numServers)

		// Create and start multiple servers concurrently
		for i := 0; i < numServers; i++ {
			servers[i] = &MockServer{}
			servers[i].On("StartServing", mock.Anything).Return(nil)
			servers[i].On("StopServing", mock.Anything).Return(nil)

			go func(server *MockServer, id int) {
				defer func() { done <- true }()

				// Start server
				err := server.StartServing(helper.Context())
				require.NoError(t, err)

				// Simulate some work
				time.Sleep(100 * time.Millisecond)

				// Stop server
				err = server.StopServing(helper.Context())
				require.NoError(t, err)
			}(servers[i], i)
		}

		// Wait for all servers to complete
		timeout := time.After(10 * time.Second)
		for i := 0; i < numServers; i++ {
			select {
			case <-done:
				// Server completed
			case <-timeout:
				t.Fatal("Timeout waiting for servers to complete")
			}
		}

		// Verify expectations
		for _, server := range servers {
			server.AssertExpectations(t)
		}
	})
}

// TestServerTypes tests different server type implementations
func TestServerTypes(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("GrpcServerType", func(t *testing.T) {
		mockGrpcServer := &MockGrpcServer{}
		grpcServer := grpc.NewServer()

		// Setup expectations
		mockGrpcServer.On("GrpcServer").Return(grpcServer)
		mockGrpcServer.On("StartServing", mock.Anything).Return(nil)
		mockGrpcServer.On("StopServing", mock.Anything).Return(nil)

		// Test interface compliance
		var server siface.IGrpcServer = mockGrpcServer
		helper.AssertNotNil(server)

		// Test methods
		resultServer := server.GrpcServer()
		helper.AssertEqual(grpcServer, resultServer)

		err := server.StartServing(helper.Context())
		helper.AssertNoError(err)

		err = server.StopServing(helper.Context())
		helper.AssertNoError(err)

		mockGrpcServer.AssertExpectations(t)
	})

	t.Run("GatewayServerType", func(t *testing.T) {
		mockGateway := &MockGatewayServer{}

		// Setup mock returns
		mux := runtime.NewServeMux()
		httpServer := &http.Server{}
		dialOpts := []grpc.DialOption{}
		endpoint := "localhost:9090"

		// Setup expectations
		mockGateway.On("GatewayRuntimeMux").Return(mux)
		mockGateway.On("GatewayServer").Return(httpServer)
		mockGateway.On("GatewayOption").Return(dialOpts)
		mockGateway.On("Endpoint").Return(endpoint)
		mockGateway.On("StartServing", mock.Anything).Return(nil)
		mockGateway.On("StopServing", mock.Anything).Return(nil)

		// Test interface compliance
		var server siface.IGatewayServer = mockGateway
		helper.AssertNotNil(server)

		// Test methods
		resultMux := server.GatewayRuntimeMux()
		helper.AssertNotNil(resultMux)

		resultServer := server.GatewayServer()
		helper.AssertNotNil(resultServer)

		resultOpts := server.GatewayOption()
		helper.AssertNotNil(resultOpts)

		resultEndpoint := server.Endpoint()
		helper.AssertEqual(endpoint, resultEndpoint)

		err := server.StartServing(helper.Context())
		helper.AssertNoError(err)

		err = server.StopServing(helper.Context())
		helper.AssertNoError(err)

		mockGateway.AssertExpectations(t)
	})
}

// TestConnectionMuxFunctionality tests connection multiplexing functionality
func TestConnectionMuxFunctionality(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ListenerCreation", func(t *testing.T) {
		mockMux := &MockConnectionMux{}
		mockListener := &MockListener{}

		// Setup expectations
		mockMux.On("GrpcListener").Return(mockListener, nil)
		mockMux.On("HTTPListener").Return(mockListener, nil)
		mockMux.On("StartServing", mock.Anything).Return(nil)
		mockMux.On("StopServing", mock.Anything).Return(nil)

		// Test interface compliance
		var mux siface.IConnectionMux = mockMux
		helper.AssertNotNil(mux)

		// Test gRPC listener
		grpcListener, err := mux.GrpcListener()
		helper.RequireNoError(err)
		helper.AssertNotNil(grpcListener)

		// Test HTTP listener
		httpListener, err := mux.HTTPListener()
		helper.RequireNoError(err)
		helper.AssertNotNil(httpListener)

		// Test server lifecycle
		err = mux.StartServing(helper.Context())
		helper.AssertNoError(err)

		err = mux.StopServing(helper.Context())
		helper.AssertNoError(err)

		mockMux.AssertExpectations(t)
	})
}

// BenchmarkServerOperations benchmarks server operations
func BenchmarkServerOperations(b *testing.B) {
	mockServer := &MockServer{}
	mockServer.On("StartServing", mock.Anything).Return(nil)
	mockServer.On("StopServing", mock.Anything).Return(nil)

	ctx := context.Background()

	b.Run("StartServing", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := mockServer.StartServing(ctx)
				if err != nil {
					b.Error(err)
				}
			}
		})
	})

	b.Run("StopServing", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := mockServer.StopServing(ctx)
				if err != nil {
					b.Error(err)
				}
			}
		})
	})
}
