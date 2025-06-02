package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	gofeatureflag "github.com/open-feature/go-sdk-contrib/providers/go-feature-flag/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
	"github.com/ponrove/ponrove-backend/internal/pkg/middleware"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
)

func main() {
	var err error
	openfeature.SetProvider(openfeature.NoopProvider{})
	if configuration.ServerConfig().OpenFeatureProviderURL != "" {
		provider, err := gofeatureflag.NewProvider(
			gofeatureflag.ProviderOptions{
				Endpoint: configuration.ServerConfig().OpenFeatureProviderURL,
			},
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize OpenFeature provider")
		}

		openfeature.SetProviderAndWait(provider)
	}

	// Add default logger to the context, which all http handlers derive their context (and logger) from.
	serverCtx := log.Logger.WithContext(context.Background())

	// Setup context to listen for SIGINT and SIGKILL signals.
	serverCtx, stop := signal.NotifyContext(serverCtx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", configuration.ServerConfig().Port),
		// Use the context that includes a notify channel for graceful shutdown.
		BaseContext:  func(_ net.Listener) context.Context { return serverCtx },
		ReadTimeout:  time.Second,
		WriteTimeout: time.Duration(configuration.ServerConfig().RequestTimeout) * time.Second,
		Handler:      createMux(),
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Info().Msgf("Starting server on %s", srv.Addr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server failed")
		} else {
			log.Info().Msg("Server stopped gracefully")
		}
	case <-serverCtx.Done():
		log.Info().Msg("Shutting down server...")
		stop()
	}

	shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), time.Duration(configuration.ServerConfig().ShutdownTimeout)*time.Second)
	defer shutdownStop()
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Warn().Msg("Shutdown context deadline exceeded, forcing shutdown")
		}
	}()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown server gracefully")
	} else {
		log.Info().Msg("Server shutdown gracefully")
	}
}

// createMux initializes the HTTP router and registers all API endpoints.
func createMux() http.Handler {
	r := chi.NewRouter()

	r.Use(
		// Middleware to log requests and responses.
		middleware.LogRequest,
	)

	api := humachi.New(r, huma.DefaultConfig("Ponrove Backend API", "1.0.0"))

	openfeatureClient := openfeature.NewClient("ponrove_backend")
	openfeature.SetEvaluationContext(openfeature.NewEvaluationContext("general", map[string]any{}))

	// Ingestion API will handle all requests related to data ingestion and processing from clients.
	huma.AutoRegister(huma.NewGroup(api, "/api/ingestion"), ingestion.NewAPI(openfeatureClient))

	// Organisations API will handle all requests related to organisations, such as creating, updating,
	// deleting, and retrieving organisation information.
	huma.AutoRegister(huma.NewGroup(api, "/api/organisations"), organisations.NewAPI(openfeatureClient))

	// Users API will handle all requests related to user management, such as creating, updating,
	// deleting, and retrieving user information.
	huma.AutoRegister(huma.NewGroup(api, "/api/users"), users.NewAPI(openfeatureClient))

	return r
}
