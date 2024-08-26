package greeter

import (
	"context"

	servicev1 "github.com/hrz8/got/pkg/pb/v1"
)

type Server struct {
	servicev1.UnimplementedGreeterServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SayHello(ctx context.Context, in *servicev1.HelloRequest) (*servicev1.HelloReply, error) {
	return &servicev1.HelloReply{Message: "Hello " + in.Name}, nil
}
