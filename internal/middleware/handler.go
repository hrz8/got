package middleware

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
)

var Module = fx.Module("middleware", fx.Provide(NewHandler))

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (m *Handler) Middleware1(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fmt.Println("middleware 1", r.URL)
		next(w, r, pathParams)
	}
}

func (m *Handler) Middleware2(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fmt.Println("middleware 2", r.URL)
		next(w, r, pathParams)
	}
}
