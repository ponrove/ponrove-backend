package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
)

// Register initializes and registers all API endpoints with the provided Huma API instance.
func Register(cfg configura.Config, api huma.API) error {
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
