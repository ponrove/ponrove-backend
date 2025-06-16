package ingestion

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/configura"
	"github.com/ponrove/octobe/driver/clickhouse"
	"github.com/ponrove/ponrove-backend/internal/database"
	"github.com/ponrove/ponrunner"
)

const (
	INGESTION_API_TEST_FLAG configura.Variable[bool] = "INGESTION_API_TEST_FLAG"
)

type server struct {
	openfeatureClient *openfeature.Client
	config            configura.Config
	clickhouse        clickhouse.Driver
}

// ingestionAPIConfig holds the configuration for the Ingestion API.
type ingestionAPIConfig struct {
	clickhouseDriver clickhouse.Driver
}

// Option is a function that modifies the api configuration.
type Option func(*ingestionAPIConfig)

// WithClickhouseDriver allows setting a custom Clickhouse driver for the ingestion API.
func WithClickhouseDriver(driver clickhouse.Driver) Option {
	return func(cfg *ingestionAPIConfig) {
		cfg.clickhouseDriver = driver
	}
}

// Register creates a new instance of the Ingestion API.
func Register(opts ...Option) ponrunner.APIBundle {
	// Init a default server configuration, then apply any options passed in.
	apiConfig := &ingestionAPIConfig{}
	for _, opt := range opts {
		opt(apiConfig)
	}

	return func(cfg configura.Config, api huma.API) error {
		err := cfg.ConfigurationKeysRegistered(
			INGESTION_API_TEST_FLAG,
		)
		if err != nil {
			return err
		}

		if apiConfig.clickhouseDriver == nil {
			clickhouseDriver, err := database.NewClickhouse(cfg)
			if err != nil {
				return err
			}

			err = clickhouseDriver.Migrate("file:///internal/database/clickhouse/migrations")
			if err != nil {
				return err
			}

			apiConfig.clickhouseDriver = clickhouseDriver
		}

		huma.AutoRegister(huma.NewGroup(api, "/api/ingestion"), &server{
			openfeatureClient: openfeature.NewClient("ingestion-api"),
			config:            cfg,
			clickhouse:        apiConfig.clickhouseDriver,
		})
		return nil
	}
}

var _ ponrunner.APIBundle = Register()

type (
	IngestionEndpointRequest  struct{}
	IngestionEndpointResponse struct {
		Status int `header:"-"`
		Body   struct {
			Message string `json:"message"`
		}
	}
)

func (a *server) RegisterPageviewEndpoint(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "Report Pageview",
		Method:      http.MethodPost,
		Path:        "/report/pageview",
		Tags:        []string{"Ingestion"},
	}, func(ctx context.Context, i *IngestionEndpointRequest) (*IngestionEndpointResponse, error) {
		// Placeholder for ingestion logic
		return &IngestionEndpointResponse{
			Status: 200,
			Body: struct {
				Message string `json:"message"`
			}{
				Message: "Pageview endpoint hit.",
			},
		}, nil
	})
}
