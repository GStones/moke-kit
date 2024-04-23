package allocate

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	allocation "agones.dev/agones/pkg/allocation/go"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
)

type MockAllocationServiceClient struct {
	hosts []string
}

func CreateMockAllocationServiceClient(hosts []string) *MockAllocationServiceClient {
	return &MockAllocationServiceClient{hosts: hosts}
}

// Allocate is a mock implementation of Allocate
func (m *MockAllocationServiceClient) Allocate(
	_ context.Context,
	_ *allocation.AllocationRequest,
	_ ...grpc.CallOption,
) (*allocation.AllocationResponse, error) {
	if len(m.hosts) <= 0 {
		return nil, fmt.Errorf("mock allocation service client url is empty")
	}
	index := rand.Intn(len(m.hosts))
	url := m.hosts[index]
	hosts := strings.Split(url, ":")
	port := int32(443)
	if len(hosts) == 2 {
		p, err := strconv.ParseInt(hosts[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("mock allocation service client url:%v is invalid", url)
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
