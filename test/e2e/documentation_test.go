package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ponrove/ponrove-backend/test/testserver"
	"github.com/stretchr/testify/suite"
)

type E2ETestDocumentationSuite struct {
	suite.Suite
	srv *httptest.Server
}

func (suite *E2ETestDocumentationSuite) SetupSuite() {
	var err error
	suite.srv, err = testserver.CreateServer()
	suite.NoError(err)
}

func (suite *E2ETestDocumentationSuite) TearDownSuite() {
	suite.srv.Close()
}

func (suite *E2ETestDocumentationSuite) TestOpenAPISpecificationIsServed() {
	resp, err := http.Get(suite.srv.URL + "/openapi.json")
	suite.NoError(err, "failed to get OpenAPI documentation")
	defer resp.Body.Close()
	suite.NotEmpty(resp.Body, "response body should not be empty")
}

func TestE2ETestDocumentationSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(E2ETestDocumentationSuite))
}
