//go:build integration && !unit

package clickhouse_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/internal/database/clickhouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a new ClickHouse connection for testing.
// It requires a CLICKHOUSE_DSN environment variable to be set.
// e.g., CLICKHOUSE_DSN="clickhouse://user:password@localhost:9000/default"
func setupTestDB(t *testing.T) *clickhouse.ClickHouse {
	t.Helper()

	dsn := os.Getenv("CLICKHOUSE_DSN")
	if dsn == "" {
		t.Skip("skipping test; CLICKHOUSE_DSN not set")
	}

	cfg := configura.NewConfigImpl()
	configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
		clickhouse.CLICKHOUSE_DSN: dsn,
	})

	db, err := clickhouse.New(cfg)
	assert.NoError(t, err)
	return db
}

func TestNew(t *testing.T) {
	t.Run("should connect successfully with valid DSN", func(t *testing.T) {
		driver := setupTestDB(t)
		defer driver.Close(context.Background())

		err := driver.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("should fail with invalid DSN", func(t *testing.T) {
		cfg := configura.NewConfigImpl()
		configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
			clickhouse.CLICKHOUSE_DSN: "invalid-dsn",
		})

		db, err := clickhouse.New(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("should fail with unreachable DSN", func(t *testing.T) {
		cfg := configura.NewConfigImpl()
		configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
			clickhouse.CLICKHOUSE_DSN: "clickhouse://user:password@192.0.2.1:9000/default",
		})

		db, err := clickhouse.New(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

func TestMigrate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close(context.Background())

	migrationsDir := t.TempDir()

	up1 := `CREATE TABLE test_migration (id UInt64) ENGINE = Memory;`
	down1 := `DROP TABLE test_migration;`
	require.NoError(t, os.WriteFile(filepath.Join(migrationsDir, "1_initial.up.sql"), []byte(up1), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(migrationsDir, "1_initial.down.sql"), []byte(down1), 0o644))

	t.Run("should run migrations successfully", func(t *testing.T) {
		migrationsURL := "file://" + migrationsDir
		err := db.Migrate(migrationsURL)
		assert.NoError(t, err)
	})

	t.Run("should handle no change", func(t *testing.T) {
		migrationsURL := "file://" + migrationsDir
		err := db.Migrate(migrationsURL)
		assert.NoError(t, err)
	})

	t.Run("should fail with invalid migrations path", func(t *testing.T) {
		migrationsURL := "file:///nonexistent/path"
		err := db.Migrate(migrationsURL)
		assert.Error(t, err)
	})

	err := db.NativeConn().Exec(context.Background(), down1)
	require.NoError(t, err)
}

func TestClose(t *testing.T) {
	db := setupTestDB(t)

	err := db.Close(context.Background())
	require.NoError(t, err)

	err = db.Ping(context.Background())
	assert.Error(t, err, "pinging a closed database should return an error")
}
