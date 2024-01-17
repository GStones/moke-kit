package allocate

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	allocation "agones.dev/agones/pkg/allocation/go"
	"google.golang.org/grpc"
)

type MockAllocationServiceClient struct {
	URL string
}

// Allocate is a mock implementation of Allocate
func (m *MockAllocationServiceClient) Allocate(
	_ context.Context,
	_ *allocation.AllocationRequest,
	_ ...grpc.CallOption,
) (*allocation.AllocationResponse, error) {
	if m.URL == "" {
		return nil, fmt.Errorf("mock allocation service client url is empty")
	}
	hosts := strings.Split(m.URL, ":")
	port := int32(443)
	if len(hosts) == 2 {
		p, err := strconv.ParseInt(hosts[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("mock allocation service client url:%v is invalid", m.URL)
		}
		port = int32(p)
	}

	res := &allocation.AllocationResponse{
		Address: hosts[0],
		Ports: []*allocation.AllocationResponse_GameServerStatusPort{
			{
				Port: port,
			},
		},
	}
	return res, nil
}
