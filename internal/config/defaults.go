package config

import (
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrunner"
)

var DefaultAPIBundles = []ponrunner.APIBundle{
	ingestion.Register,
	hub.Register,
}
