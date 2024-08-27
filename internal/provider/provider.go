package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hrz8/got/config"
	"github.com/hrz8/got/internal/greeter"
	"github.com/hrz8/got/internal/health"
	"github.com/hrz8/got/internal/storage/postgres"
	"github.com/hrz8/got/pkg/grpcserver"
	"github.com/hrz8/got/pkg/httpserver"
	"github.com/hrz8/got/pkg/logger"
	servicev1 "github.com/hrz8/got/pkg/pb/service/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func LogLevel(cfg *config.Config) logger.LogLevel {
	return logger.LogLevel(cfg.LogLevel)
}

func NewGRPCClient(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger) *grpc.ClientConn {
	logger.Info("registering grpc client")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	addr := net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", cfg.GRPCPort))
	cli, err := grpc.NewClient(addr, opts...)
	if err != nil {
		logger.Error("error creating grpc client", slog.String("err", err.Error()))
		os.Exit(1)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if err = cli.Close(); err != nil {
				slog.Error("failed to close grpc client", slog.Any("error", err))
				return err
			}
			logger.Info("gracefully closing grpc client")
			return nil
		},
	})

	return cli
}

func registerGatewayHandlers(cli *grpc.ClientConn) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	ctx := context.TODO()

	if err := servicev1.RegisterHealthServiceHandler(ctx, mux, cli); err != nil {
		return nil, err
	}
	if err := servicev1.RegisterGreeterServiceHandler(ctx, mux, cli); err != nil {
		return nil, err
	}

	return mux, nil
}

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, cliConn *grpc.ClientConn, logger *logger.Logger) *httpserver.Server {
	logger.Info("registering http server", slog.Any("port", cfg.HTTPPort))

	mux, err := registerGatewayHandlers(cliConn)
	if err != nil {
		logger.Error("error registering gateway", slog.String("err", err.Error()))
		os.Exit(1)
	}

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

	servicev1.RegisterHealthServiceServer(grpcServer.Server, health.NewServer())
	servicev1.RegisterGreeterServiceServer(grpcServer.Server, greeter.NewServer())
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

func NewDB(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger) *postgres.Postgres {
	pg := postgres.New(
		cfg.DatabaseURL,
		cfg.DatabaseURLReader,
		postgres.MaxOpenConnections(20),
		postgres.MaxIdleConnections(1),
		postgres.MaxConnectionLifeTime(300*time.Second),
		postgres.MaxConnectionIdleTime(60*time.Second),
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := pg.Connect(context.Background()); err != nil {
					logger.Error("cannot connect to database", slog.String("err", err.Error()))
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := pg.Close(); err != nil {
				logger.Error("failed to close the database", slog.String("err", err.Error()))
				return err
			}
			logger.Info("gracefully closing database connection")
			return nil
		},
	})

	return pg
}
