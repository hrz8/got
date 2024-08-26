package health

import (
	"context"

	servicev1 "github.com/hrz8/got/pkg/pb/v1"
)

type Server struct {
	servicev1.UnimplementedHealthServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Check(ctx context.Context, in *servicev1.HealthRequest) (*servicev1.HealthResponse, error) {
	return &servicev1.HealthResponse{Status: "OK"}, nil
}
