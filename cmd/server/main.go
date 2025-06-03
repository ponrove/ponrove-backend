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

	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration/flags"
	"github.com/ponrove/ponrove-backend/internal/pkg/mux"
)

func main() {
	var err error
	cfg := configuration.New()

	// With configuration loaded, we can now set up the OpenFeature provider. If no provider is set, the openfeature
	// default provider will be nooped, which fallbacks to environment variables and other defaults.
	err = flags.SetOpenFeatureProvider(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set openfeature provider")
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
		log.Info().Msgf("starting server on %s", srv.Addr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server failed")
		} else {
			log.Info().Msg("server stopped gracefully")
		}
	case <-serverCtx.Done():
		log.Info().Msg("shutting down server...")
		stop()
	}

	shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), time.Duration(cfg.ServerShutdownTimeout)*time.Second)
	defer shutdownStop()
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Warn().Msg("shutdown context deadline exceeded, forcing shutdown")
		}
	}()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown server gracefully")
	} else {
		log.Info().Msg("server shutdown gracefully")
	}
}
