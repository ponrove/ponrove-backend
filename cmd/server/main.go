package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	gofeatureflag "github.com/open-feature/go-sdk-contrib/providers/go-feature-flag/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
	"github.com/ponrove/ponrove-backend/internal/pkg/mux"
)

func main() {
	var err error
	cfg := configuration.New()
	openfeature.SetProvider(openfeature.NoopProvider{})
	if cfg.ServerOpenFeatureProviderURL != "" {
		provider, err := gofeatureflag.NewProvider(
			gofeatureflag.ProviderOptions{
				Endpoint: cfg.ServerOpenFeatureProviderURL,
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
		Addr: fmt.Sprintf(":%d", cfg.ServerPort),
		// Use the context that includes a notify channel for graceful shutdown.
		BaseContext:  func(_ net.Listener) context.Context { return serverCtx },
		ReadTimeout:  time.Second,
		WriteTimeout: time.Duration(cfg.ServerRequestTimeout) * time.Second,
		Handler:      mux.New(cfg),
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

	shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), time.Duration(cfg.ServerShutdownTimeout)*time.Second)
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
