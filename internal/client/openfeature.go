package client

import (
	"errors"
	"fmt"
	"net/url"

	gofeatureflag "github.com/open-feature/go-sdk-contrib/providers/go-feature-flag/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/config"
)

var (
	ErrOpenFeatureProviderURLNotSet   = errors.New("openfeature provider url not set")
	ErrUnsupportedOpenFeatureProvider = errors.New("unsupported openfeature provider")
	ErrInvalidOpenFeatureProviderURL  = errors.New("invalid openfeature provider url")
)

// SetOpenFeatureProvider initializes the OpenFeature provider based on the server configuration. It's possible to add
// more providers in the future, but for now we only support the Go Feature Flag provider.
func SetOpenFeatureProvider(cfg configura.Config) error {
	openfeature.SetProvider(openfeature.NoopProvider{})
	if cfg.String(config.SERVER_OPENFEATURE_PROVIDER_NAME) == "" || cfg.String(config.SERVER_OPENFEATURE_PROVIDER_NAME) == "NoopProvider" {
		return nil // No provider configured, using noop provider.
	}

	// If the provider URL is not set, we cannot initialize the provider. Return an error to indicate this.
	if cfg.String(config.SERVER_OPENFEATURE_PROVIDER_URL) == "" {
		return fmt.Errorf("%w: %s", ErrOpenFeatureProviderURLNotSet, cfg.String(config.SERVER_OPENFEATURE_PROVIDER_URL))
	}
	var err error

	// parse url
	_, err = url.ParseRequestURI(cfg.String(config.SERVER_OPENFEATURE_PROVIDER_URL))
	if err != nil {
		return fmt.Errorf("%w: %s: %v", ErrInvalidOpenFeatureProviderURL, cfg.String(config.SERVER_OPENFEATURE_PROVIDER_URL), err)
	}

	var provider openfeature.FeatureProvider
	switch cfg.String(config.SERVER_OPENFEATURE_PROVIDER_NAME) {
	case "go-feature-flag":
		// Currently, error can only occur if the URL is empty, which is handled above. Surpress the error here.
		provider, _ = gofeatureflag.NewProvider(
			gofeatureflag.ProviderOptions{
				Endpoint: cfg.String(config.SERVER_OPENFEATURE_PROVIDER_URL),
			},
		)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedOpenFeatureProvider, cfg.String(config.SERVER_OPENFEATURE_PROVIDER_NAME))
	}

	return openfeature.SetProviderAndWait(provider)
}
