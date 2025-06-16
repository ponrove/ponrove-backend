package config

import (
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/internal/database"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrunner"
)

const (
	// Server configuration variables
	SERVER_OPENFEATURE_PROVIDER_NAME configura.Variable[string] = "SERVER_OPENFEATURE_PROVIDER_NAME"
	SERVER_OPENFEATURE_PROVIDER_URL  configura.Variable[string] = "SERVER_OPENFEATURE_PROVIDER_URL"
)

var serverConfigInstance *configura.ConfigImpl

// ServerConfig returns a singleton instance of the server configuration.
func New() configura.Config {
	if serverConfigInstance == nil {
		serverConfigInstance = configura.NewConfigImpl()
		/* Open Feature */
		configura.LoadEnvironment(serverConfigInstance, SERVER_OPENFEATURE_PROVIDER_NAME, "")
		configura.LoadEnvironment(serverConfigInstance, SERVER_OPENFEATURE_PROVIDER_URL, "")
		/* HTTP Server */
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_PORT, int64(8080))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_REQUEST_TIMEOUT, int64(30))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_SHUTDOWN_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_READ_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_WRITE_TIMEOUT, int64(10))
		/* Logging */
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_LOG_LEVEL, "info")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_LOG_FORMAT, "json")
		/* Clickhouse */
		configura.LoadEnvironment(serverConfigInstance, database.CLICKHOUSE_DSN, "")
		/* Open Telemetry */
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_LOGS_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_METRICS_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_TRACES_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_SERVICE_NAME, "ponrove-backend")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_PROTOCOL, "grpc")
		/* Ingestion API configuration */
		configura.LoadEnvironment(serverConfigInstance, ingestion.INGESTION_API_TEST_FLAG, false)
		/* Hub API configuration */
		configura.LoadEnvironment(serverConfigInstance, hub.HUB_API_TEST_FLAG, false)
	}

	return serverConfigInstance
}
