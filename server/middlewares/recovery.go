package middlewares

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpcPanicRecoveryHandler is a grpc recovery handler that returns an Internal error
func grpcPanicRecoveryHandler(p any) (err error) {
	return status.Errorf(codes.Internal, "recovered from panic: %v", p)
}
