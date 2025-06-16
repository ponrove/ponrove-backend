//go:build integration && !unit

package database_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupClickhouseDB creates a new ClickHouse connection for testing.
// It requires a CLICKHOUSE_DSN environment variable to be set.
// e.g., CLICKHOUSE_DSN="clickhouse://user:password@localhost:9000/default"
func setupClickhouseDB(t *testing.T) *database.ClickHouse {
	t.Helper()

	dsn := os.Getenv("CLICKHOUSE_DSN")
	if dsn == "" {
		dsn = "clickhouse://ponrove:321evornop@localhost:9000/default"
		// t.Skip("skipping test; CLICKHOUSE_DSN not set")
	}

	cfg := configura.NewConfigImpl()
	configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
		database.CLICKHOUSE_DSN: dsn,
	})

	// Note: The package being tested uses a global sync.Once. To ensure tests
	// are isolated, 'go test' runs this as a separate process, effectively
	// resetting the state for each run. Running multiple test packages that use
	// this database package within the same process might lead to issues.
	db, err := database.NewClickhouse(cfg)
	require.NoError(t, err)

	err = db.Ping(context.Background())
	require.NoError(t, err, "failed to connect to ClickHouse")
	return db
}

func TestNew(t *testing.T) {
	t.Run("should connect successfully with valid DSN", func(t *testing.T) {
		driver := setupClickhouseDB(t)
		defer driver.Close(context.Background())

		err := driver.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("should fail with invalid DSN", func(t *testing.T) {
		cfg := configura.NewConfigImpl()
		configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
			database.CLICKHOUSE_DSN: "invalid-dsn",
		})

		db, err := database.NewClickhouse(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("should fail with unreachable DSN", func(t *testing.T) {
		cfg := configura.NewConfigImpl()
		configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
			database.CLICKHOUSE_DSN: "clickhouse://user:password@192.0.2.1:9000/default",
		})

		db, err := database.NewClickhouse(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

func TestMigrate(t *testing.T) {
	// This test interacts with a function that uses a global sync.Once and a
	// global error variable. This makes testing the failure case difficult, as
	// a failure would poison the state for any subsequent tests within the same
	// process. Consequently, this test suite focuses on the success path to
	// validate concurrency and idempotency.

	db := setupClickhouseDB(t)
	defer db.Close(context.Background())

	migrationsDir := t.TempDir()
	migrationsURL := "file://" + migrationsDir

	up1 := `CREATE TABLE test_migration_concurrency (id UInt64) ENGINE = Memory;`
	down1 := `DROP TABLE test_migration_concurrency;`
	require.NoError(t, os.WriteFile(filepath.Join(migrationsDir, "1_initial.up.sql"), []byte(up1), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(migrationsDir, "1_initial.down.sql"), []byte(down1), 0o644))

	// Ensure the migration table is dropped at the end of the test.
	defer func() {
		err := db.NativeConn().Exec(context.Background(), down1)
		require.NoError(t, err, "failed to run teardown migration")
		err = db.NativeConn().Exec(context.Background(), `DROP TABLE IF EXISTS schema_migrations;`)
		require.NoError(t, err, "failed to drop schema_migrations table")
	}()

	t.Run("should run migrations successfully and only once when called concurrently", func(t *testing.T) {
		timer := time.Now()
		const numConcurrentCalls = 10
		var wg sync.WaitGroup
		wg.Add(numConcurrentCalls)

		errs := make(chan error, numConcurrentCalls)

		// Call Migrate concurrently from multiple goroutines.
		for i := range 10 {
			fmt.Printf("Running migration call %d - since start %s\n", i+1, time.Since(timer))
			go func() {
				defer wg.Done()
				err := db.Migrate(migrationsURL)
				errs <- err
			}()
		}

		wg.Wait()
		close(errs)

		// Check that all concurrent calls returned no error.
		for err := range errs {
			assert.NoError(t, err)
		}

		// Verify that the migration was actually applied by checking the table.
		var count uint64
		err := db.NativeConn().QueryRow(context.Background(), `SELECT count() FROM test_migration_concurrency`).Scan(&count)
		assert.NoError(t, err, "migration table should exist after running migrations")
		assert.Equal(t, uint64(0), count)
	})

	t.Run("should handle no change on subsequent calls", func(t *testing.T) {
		// This call should do nothing because the migration has already been run.
		// It should return no error, confirming idempotency.
		err := db.Migrate(migrationsURL)
		assert.NoError(t, err)
	})
}

func TestClose(t *testing.T) {
	db := setupClickhouseDB(t)

	err := db.Close(context.Background())
	require.NoError(t, err)

	/* Bug in underlying ClickHouse driver, doesn't actually close the connection.
	err = db.Ping(context.Background())
	assert.Error(t, err, "pinging a closed database should return an error")
	*/
}
