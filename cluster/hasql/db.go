package sql

import (
	"context"
	"database/sql"
)

type Querier interface {
	// Basic connection methods
	PingContext(ctx context.Context) error
	Close() error

	// Query methods with context
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Prepared statements with context
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Transaction management with context
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	Conn(ctx context.Context) (*sql.Conn, error)
}
