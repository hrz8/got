package provider

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	Middleware "github.com/hrz8/got/internal/middleware"
	User "github.com/hrz8/got/internal/user"
	servicev1 "github.com/hrz8/got/pkg/pb/service/v1"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
)

func NewHTTPRouter(
	user *User.Handler,
) *chi.Mux {
	// chi routers
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/", user.List)
	})

	return r
}

type customRoute struct {
	meth    string
	pattern string
}

func NewGatewayMux(
	cli *grpc.ClientConn,
	middleware *Middleware.Handler,
	user *User.Handler,
) *runtime.ServeMux {
	healthClient := grpchealth.NewHealthClient(cli)
	opts := []runtime.ServeMuxOption{
		runtime.WithHealthzEndpoint(healthClient), // healthz will not injected with middleware
		runtime.WithMiddlewares(middleware.Middleware1, middleware.Middleware2),
	}

	mux := runtime.NewServeMux(opts...)
	ctx := context.TODO()

	// grpc services
	if err := servicev1.RegisterGreeterServiceHandler(ctx, mux, cli); err != nil {
		return nil
	}

	// extras routes
	if err := mux.HandlePath("GET", "/users/{id}", user.Detail); err != nil {
		return nil
	}

	return mux
}
