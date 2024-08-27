package provider

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/hrz8/got/config"
	"github.com/hrz8/got/internal/storage/postgres"
	"github.com/hrz8/got/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
