package users

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/open-feature/go-sdk/openfeature"
)

const (
	APIVersion = "1.0.0"
	APIName    = "Users API"
)

type api struct {
	openfeatureClient *openfeature.Client
}

// NewAPI creates a new instance of the Users API with the provided OpenFeature client for configuration and feature
// flag evaluation.
func NewAPI(openfeatureClient *openfeature.Client) *api {
	return &api{
		openfeatureClient: openfeatureClient,
	}
}

// NewAPIHandler creates a new HTTP handler for the Users API using the provided OpenFeature client.
func NewAPIHandler(openfeatureClient *openfeature.Client) http.Handler {
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig(APIName, APIVersion))
	huma.AutoRegister(api, NewAPI(openfeatureClient))

	return r
}

type (
	RootEndpointRequest  struct{}
	RootEndpointResponse struct {
		Status int `header:"-"`
		Body   struct {
			Message string `json:"message"`
		}
	}
)

func (a *api) RegisterRootEndpoint(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "UsersRoot",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, i *RootEndpointRequest) (*RootEndpointResponse, error) {
		return &RootEndpointResponse{
			Status: 200,
			Body: struct {
				Message string `json:"message"`
			}{
				Message: "Users API root endpoint.",
			},
		}, nil
	})
}
