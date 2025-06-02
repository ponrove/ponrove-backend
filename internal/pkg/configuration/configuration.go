package configuration

import "github.com/Kansuler/env"

// Variable represents a configuration variable.
type Variable string

const (
	// Server configuration variables
	SERVER_OPENFEATURE_PROVIDER_URL Variable = "SERVER_OPENFEATURE_PROVIDER_URL"
	SERVER_PORT                     Variable = "SERVER_PORT"
	SERVER_REQUEST_TIMEOUT          Variable = "SERVER_REQUEST_TIMEOUT"
	SERVER_SHUTDOWN_TIMEOUT         Variable = "SERVER_SHUTDOWN_TIMEOUT"

	// Ingestion API configuration variables
	INGESTION_API_TEST_FLAG Variable = "INGESTION_API_TEST_FLAG" // Bootstrap flag, will become obsolete

	// Organisations API configuration variables
	ORGANISATIONS_API_TEST_FLAG Variable = "ORGANISATIONS_API_TEST_FLAG" // Bootstrap flag, will become obsolete

	// Users API configuration variables
	USERS_API_TEST_FLAG Variable = "USERS_API_TEST_FLAG" // Bootstrap flag, will become obsolete
)

type ServerConfig struct {
	ServerOpenFeatureProviderURL string
	ServerPort                   int64
	ServerRequestTimeout         int64
	ServerShutdownTimeout        int64
	IngestionApiTestFlag         bool // Bootstrap flag, will become obsolete
	OrganisationsApiTestFlag     bool // Bootstrap flag, will become obsolete
	UsersApiTestFlag             bool // Bootstrap flag, will become obsolete
}

var serverConfigInstance *ServerConfig

// ServerConfig returns a singleton instance of the server configuration.
func New() ServerConfig {
	if serverConfigInstance == nil {
		serverConfigInstance = &ServerConfig{
			ServerOpenFeatureProviderURL: env.String(string(SERVER_OPENFEATURE_PROVIDER_URL), ""),
			ServerPort:                   env.Int64(string(SERVER_PORT), 8080),
			ServerRequestTimeout:         env.Int64(string(SERVER_REQUEST_TIMEOUT), 30),
			ServerShutdownTimeout:        env.Int64(string(SERVER_SHUTDOWN_TIMEOUT), 10),
			IngestionApiTestFlag:         env.Bool(string(INGESTION_API_TEST_FLAG), false),
			OrganisationsApiTestFlag:     env.Bool(string(ORGANISATIONS_API_TEST_FLAG), false),
			UsersApiTestFlag:             env.Bool(string(USERS_API_TEST_FLAG), false),
		}
	}

	return *serverConfigInstance
}
