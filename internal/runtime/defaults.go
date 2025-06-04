package runtime

import (
	"github.com/ponrove/ponrove-backend/pkg/api"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
)

var DefaultAPIBundles = []api.APIBundle{
	ingestion.Register,
	hub.Register,
}
