package organisations

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/pkg/shared"
)

const (
	ORGANISATIONS_API_TEST_FLAG shared.Variable[bool] = "ORGANISATIONS_API_TEST_FLAG" // Bootstrap flag, will become obsolete
)

type server struct {
	openfeatureClient *openfeature.Client
	config            shared.Config
}

// Register creates a new instance of the Organisations API.
func Register(cfg shared.Config, api huma.API) error {
	err := cfg.ConfigExists(
		ORGANISATIONS_API_TEST_FLAG,
	)
	if err != nil {
		return err
	}
	huma.AutoRegister(huma.NewGroup(api, "/api/organisations"), &server{
		openfeatureClient: openfeature.NewClient("organisations-api"),
		config:            cfg,
	})
	return err
}

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
		OperationID: "OrganisationsRoot",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"Organisations"},
	}, func(ctx context.Context, i *RootEndpointRequest) (*RootEndpointResponse, error) {
		testflag, err := a.openfeatureClient.BooleanValue(ctx, "test-flag", a.config.GetBool(ORGANISATIONS_API_TEST_FLAG), openfeature.EvaluationContext{})
		if err != nil {
			return nil, err
		}
		return &RootEndpointResponse{
			Status: 200,
			Body: struct {
				Message         string `json:"message"`
				TestFeatureFlag bool   `json:"test_feature_flag"`
			}{
				Message:         "Organisations API root endpoint.",
				TestFeatureFlag: testflag,
			},
		}, nil
	})
}
