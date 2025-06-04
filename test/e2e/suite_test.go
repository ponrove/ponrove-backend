package e2e

import (
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	openfeatureClient *openfeature.Client
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
