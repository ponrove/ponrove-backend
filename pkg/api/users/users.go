package users

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
)

const (
	APIVersion = "1.0.0"
	APIName    = "Users API"
)

type api struct {
	openfeatureClient *openfeature.Client
	config            configuration.ServerConfig
}

// NewAPI creates a new instance of the Users API with the provided OpenFeature client for configuration and feature
// flag evaluation.
func NewAPI(openfeatureClient *openfeature.Client, cfg configuration.ServerConfig) *api {
	return &api{
		openfeatureClient: openfeatureClient,
		config:            cfg,
	}
}

// NewAPIHandler creates a new HTTP handler for the Users API using the provided OpenFeature client.
func NewAPIHandler(openfeatureClient *openfeature.Client, cfg configuration.ServerConfig) http.Handler {
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig(APIName, APIVersion))
	huma.AutoRegister(api, NewAPI(openfeatureClient, cfg))

	return r
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
func (a *api) RegisterRootEndpoint(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "UsersRoot",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, i *RootEndpointRequest) (*RootEndpointResponse, error) {
		testflag, err := a.openfeatureClient.BooleanValue(ctx, "test-flag", a.config.UsersApiTestFlag, openfeature.EvaluationContext{})
		if err != nil {
			return nil, err
		}
		return &RootEndpointResponse{
			Status: 200,
			Body: struct {
				Message         string `json:"message"`
				TestFeatureFlag bool   `json:"test_feature_flag"`
			}{
				Message:         "Users API root endpoint.",
				TestFeatureFlag: testflag,
			},
		}, nil
	})
}
