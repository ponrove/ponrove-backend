package ingestion_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/pkg/shared"
	"github.com/ponrove/ponrove-backend/test/testserver"
	"github.com/stretchr/testify/suite"
)

type IngestionAPITestSuite struct {
	suite.Suite
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *IngestionAPITestSuite) TestRootEndpointFeatureFlagTrue() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}

	srv, err := testserver.CreateServer(
		testserver.WithConfig(shared.ConfigImpl{
			Bool: map[shared.Variable[bool]]bool{
				ingestion.INGESTION_API_TEST_FLAG: true,
			},
		}),
		testserver.WithAPI(ingestion.Register),
	)
	suite.NoError(err)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/ingestion/")
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

	srv, err := testserver.CreateServer(
		testserver.WithConfig(shared.ConfigImpl{
			Bool: map[shared.Variable[bool]]bool{
				ingestion.INGESTION_API_TEST_FLAG: false,
			},
		}),
		testserver.WithAPI(ingestion.Register),
	)
	suite.NoError(err)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/ingestion/")
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
