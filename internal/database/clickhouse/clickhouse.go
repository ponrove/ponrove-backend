package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ponrove/configura"
	"github.com/ponrove/octobe"
	"github.com/ponrove/octobe/driver/clickhouse"
	"github.com/ponrove/octobe/driver/postgres"
)

const (
	CLICKHOUSE_DSN configura.Variable[string] = "CLICKHOUSE_DSN"
)

// Client defines the interface for a ClickHouse client. It's implemented by
// both the real and mock clients.
type Client interface {
	clickhouse.Driver
	NativeConn() clickhouse.NativeConn
	SQLConn() postgres.SQL
	Migrate(migrationsURL string) error
}

// ClickHouse is a wrapper around the Octobe ClickHouse driver that provides
// additional functionality, such as database migrations.
type ClickHouse struct {
	clickhouse.Driver
	nativeConn clickhouse.NativeConn
	sqlConn    *sql.DB
}

// Ensure ClickHouse implements the Client interface.
var _ Client = &ClickHouse{}

// Close closes the underlying ClickHouse native and sql connection.
func (c *ClickHouse) Close(ctx context.Context) error {
	if c.sqlConn != nil {
		err := c.sqlConn.Close()
		if err != nil {
			return fmt.Errorf("failed to close ClickHouse SQL connection: %w", err)
		}
	}

	if c.nativeConn != nil {
		err := c.nativeConn.Close()
		if err != nil {
			return fmt.Errorf("failed to close ClickHouse native connection: %w", err)
		}
	}

	return nil
}

func (c *ClickHouse) NativeConn() clickhouse.NativeConn {
	return c.nativeConn
}

func (c *ClickHouse) SQLConn() postgres.SQL {
	return c.sqlConn
}

// Migrate runs database migrations on the ClickHouse database.
func (c *ClickHouse) Migrate(migrationsURL string) error {
	driver, err := migratedb.WithInstance(c.sqlConn, &migratedb.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"clickhouse",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func New(cfg configura.Config) (*ClickHouse, error) {
	opts, err := ch.ParseDSN(cfg.String(CLICKHOUSE_DSN))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ClickHouse DSN: %w", err)
	}

	var dialCount int
	opts.DialContext = func(ctx context.Context, addr string) (net.Conn, error) {
		dialCount++
		var d net.Dialer
		return d.DialContext(ctx, "tcp", addr)
	}

	nativeConn, err := ch.Open(opts)
	if err != nil {
		return nil, err
	}

	if err := nativeConn.Ping(context.Background()); err != nil {
		nativeConn.Close()
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	sqlConn := ch.OpenDB(opts)

	if err := sqlConn.Ping(); err != nil {
		nativeConn.Close()
		return nil, fmt.Errorf("failed to open ClickHouse SQL connection: %w", err)
	}

	octdriv, err := octobe.New(clickhouse.OpenNativeWithConn(nativeConn))
	if err != nil {
		nativeConn.Close()
		sqlConn.Close()
		return nil, fmt.Errorf("failed to create Octobe ClickHouse driver: %w", err)
	}

	return &ClickHouse{
		Driver:     octdriv,
		nativeConn: nativeConn,
		sqlConn:    sqlConn,
	}, nil
}
