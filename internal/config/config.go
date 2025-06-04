package config

import (
	"github.com/ponrove/configura"
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
		configura.LoadEnvironment(serverConfigInstance, SERVER_OPENFEATURE_PROVIDER_NAME, "")
		configura.LoadEnvironment(serverConfigInstance, SERVER_OPENFEATURE_PROVIDER_URL, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_PORT, int64(8080))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_REQUEST_TIMEOUT, int64(30))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_SHUTDOWN_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_READ_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_WRITE_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ingestion.INGESTION_API_TEST_FLAG, false)
		configura.LoadEnvironment(serverConfigInstance, hub.HUB_API_TEST_FLAG, false)
	}

	return *serverConfigInstance
}
