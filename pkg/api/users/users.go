package users

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/configura"
)

const (
	USERS_API_TEST_FLAG configura.Variable[bool] = "USERS_API_TEST_FLAG" // Bootstrap flag, will become obsolete
)

type UsersApiConfig struct {
	UsersApiTestFlag bool // Bootstrap flag, will become obsolete
}

type server struct {
	openfeatureClient *openfeature.Client
	config            configura.Config
}

// Register creates a new instance of the Users API.
func Register(cfg configura.Config, api huma.API) error {
	err := cfg.ConfigurationKeysRegistered(
		USERS_API_TEST_FLAG,
	)
	if err != nil {
		return err
	}

	huma.AutoRegister(huma.NewGroup(api, "/api/users"), &server{
		openfeatureClient: openfeature.NewClient("users-api"),
		config:            cfg,
	})
	return nil
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
		OperationID: "UsersRoot",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, i *RootEndpointRequest) (*RootEndpointResponse, error) {
		testflag, err := a.openfeatureClient.BooleanValue(ctx, "test-flag", a.config.Bool(USERS_API_TEST_FLAG), openfeature.EvaluationContext{})
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
