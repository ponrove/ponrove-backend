package testserver

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/internal/config"
	"github.com/ponrove/ponrove-backend/internal/runtime"
	"github.com/ponrove/ponrove-backend/pkg/api"
	"github.com/rs/zerolog/log"
)

// testServerConfig is a configuration struct that holds the options for the test server, and apply them after potential
// modifications from the Option functions.
type testServerConfig struct {
	captureLog    io.Writer
	serviceConfig configura.Config
	apiBundles    []api.APIBundle
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

// WithAPIBundle allows the caller to pass in a custom API bundle to the server. If one API bundle is provided, it will
// overwrite the default API bundles. If multiple API bundles are provided, they will be appended.
func WithAPIBundle(bundle api.APIBundle) Option {
	return func(cfg *testServerConfig) {
		if cfg.apiBundles == nil {
			cfg.apiBundles = []api.APIBundle{} // Ensure api is initialized to an empty slice if not provided
		}
		cfg.apiBundles = append(cfg.apiBundles, bundle)
	}
}

// WithConfig allows the caller to pass in a custom configuration for the server from an existing configura.Config.
func WithConfig(cfg configura.Config) Option {
	return func(tsc *testServerConfig) {
		tsc.serviceConfig = cfg
	}
}

// CreateServer creates a new test server with the provided options. The options are very useful to pass in unique
// configurations, custom mocks, or other settings.
func CreateServer(opts ...Option) (*httptest.Server, error) {
	r := chi.NewRouter()

	// Init a default server configuration, then apply any options passed in.
	cfg := &testServerConfig{
		captureLog:    nil,
		serviceConfig: config.New(),
		apiBundles:    nil,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// If no bundles are provided, use the default API bundles.
	if cfg.apiBundles == nil {
		cfg.apiBundles = runtime.DefaultAPIBundles
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

	// If API packages are provided, register them with the Huma API instance.
	api.RegisterAPIBundles(cfg.serviceConfig, h, cfg.apiBundles...)

	// Start a test server with the application router
	return httptest.NewServer(r), nil
}
