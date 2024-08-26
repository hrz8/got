package grpcserver

import (
	"fmt"
	"net"
)

type Option func(*Server)

func Port(port uint16) Option {
	return func(s *Server) {
		s.addr = net.JoinHostPort("", fmt.Sprintf("%d", port))
	}
}
