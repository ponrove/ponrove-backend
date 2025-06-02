package mux

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
	"github.com/ponrove/ponrove-backend/internal/pkg/middleware"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
)

// New initializes the HTTP router and registers all API endpoints.
func New(cfg configuration.ServerConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(
		// Middleware to log requests and responses.
		middleware.LogRequest,
	)

	api := humachi.New(r, huma.DefaultConfig("Ponrove Backend API", "1.0.0"))

	openfeatureClient := openfeature.NewClient("ponrove_backend")
	openfeature.SetEvaluationContext(openfeature.NewEvaluationContext("general", map[string]any{}))

	// Ingestion API will handle all requests related to data ingestion and processing from clients.
	huma.AutoRegister(huma.NewGroup(api, "/api/ingestion"), ingestion.NewAPI(openfeatureClient, cfg))

	// Organisations API will handle all requests related to organisations, such as creating, updating,
	// deleting, and retrieving organisation information.
	huma.AutoRegister(huma.NewGroup(api, "/api/organisations"), organisations.NewAPI(openfeatureClient, cfg))

	// Users API will handle all requests related to user management, such as creating, updating,
	// deleting, and retrieving user information.
	huma.AutoRegister(huma.NewGroup(api, "/api/users"), users.NewAPI(openfeatureClient, cfg))

	return r
}
