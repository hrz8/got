package health

import (
	"context"

	servicev1 "github.com/hrz8/got/pkg/pb/service/v1"
)

type Server struct {
	servicev1.UnimplementedHealthServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Check(ctx context.Context, in *servicev1.CheckRequest) (*servicev1.CheckResponse, error) {
	return &servicev1.CheckResponse{Status: "OK"}, nil
}
