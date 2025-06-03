package testserver

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/ponrove-backend/internal/config"
	"github.com/ponrove/ponrove-backend/pkg/api"
	"github.com/ponrove/ponrove-backend/pkg/shared"
	"github.com/rs/zerolog/log"
)

// testServerConfig is a configuration struct that holds the options for the test server, and apply them after potential
// modifications from the Option functions.
type testServerConfig struct {
	captureLog    io.Writer
	serviceConfig shared.Config
	apiPackages   []func(config shared.Config, api huma.API)
}

// Option is a function that modifies the server configuration.
type Option func(*testServerConfig)

// CaptureLogToWriter is an option that allows the caller to capture the log output from the server to a writer. This is
// useful for testing the log output of the server inside the tests, if needed.
func CaptureLogToWriter(w io.Writer) Option {
	return func(cfg *testServerConfig) {
		cfg.captureLog = w
	}
}

// WithConfig allows the caller to pass in a custom configuration for the server. This is useful for testing with
// specific apis enabled.
func WithAPI(api func(config shared.Config, api huma.API)) Option {
	return func(cfg *testServerConfig) {
		if cfg.apiPackages == nil {
			cfg.apiPackages = []func(config shared.Config, api huma.API){} // Ensure api is initialized to an empty slice if not provided
		}
		cfg.apiPackages = append(cfg.apiPackages, api)
	}
}

// CreateServer creates a new test server with the provided options. The options are very useful to pass in unique
// configurations, custom mocks, or other settings.
func CreateServer(opts ...Option) *httptest.Server {
	r := chi.NewRouter()

	// Init a default server configuration, then apply any options passed in.
	cfg := &testServerConfig{
		captureLog:    nil,
		serviceConfig: config.New(),
		apiPackages:   nil,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Hook up a middleware that captures the log output from the request, if provided through the options.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := log.With().Logger().Output(nil) // suppress log output in tests

			if cfg.captureLog != nil {
				// If a writer is provided, output the log to that writer. This is useful for testing the actual log
				// output.
				logger = logger.Output(cfg.captureLog)
			}

			next.ServeHTTP(w, r.WithContext(logger.WithContext(ctx)))
		})
	})

	// Attach a new Huma API instance to the router, which will handle the API routes.
	h := humachi.New(r, huma.DefaultConfig("Ponrove Backend API", "1.0.0"))

	// If no API packages are provided, register the default API packages.
	if cfg.apiPackages == nil {
		api.Register(cfg.serviceConfig, h)
	}

	// If API packages are provided, register them with the Huma API instance.
	for _, apiPackage := range cfg.apiPackages {
		apiPackage(cfg.serviceConfig, h)
	}

	// Start a test server with the application router
	return httptest.NewServer(r)
}
