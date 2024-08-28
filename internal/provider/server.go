package provider

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hrz8/got/config"
	"github.com/hrz8/got/internal/greeter"
	"github.com/hrz8/got/internal/health"
	"github.com/hrz8/got/pkg/grpcserver"
	"github.com/hrz8/got/pkg/httpserver"
	"github.com/hrz8/got/pkg/logger"
	servicev1 "github.com/hrz8/got/pkg/pb/service/v1"
	"go.uber.org/fx"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger, mux *chi.Mux, gwmux *runtime.ServeMux) *httpserver.Server {
	logger.Info("registering http server", slog.Any("port", cfg.HTTPPort))

	if gwmux == nil {
		logger.Error("error registering gateway mux")
		os.Exit(1)
	}

	mux.Mount("/api", http.StripPrefix("/api", gwmux))
	httpServer := httpserver.New(
		mux,
		httpserver.Port(cfg.HTTPPort),
		httpserver.ShutdownTimeout(cfg.ShutdownTimeout),
		httpserver.ReadHeaderTimeout(5*time.Second),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.IdleTimeout(15*time.Second),
		httpserver.AllowedOrigins(cfg.AllowedOrigins),
		httpserver.AllowedHeaders(cfg.AllowedHeaders),
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			httpServer.Run()
			go func() {
				err := <-httpServer.Notify()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("error starting http server", slog.String("err", err.Error()))
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := httpServer.Shutdown(ctx); err != nil {
				logger.Error("can't shutdown http server gracefully", slog.String("err", err.Error()))
				return err
			}
			logger.Info("gracefully shutdown http server")
			return nil
		},
	})

	return httpServer
}

func NewGRPCServer(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger) *grpcserver.Server {
	logger.Info("registering grpc server", slog.Any("port", cfg.GRPCPort))

	grpcServer := grpcserver.New(
		grpcserver.Port(cfg.GRPCPort),
	)

	servicev1.RegisterGreeterServiceServer(grpcServer.Server, greeter.NewServer())
	grpchealth.RegisterHealthServer(grpcServer.Server, health.NewServer())
	reflection.Register(grpcServer.Server)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			grpcServer.Run()
			go func() {
				err := <-grpcServer.Notify()
				if err != nil {
					logger.Error("error starting grpc server", slog.String("err", err.Error()))
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := grpcServer.Shutdown(); err != nil {
				logger.Error("can't shutdown grpc server gracefully", slog.String("err", err.Error()))
				return err
			}
			logger.Info("gracefully shutdown grpc server")
			return nil
		},
	})

	return grpcServer
}
