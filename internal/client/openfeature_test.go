package client_test

import (
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/client"
	"github.com/ponrove/ponrove-backend/internal/config"
	"github.com/ponrove/ponrove-backend/pkg/shared"
	"github.com/stretchr/testify/assert"
)

type TestSetOpenFeatureProviderTestCase struct {
	name     string
	expected string
	err      error
	cfg      shared.Config
}

var TestSetOpenFeatureProviderTestCases = []TestSetOpenFeatureProviderTestCase{
	{
		name:     "Default NoopProvider",
		expected: "NoopProvider",
		cfg:      shared.ConfigImpl{},
		err:      nil,
	},
	{
		name:     "NoopProvider without URL",
		expected: "NoopProvider",
		cfg: &shared.ConfigImpl{
			Str: map[shared.Variable[string]]string{
				config.SERVER_OPENFEATURE_PROVIDER_NAME: "NoopProvider",
			},
		},
	},
	{
		name:     "Go Feature Flag Provider",
		expected: "GO Feature Flag Provider",
		cfg: &shared.ConfigImpl{
			Str: map[shared.Variable[string]]string{
				config.SERVER_OPENFEATURE_PROVIDER_NAME: "go-feature-flag",
				config.SERVER_OPENFEATURE_PROVIDER_URL:  "http://custom-provider.example.com",
			},
		},
		err: nil,
	},
	{
		name: "Provider name given, missing url",
		err:  client.ErrOpenFeatureProviderURLNotSet,
		cfg: &shared.ConfigImpl{
			Str: map[shared.Variable[string]]string{
				config.SERVER_OPENFEATURE_PROVIDER_NAME: "go-feature-flag",
			},
		},
	},
	{
		name:     "Provider URL given, missing name",
		expected: "NoopProvider",
		cfg: &shared.ConfigImpl{
			Str: map[shared.Variable[string]]string{
				config.SERVER_OPENFEATURE_PROVIDER_URL: "http://custom-provider.example.com",
			},
		},
	},
	{
		name: "Provider URL invalid",
		err:  client.ErrInvalidOpenFeatureProviderURL,
		cfg: &shared.ConfigImpl{
			Str: map[shared.Variable[string]]string{
				config.SERVER_OPENFEATURE_PROVIDER_NAME: "go-feature-flag",
				config.SERVER_OPENFEATURE_PROVIDER_URL:  "http:/i\nvalid-url",
			},
		},
	},
}

// TestSetOpenFeatureProvider tests the SetOpenFeatureProvider function to ensure it correctly sets the OpenFeature
// provider based on the provided configuration.
func TestSetOpenFeatureProvider(t *testing.T) {
	for _, tc := range TestSetOpenFeatureProviderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.SetOpenFeatureProvider(tc.cfg)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}

			assert.NoError(t, err)
			metadata := openfeature.NamedProviderMetadata("")
			assert.Equal(t, tc.expected, metadata.Name)

			openfeatureClient := openfeature.NewClient("test")
			assert.Equal(t, "test", openfeatureClient.Metadata().Domain())
		})
	}
}
