package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/internal/middleware"
	"github.com/ponrove/ponrove-backend/pkg/api"
)

func New(cfg configura.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(
		// Middleware to log requests and responses.
		middleware.LogRequest,
	)

	h := humachi.New(r, huma.DefaultConfig("Ponrove Backend API", "1.0.0"))
	return r, api.Register(cfg, h)
}
