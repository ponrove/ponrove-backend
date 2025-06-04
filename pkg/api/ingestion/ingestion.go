package ingestion

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrunner"
)

const (
	INGESTION_API_TEST_FLAG configura.Variable[bool] = "INGESTION_API_TEST_FLAG"
)

type server struct {
	openfeatureClient *openfeature.Client
	config            configura.Config
}

// Register creates a new instance of the Ingestion API.
func Register(cfg configura.Config, api huma.API) error {
	err := cfg.ConfigurationKeysRegistered(
		INGESTION_API_TEST_FLAG,
	)
	if err != nil {
		return err
	}
	huma.AutoRegister(huma.NewGroup(api, "/api/ingestion"), &server{
		openfeatureClient: openfeature.NewClient("ingestion-api"),
		config:            cfg,
	})
	return nil
}

var _ ponrunner.APIBundle = Register

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
		OperationID: "IngestionRoot",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"Ingestion"},
	}, func(ctx context.Context, i *RootEndpointRequest) (*RootEndpointResponse, error) {
		testflag, err := a.openfeatureClient.BooleanValue(ctx, "test-flag", a.config.Bool(INGESTION_API_TEST_FLAG), openfeature.EvaluationContext{})
		if err != nil {
			return nil, err
		}
		return &RootEndpointResponse{
			Status: 200,
			Body: struct {
				Message         string `json:"message"`
				TestFeatureFlag bool   `json:"test_feature_flag"`
			}{
				Message:         "Ingestion API root endpoint.",
				TestFeatureFlag: testflag,
			},
		}, nil
	})
}
