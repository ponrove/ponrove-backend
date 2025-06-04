package main

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/ponrove/ponrove-backend/internal/client"
	"github.com/ponrove/ponrove-backend/internal/config"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrunner"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	cfg := config.New()

	// With configuration loaded, we can now set up the OpenFeature provider. If no provider is set, the openfeature
	// default provider will be nooped, which fallbacks to environment variables and other defaults.
	err = client.SetOpenFeatureProvider(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set openfeature provider")
	}

	// Add default logger to the context, which all http handlers derive their context (and logger) from.
	ctx := log.Logger.WithContext(context.Background())

	router := chi.NewRouter()

	// Start the runtime with the provided configuration and API bundles.
	err = ponrunner.Start(ctx, cfg, router, ingestion.Register)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start runtime")
	}
}
