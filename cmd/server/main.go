package main

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/config"
	"github.com/ponrove/ponrunner"
)

func main() {
	var err error
	cfg := config.New()
	ctx := context.Background()

	router := chi.NewRouter()

	// Start the runtime with the provided configuration and API bundles.
	err = ponrunner.Start(ctx, cfg, router, func(c configura.Config, r chi.Router, a huma.API) error {
		err := ponrunner.RegisterAPIBundles(c, a, config.DefaultAPIBundles...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to start runtime", slog.Any("error", err))
	}
}
