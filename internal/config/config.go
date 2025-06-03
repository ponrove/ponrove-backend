package config

import (
	"github.com/Kansuler/env"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
	"github.com/ponrove/ponrove-backend/pkg/shared"
)

const (
	// Server configuration variables
	SERVER_OPENFEATURE_PROVIDER_NAME shared.Variable[string] = "SERVER_OPENFEATURE_PROVIDER_NAME"
	SERVER_OPENFEATURE_PROVIDER_URL  shared.Variable[string] = "SERVER_OPENFEATURE_PROVIDER_URL"
	SERVER_PORT                      shared.Variable[int64]  = "SERVER_PORT"
	SERVER_REQUEST_TIMEOUT           shared.Variable[int64]  = "SERVER_REQUEST_TIMEOUT"
	SERVER_SHUTDOWN_TIMEOUT          shared.Variable[int64]  = "SERVER_SHUTDOWN_TIMEOUT"
)

var serverConfigInstance *shared.ConfigImpl

// ServerConfig returns a singleton instance of the server configuration.
func New() shared.Config {
	if serverConfigInstance == nil {
		serverConfigInstance = &shared.ConfigImpl{
			Str: map[shared.Variable[string]]string{
				SERVER_OPENFEATURE_PROVIDER_URL:  env.String(string(SERVER_OPENFEATURE_PROVIDER_URL), ""),
				SERVER_OPENFEATURE_PROVIDER_NAME: env.String(string(SERVER_OPENFEATURE_PROVIDER_NAME), "NoopProvider"),
			},
			I64: map[shared.Variable[int64]]int64{
				SERVER_PORT:             env.Int64(string(SERVER_PORT), 8080),
				SERVER_REQUEST_TIMEOUT:  env.Int64(string(SERVER_REQUEST_TIMEOUT), 30),
				SERVER_SHUTDOWN_TIMEOUT: env.Int64(string(SERVER_SHUTDOWN_TIMEOUT), 10),
			},
			Bool: map[shared.Variable[bool]]bool{
				ingestion.INGESTION_API_TEST_FLAG:         env.Bool(string(ingestion.INGESTION_API_TEST_FLAG), false),
				organisations.ORGANISATIONS_API_TEST_FLAG: env.Bool(string(organisations.ORGANISATIONS_API_TEST_FLAG), false),
				users.USERS_API_TEST_FLAG:                 env.Bool(string(users.USERS_API_TEST_FLAG), false),
			},
		}
	}

	return *serverConfigInstance
}
