package flags_test

import (
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration/flags"
	"github.com/stretchr/testify/assert"
)

// TestSetOpenFeatureProvider tests the SetOpenFeatureProvider function to ensure it correctly sets the OpenFeature
// provider based on the provided configuration.
func TestSetOpenFeatureProvider(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		err      error
		cfg      configuration.ServerConfig
	}{
		{
			name:     "Default NoopProvider",
			expected: "NoopProvider",
			cfg:      configuration.ServerConfig{},
			err:      nil,
		},
		{
			name:     "Go Feature Flag Provider",
			expected: "GO Feature Flag Provider",
			cfg: configuration.ServerConfig{
				ServerOpenFeatureProviderName: "go-feature-flag",
				ServerOpenFeatureProviderURL:  "http://custom-provider.example.com",
			},
			err: nil,
		},
		{
			name: "Provider name given, missing url",
			err:  flags.ErrOpenFeatureProviderURLNotSet,
			cfg: configuration.ServerConfig{
				ServerOpenFeatureProviderName: "go-feature-flag",
			},
		},
		{
			name:     "Provider URL given, missing name",
			expected: "NoopProvider",
			cfg: configuration.ServerConfig{
				ServerOpenFeatureProviderURL: "http://custom-provider.example.com",
			},
		},
		{
			name: "Provider URL invalid",
			err:  flags.ErrInvalidOpenFeatureProviderURL,
			cfg: configuration.ServerConfig{
				ServerOpenFeatureProviderName: "go-feature-flag",
				ServerOpenFeatureProviderURL:  "http:/i\nvalid-url",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := flags.SetOpenFeatureProvider(tc.cfg)
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
