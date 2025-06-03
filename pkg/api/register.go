package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
	"github.com/ponrove/ponrove-backend/pkg/shared"
)

// Register initializes and registers all API endpoints with the provided Huma API instance.
func Register(cfg shared.Config, api huma.API) error {
	var err error
	err = ingestion.Register(cfg, api)
	if err != nil {
		return err
	}
	err = organisations.Register(cfg, api)
	if err != nil {
		return err
	}
	err = users.Register(cfg, api)
	if err != nil {
		return err
	}
	return nil
}
