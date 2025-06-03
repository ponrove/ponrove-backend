package ingestion_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/ponrove/ponrove-backend/internal/pkg/configuration"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/test/testserver"
	"github.com/stretchr/testify/suite"
)

type IngestionAPITestSuite struct {
	suite.Suite
	openfeatureClient *openfeature.Client
}

func (suite *IngestionAPITestSuite) SetupTest() {
	// Initialize the OpenFeature client with a Noop provider for testing.
	openfeature.SetProvider(openfeature.NoopProvider{})
	suite.openfeatureClient = openfeature.NewClient("ingestion-test-client")
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *IngestionAPITestSuite) TestRootEndpointFeatureFlagTrue() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}

	srv := testserver.CreateServer(
		testserver.WithMux(func() http.Handler {
			return ingestion.NewAPIHandler(suite.openfeatureClient, configuration.IngestionApiConfig{
				IngestionApiTestFlag: true,
			})
		}),
	)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Ingestion API root endpoint.", body.Message)
	suite.True(body.TestFeatureFlag, "Expected test_feature_flag to be true")
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *IngestionAPITestSuite) TestRootEndpointFeatureFlagFalse() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}

	srv := testserver.CreateServer(
		testserver.WithMux(func() http.Handler {
			return ingestion.NewAPIHandler(suite.openfeatureClient, configuration.IngestionApiConfig{
				IngestionApiTestFlag: false,
			})
		}),
	)
	defer srv.Close()
	resp, err := http.Get(srv.URL)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Ingestion API root endpoint.", body.Message)
	suite.False(body.TestFeatureFlag, "Expected test_feature_flag to be true")
}

func TestIngestionAPITestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IngestionAPITestSuite))
}
