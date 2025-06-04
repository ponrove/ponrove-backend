package users_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/api/users"
	"github.com/ponrove/ponrove-backend/test/testserver"
	"github.com/stretchr/testify/suite"
)

type UsersAPITestSuite struct {
	suite.Suite
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *UsersAPITestSuite) TestRootEndpointFeatureFlagTrue() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}
	cfg := configura.NewConfigImpl()
	cfg.RegBool[users.USERS_API_TEST_FLAG] = true
	srv, err := testserver.CreateServer(
		testserver.WithConfig(cfg),
		testserver.WithAPIBundle(users.Register),
	)
	suite.NoError(err)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/users/")
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Users API root endpoint.", body.Message)
	suite.True(body.TestFeatureFlag, "Expected test_feature_flag to be true")
}

// Bootstrap Test for foundational logic, this will become obsolete.
func (suite *UsersAPITestSuite) TestRootEndpointFeatureFlagFalse() {
	var body struct {
		Schema          string `json:"$schema"`
		Message         string `json:"message"`
		TestFeatureFlag bool   `json:"test_feature_flag"`
	}
	cfg := configura.NewConfigImpl()
	cfg.RegBool[users.USERS_API_TEST_FLAG] = false
	srv, err := testserver.CreateServer(
		testserver.WithConfig(cfg),
		testserver.WithAPIBundle(users.Register),
	)
	suite.NoError(err)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/users/")
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NoError(json.NewDecoder(resp.Body).Decode(&body))
	suite.Contains(resp.Header.Get("Content-Type"), "application/json")
	suite.NotEmpty(body.Schema)
	suite.Equal("Users API root endpoint.", body.Message)
	suite.False(body.TestFeatureFlag, "Expected test_feature_flag to be true")
}

func TestUsersAPITestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UsersAPITestSuite))
}
