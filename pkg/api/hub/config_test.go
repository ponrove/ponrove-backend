package hub_test

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestRequiredConfiguration checks that the server package loads the configuration required by the organisations API.
func TestRequiredConfiguration(t *testing.T) {
	t.Parallel()

	cfg := config.New()
	_, driver := setupDB(t)
	err := hub.Register(
		hub.WithClickhouseDriver(driver),
	)(cfg, humachi.New(chi.NewRouter(), huma.DefaultConfig("", "")))
	assert.NoError(t, err, "failed to register users API")

	cleanCfg := configura.NewConfigImpl()
	err = hub.Register(
		hub.WithClickhouseDriver(driver),
	)(cleanCfg, humachi.New(chi.NewRouter(), huma.DefaultConfig("", "")))
	assert.Error(t, err, "expected error when registering users API with empty configuration")
	assert.ErrorIs(t, err, configura.ErrMissingVariable)
}
