package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
	"github.com/ponrove/ponrove-backend/pkg/shared"
)

// Register initializes and registers all API endpoints with the provided Huma API instance.
func Register(cfg shared.Config, api huma.API) {
	ingestion.Register(cfg, api)
	organisations.Register(cfg, api)
	users.Register(cfg, api)
}
