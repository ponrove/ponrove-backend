package main

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/config"
	"github.com/ponrove/ponrunner"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	cfg := config.New()

	// Add default logger to the context, which all http handlers derive their context (and logger) from.
	ctx := log.Logger.WithContext(context.Background())

	router := chi.NewRouter()

	// Start the runtime with the provided configuration and API bundles.
	err = ponrunner.Start(ctx, cfg, router, hub.Register)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start runtime")
	}
}
