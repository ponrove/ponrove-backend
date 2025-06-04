package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/ponrove/configura"
)

type APIBundle func(configura.Config, huma.API) error

func RegisterAPIBundles(cfg configura.Config, api huma.API, bundles ...APIBundle) error {
	for _, bundle := range bundles {
		if err := bundle(cfg, api); err != nil {
			return err
		}
	}

	return nil
}
