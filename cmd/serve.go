package cmd

import (
	"github.com/hrz8/got/config"
	"github.com/hrz8/got/internal/container"
	Middleware "github.com/hrz8/got/internal/middleware"
	"github.com/hrz8/got/internal/provider"
	User "github.com/hrz8/got/internal/user"
	"github.com/hrz8/got/pkg/grpcserver"
	"github.com/hrz8/got/pkg/httpserver"
	"github.com/hrz8/got/pkg/logger"
	"github.com/spf13/cobra"
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

	c.AddModules(
		User.Module,
		Middleware.Module,
	)
	c.AddProviders(
		config.New,
		provider.LogLevel,
		logger.New,
		provider.NewGatewayMux,
		provider.NewHTTPRouter,
		provider.NewDB,
		provider.NewGRPCClient,
	)
	c.AddServers(provider.NewHTTPServer, provider.NewProfilerServer, provider.NewGRPCServer)
	c.AddInvokers(func(*httpserver.Server) {}, func(*grpcserver.Server) {}, func(*provider.ProfilerServer) {})
	c.Run()

	return nil
}
