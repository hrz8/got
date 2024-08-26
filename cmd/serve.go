package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/hrz8/got/config"
	"github.com/hrz8/got/internal/container"
	"github.com/hrz8/got/internal/storage/postgres"
	"github.com/hrz8/got/pkg/httpserver"
	"github.com/hrz8/got/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve the server",
		RunE:  serve,
	}
	return cmd
}

func serve(cmd *cobra.Command, args []string) error {
	c := container.NewContainer()

	c.AddProviders(config.New, LogLevel, logger.New, NewDB)
	c.AddServers(NewHTTPServer)
	c.AddInvokers(func(*httpserver.Server) {})
	c.Run()

	return nil
}

func LogLevel(cfg *config.Config) logger.LogLevel {
	return logger.LogLevel(cfg.LogLevel)
}

func NewHTTPServer(lc fx.Lifecycle, logger *logger.Logger, cfg *config.Config) *httpserver.Server {
	logger.Info("registering http server", slog.Any("addr", cfg.HTTPPort))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	})

	httpServer := httpserver.New(
		mux,
		httpserver.Port(cfg.HTTPPort),
		httpserver.ShutdownTimeout(cfg.ShutdownTimeout),
		httpserver.ReadHeaderTimeout(5),
		httpserver.ReadTimeout(10),
		httpserver.WriteTimeout(10),
		httpserver.IdleTimeout(15),
		httpserver.AllowedOrigins(cfg.AllowedOrigins),
		httpserver.AllowedHeaders(cfg.AllowedHeaders),
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			httpServer.Run()
			go func() {
				err := <-httpServer.Notify()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("error starting http server", "err", err)
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := httpServer.Shutdown(); err != nil {
				logger.Error("can't shutdown http server gracefully", slog.String("err", err.Error()))
				return err
			}
			logger.Info("gracefully shutdown http server")
			return nil
		},
	})

	return httpServer
}

func NewDB(lc fx.Lifecycle, logger *logger.Logger, cfg *config.Config) *postgres.Postgres {
	pg := postgres.New(
		cfg.DatabaseURL,
		cfg.DatabaseURLReader,
		postgres.MaxOpenConnections(20),
		postgres.MaxIdleConnections(1),
		postgres.MaxConnectionLifeTime(300),
		postgres.MaxConnectionIdleTime(60),
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := pg.Connect(context.Background()); err != nil {
					logger.Error("cannot connect to database", "err", err)
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
