package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/hrz8/got/config"
	"github.com/hrz8/got/pkg/logger"
	"go.uber.org/fx"
)

type ProfilerServer struct {
	*http.Server
}

func NewProfilerServer(lc fx.Lifecycle, cfg *config.Config, logger *logger.Logger) *ProfilerServer {
	if !cfg.EnableProfiler {
		return nil
	}

	logger.Info("registering pprof server", slog.Any("port", cfg.ProfilerPort))

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))

	pprofserver := &ProfilerServer{
		Server: &http.Server{
			Addr:         net.JoinHostPort("", fmt.Sprintf("%d", cfg.ProfilerPort)),
			Handler:      mux,
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 20 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := pprofserver.ListenAndServe(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						logger.Error("error starting pprof server", slog.String("err", err.Error()))
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := pprofserver.Server.Shutdown(timeoutCtx); err != nil {
				logger.Error("can't shutdown http pprof gracefully", slog.String("err", err.Error()))
			}
			return nil
		},
	})

	return pprofserver
}
