package clickhouse

import (
	"context"
	"fmt"
	"net"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ponrove/configura"
	"github.com/ponrove/octobe"
	"github.com/ponrove/octobe/driver/clickhouse"
)

const (
	CLICKHOUSE_DSN configura.Variable[string] = "CLICKHOUSE_DSN"
)

func New(cfg configura.Config) (clickhouse.Driver, error) {
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

	conn, err := ch.Open(opts)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	octdriv, err := octobe.New(clickhouse.OpenNativeWithConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create Octobe ClickHouse driver: %w", err)
	}

	return octdriv, nil
}
