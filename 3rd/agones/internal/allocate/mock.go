package allocate

import (
	"context"

	allocation "agones.dev/agones/pkg/allocation/go"
	"google.golang.org/grpc"
)

type MockAllocationServiceClient struct {
}

// Allocate is a mock implementation of Allocate
func (m *MockAllocationServiceClient) Allocate(
	_ context.Context,
	_ *allocation.AllocationRequest,
	_ ...grpc.CallOption,
) (*allocation.AllocationResponse, error) {
	return nil, nil
}
