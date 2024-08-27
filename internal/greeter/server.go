package greeter

import (
	"context"

	servicev1 "github.com/hrz8/got/pkg/pb/service/v1"
)

type Server struct {
	servicev1.UnimplementedGreeterServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SayHello(ctx context.Context, in *servicev1.SayHelloRequest) (*servicev1.SayHelloResponse, error) {
	return &servicev1.SayHelloResponse{Message: "Hello " + in.Name}, nil
}
