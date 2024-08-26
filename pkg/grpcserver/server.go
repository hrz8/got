package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Server   *grpc.Server
	addr     string
	listener net.Listener
	notify   chan error
}

func New(opts ...Option) *Server {
	grpcServer := grpc.NewServer()

	s := &Server{
		Server: grpcServer,
		addr:   net.JoinHostPort("", fmt.Sprintf("%d", defaultPort)),
		notify: make(chan error, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() {
	go func() {
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			s.notify <- err
			close(s.notify)
			return
		}

		s.listener = lis
		s.notify <- s.Server.Serve(s.listener)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	s.Server.GracefulStop()
	return nil
}
