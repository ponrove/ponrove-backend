package configuration

import "github.com/Kansuler/env"

const (
	SERVER_OPENFEATURE_PROVIDER_URL Variable = "SERVER_OPENFEATURE_PROVIDER_URL"
	SERVER_PORT                     Variable = "SERVER_PORT"
	SERVER_REQUEST_TIMEOUT          Variable = "SERVER_REQUEST_TIMEOUT"
	SERVER_SHUTDOWN_TIMEOUT         Variable = "SERVER_SHUTDOWN_TIMEOUT"
)

type serverConfig struct {
	OpenFeatureProviderURL string
	Port                   int64
	RequestTimeout         int64
	ShutdownTimeout        int64
}

var serverConfigInstance *serverConfig

// ServerConfig returns a singleton instance of the server configuration.
func ServerConfig() *serverConfig {
	if serverConfigInstance == nil {
		serverConfigInstance = &serverConfig{
			OpenFeatureProviderURL: env.String(string(SERVER_OPENFEATURE_PROVIDER_URL), ""),
			Port:                   env.Int64(string(SERVER_PORT), 8080),
			RequestTimeout:         env.Int64(string(SERVER_REQUEST_TIMEOUT), 30),
			ShutdownTimeout:        env.Int64(string(SERVER_SHUTDOWN_TIMEOUT), 10),
		}
	}

	return serverConfigInstance
}
