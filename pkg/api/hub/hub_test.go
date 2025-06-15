package hub_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ponrove/configura"
	"github.com/ponrove/octobe"
	"github.com/ponrove/octobe/driver/clickhouse"
	"github.com/ponrove/octobe/driver/clickhouse/mock"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/test/testserver"
	"github.com/stretchr/testify/suite"
)

func setupDB(t *testing.T) (*mock.Mock, clickhouse.Driver) {
	t.Helper()
	nativeConn := mock.NewMock()
	octdriv, err := octobe.New(clickhouse.OpenNativeWithConn(nativeConn))
	if err != nil {
		t.Fatalf("failed to create ClickHouse driver: %v", err)
	}
	return nativeConn, octdriv
}

type HubAPITestSuite struct {
	suite.Suite
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *HubAPITestSuite) TestRootEndpointFeatureFlagTrue() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}
	cfg := configura.NewConfigImpl()
	err := configura.WriteConfiguration(cfg, map[configura.Variable[bool]]bool{
		hub.HUB_API_TEST_FLAG: true,
	})
	suite.NoError(err)

	_, driver := setupDB(suite.T())
	srv, err := testserver.CreateServer(
		testserver.WithConfig(cfg),
		testserver.WithAPIBundle(hub.Register(hub.WithClickhouseDriver(driver))),
	)
	suite.NoError(err)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/hub/")
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Hub API root endpoint.", body.Message)
	suite.True(body.TestFeatureFlag, "Expected test_feature_flag to be true")
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *HubAPITestSuite) TestRootEndpointFeatureFlagFalse() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}
	cfg := configura.NewConfigImpl()
	err := configura.WriteConfiguration(cfg, map[configura.Variable[bool]]bool{
		hub.HUB_API_TEST_FLAG: false,
	})
	suite.NoError(err)

	_, driver := setupDB(suite.T())
	srv, err := testserver.CreateServer(
		testserver.WithConfig(cfg),
		testserver.WithAPIBundle(hub.Register(hub.WithClickhouseDriver(driver))),
	)
	suite.NoError(err)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/hub/")
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Hub API root endpoint.", body.Message)
	suite.False(body.TestFeatureFlag, "Expected test_feature_flag to be true")
}

func TestHubAPITestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HubAPITestSuite))
}
