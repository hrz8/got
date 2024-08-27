package user

import (
	"fmt"
	"net/http"

	"go.uber.org/fx"
)

var Module = fx.Module("user", fx.Provide(NewHandler))

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User List")
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Fprint(w, "User ID: "+pathParams["id"])
}
