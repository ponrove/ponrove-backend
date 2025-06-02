package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// Custom response writer to capture the status code and response size, for logging.
type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Ensure the captureResponseWriter implements the http.ResponseWriter interface at compile time.
var _ http.ResponseWriter = &captureResponseWriter{}

func (crw *captureResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

func (crw *captureResponseWriter) Write(b []byte) (int, error) {
	if crw.statusCode == 0 {
		crw.statusCode = http.StatusOK
	}
	size, err := crw.ResponseWriter.Write(b)
	crw.size += size
	return size, err
}

// LogRequest is a middleware that logs the request details on each request.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Fetch a logger from context (or default logger), and make a child logger from it.
		childLogger := zerolog.Ctx(r.Context()).With().Logger()

		// Add the logger to the request context, to pass it downstream.
		r = r.WithContext(childLogger.WithContext(r.Context()))

		// Wrap the response writer to capture the status code and response size, on writes.
		crw := &captureResponseWriter{ResponseWriter: w}
		next.ServeHTTP(crw, r)

		// Fetch logger from context, potentially has new fields and contextual data added by the endpoints.
		postLogger := zerolog.Ctx(r.Context())
		postLogger.Info().
			Dur("duration", time.Since(start)).
			Str("requestMethod", r.Method).
			Str("requestUrl", r.URL.String()).
			Str("userAgent", r.Header.Get("User-Agent")).
			Str("requestSize", strconv.FormatInt(r.ContentLength, 10)).
			Str("remoteIp", r.RemoteAddr).
			Str("referer", r.Header.Get("Referer")).
			Str("protocol", r.Proto).
			Msgf("[%s][%d] %s", r.Method, crw.statusCode, r.URL.Path)
	})
}
