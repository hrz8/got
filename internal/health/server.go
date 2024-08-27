package health

import (
	"context"

	servicev1 "github.com/hrz8/got/pkg/pb/service/v1"
	"google.golang.org/grpc/codes"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type Server struct {
	servicev1.UnimplementedHealthServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Check(context.Context, *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	return &grpchealth.HealthCheckResponse{Status: grpchealth.HealthCheckResponse_SERVING}, nil
}

func (s *Server) Watch(_ *grpchealth.HealthCheckRequest, _ grpchealth.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}
