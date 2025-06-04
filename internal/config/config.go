package config

import (
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
)

const (
	// Server configuration variables
	SERVER_OPENFEATURE_PROVIDER_NAME configura.Variable[string] = "SERVER_OPENFEATURE_PROVIDER_NAME"
	SERVER_OPENFEATURE_PROVIDER_URL  configura.Variable[string] = "SERVER_OPENFEATURE_PROVIDER_URL"
	SERVER_PORT                      configura.Variable[int64]  = "SERVER_PORT"
	SERVER_REQUEST_TIMEOUT           configura.Variable[int64]  = "SERVER_REQUEST_TIMEOUT"
	SERVER_SHUTDOWN_TIMEOUT          configura.Variable[int64]  = "SERVER_SHUTDOWN_TIMEOUT"
)

var serverConfigInstance *configura.ConfigImpl

// ServerConfig returns a singleton instance of the server configuration.
func New() configura.Config {
	if serverConfigInstance == nil {
		serverConfigInstance = configura.NewConfigImpl()
		configura.LoadEnvironment(serverConfigInstance, SERVER_OPENFEATURE_PROVIDER_NAME, "")
		configura.LoadEnvironment(serverConfigInstance, SERVER_OPENFEATURE_PROVIDER_URL, "")
		configura.LoadEnvironment(serverConfigInstance, SERVER_PORT, int64(8080))
		configura.LoadEnvironment(serverConfigInstance, SERVER_REQUEST_TIMEOUT, int64(30))
		configura.LoadEnvironment(serverConfigInstance, SERVER_SHUTDOWN_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ingestion.INGESTION_API_TEST_FLAG, false)
		configura.LoadEnvironment(serverConfigInstance, organisations.ORGANISATIONS_API_TEST_FLAG, false)
		configura.LoadEnvironment(serverConfigInstance, users.USERS_API_TEST_FLAG, false)
	}

	return *serverConfigInstance
}
