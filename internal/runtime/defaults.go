package runtime

import (
	"github.com/ponrove/ponrove-backend/pkg/api"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/api/organisations"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
)

var DefaultAPIBundles = []api.APIBundle{
	ingestion.Register,
	organisations.Register,
	users.Register,
}
