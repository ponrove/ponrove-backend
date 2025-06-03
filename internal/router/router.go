package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/ponrove-backend/internal/middleware"
	"github.com/ponrove/ponrove-backend/pkg/api"
	"github.com/ponrove/ponrove-backend/pkg/shared"
)

func New(cfg shared.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(
		// Middleware to log requests and responses.
		middleware.LogRequest,
	)

	h := humachi.New(r, huma.DefaultConfig("Ponrove Backend API", "1.0.0"))
	return r, api.Register(cfg, h)
}
