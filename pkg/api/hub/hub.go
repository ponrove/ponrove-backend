package hub

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/configura"
	ch "github.com/ponrove/octobe/driver/clickhouse"
	"github.com/ponrove/ponrove-backend/internal/database/clickhouse"
	"github.com/ponrove/ponrunner"
)

const (
	HUB_API_TEST_FLAG configura.Variable[bool] = "HUB_API_TEST_FLAG" // Bootstrap flag, will become obsolete
)

type server struct {
	openfeatureClient *openfeature.Client
	config            configura.Config
	clickhouse        ch.Driver
}

// ingestionAPIConfig holds the configuration for the Ingestion API.
type hubAPIConfig struct {
	clickhouseDriver ch.Driver
}

// Option is a function that modifies the api configuration.
type Option func(*hubAPIConfig)

// WithClickhouseDriver allows setting a custom Clickhouse driver for the ingestion API.
func WithClickhouseDriver(driver ch.Driver) Option {
	return func(cfg *hubAPIConfig) {
		cfg.clickhouseDriver = driver
	}
}

// Register creates a new instance of the Hub API.
func Register(opts ...Option) ponrunner.APIBundle {
	// Init a default server configuration, then apply any options passed in.
	apiConfig := &hubAPIConfig{}
	for _, opt := range opts {
		opt(apiConfig)
	}

	return func(cfg configura.Config, api huma.API) error {
		err := cfg.ConfigurationKeysRegistered(
			HUB_API_TEST_FLAG,
		)
		if err != nil {
			return err
		}

		if apiConfig.clickhouseDriver == nil {
			clickhouseDriver, err := clickhouse.New(cfg)
			if err != nil {
				return err
			}
			apiConfig.clickhouseDriver = clickhouseDriver
		}

		huma.AutoRegister(huma.NewGroup(api, "/api/hub"), &server{
			openfeatureClient: openfeature.NewClient("hub-api"),
			config:            cfg,
		})
		return err
	}
}

var _ ponrunner.APIBundle = Register()

type (
	RootEndpointRequest  struct{}
	RootEndpointResponse struct {
		Status int `header:"-"`
		Body   struct {
			Message         string `json:"message"`
			TestFeatureFlag bool   `json:"test_feature_flag"`
		}
	}
)

// Bootstrap endpoint for foundational logic, this will become obsolete.
func (a *server) RegisterRootEndpoint(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "HubRoot",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"Hub"},
	}, func(ctx context.Context, i *RootEndpointRequest) (*RootEndpointResponse, error) {
		testflag, err := a.openfeatureClient.BooleanValue(ctx, "test-flag", a.config.Bool(HUB_API_TEST_FLAG), openfeature.EvaluationContext{})
		if err != nil {
			return nil, err
		}
		return &RootEndpointResponse{
			Status: 200,
			Body: struct {
				Message         string `json:"message"`
				TestFeatureFlag bool   `json:"test_feature_flag"`
			}{
				Message:         "Hub API root endpoint.",
				TestFeatureFlag: testflag,
			},
		}, nil
	})
}
