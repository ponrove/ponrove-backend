package ingestion_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ponrove/configura"
	"github.com/ponrove/octobe"
	"github.com/ponrove/octobe/driver/clickhouse"
	"github.com/ponrove/octobe/driver/clickhouse/mock"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	"github.com/ponrove/ponrove-backend/test/testserver"
	"github.com/stretchr/testify/suite"
)

type IngestionAPITestSuite struct {
	suite.Suite
}

func setupDB(t *testing.T) (*mock.Mock, clickhouse.Driver) {
	t.Helper()
	nativeConn := mock.NewMock()
	octdriv, err := octobe.New(clickhouse.OpenNativeWithConn(nativeConn))
	if err != nil {
		t.Fatalf("failed to create ClickHouse driver: %v", err)
	}
	return nativeConn, octdriv
}

// Test for reporting pageview endpoint
// Under development, this will change in the future.
func (suite *IngestionAPITestSuite) RegisterPageviewEndpoint() {
	var body struct {
		Schema  string `json:"$schema"`
		Message string `json:"message"`
	}

	cfg := configura.NewConfigImpl()
	err := configura.WriteConfiguration(cfg, map[configura.Variable[bool]]bool{})
	suite.NoError(err)

	_, driver := setupDB(suite.T())
	srv, err := testserver.CreateServer(
		testserver.WithConfig(cfg),
		testserver.WithAPIBundle(ingestion.Register(ingestion.WithClickhouseDriver(driver))),
	)
	suite.NoError(err)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/ingestion/report/pageview")
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Pageview endpoint hit.", body.Message)
}

func TestIngestionAPITestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IngestionAPITestSuite))
}
