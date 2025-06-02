package testserver

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// serverConfig is a configuration struct that holds the options for the test server, and apply them after potential
// modifications from the Option functions.
type serverConfig struct {
	captureLog io.Writer
	mux        http.Handler
}

// Option is a function that modifies the server configuration.
type Option func(*serverConfig)

// CaptureLogToWriter is an option that allows the caller to capture the log output from the server to a writer. This is
// useful for testing the log output of the server inside the tests, if needed.
func CaptureLogToWriter(w io.Writer) Option {
	return func(cfg *serverConfig) {
		cfg.captureLog = w
	}
}

// WithMux allows the caller to provide a custom HTTP handler (mux) for the server. This is useful for testing specific
// routes or behaviors without having to set up the entire application context. If no handler is provided, it defaults
// to a NotFound handler.
func WithMux(handler func() http.Handler) Option {
	return func(cfg *serverConfig) {
		if handler != nil {
			cfg.mux = handler()
		} else {
			cfg.mux = http.NotFoundHandler() // Default to a NotFound handler if no handler is provided
		}
	}
}

// CreateServer creates a new test server with the provided options. The options are very useful to pass in unique
// configurations, custom mocks, or other settings.
func CreateServer(opts ...Option) *httptest.Server {
	// Init a default server configuration, then apply any options passed in.
	cfg := &serverConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	r := chi.NewRouter()

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

	r.Mount("/", cfg.mux)

	// Start a test server with the application router
	return httptest.NewServer(r)
}
