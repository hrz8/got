package httpserver

import (
	"fmt"
	"net"
	"time"
)

type Option func(*Server)

func Port(port uint16) Option {
	return func(s *Server) {
		s.Server.Addr = net.JoinHostPort("", fmt.Sprintf("%d", port))
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func ReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.ReadHeaderTimeout = timeout
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.WriteTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.IdleTimeout = timeout
	}
}

func AllowedOrigins(origins []string) Option {
	return func(s *Server) {
		s.allowedOrigins = origins
	}
}

func AllowedHeaders(headers []string) Option {
	return func(s *Server) {
		s.allowedHeaders = headers
	}
}
