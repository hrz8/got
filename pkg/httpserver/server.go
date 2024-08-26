package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type Server struct {
	Server          *http.Server
	allowedOrigins  []string
	allowedHeaders  []string
	notify          chan error
	shutdownTimeout time.Duration
}

func New(mux http.Handler, opts ...Option) *Server {
	handler := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodPut,
		},
	}).Handler(mux)

	httpServer := &http.Server{
		Addr:              net.JoinHostPort("", fmt.Sprintf("%d", defaultPort)),
		Handler:           handler,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}

	s := &Server{
		Server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() {
	go func() {
		s.notify <- s.Server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
