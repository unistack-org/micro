package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"sync"
	"time"
)

// OpenDBWithCluster creates a [*sql.DB] that uses the [ClusterQuerier]
func OpenDBWithCluster(db ClusterQuerier) (*sql.DB, error) {
	driver := NewClusterDriver(db)
	connector, err := driver.OpenConnector("")
	if err != nil {
		return nil, err
	}
	return sql.OpenDB(connector), nil
}

// ClusterDriver implements [driver.Driver] and driver.Connector for an existing [Querier]
type ClusterDriver struct {
	db ClusterQuerier
}

// NewClusterDriver creates a new [driver.Driver] that uses an existing [ClusterQuerier]
func NewClusterDriver(db ClusterQuerier) *ClusterDriver {
	return &ClusterDriver{db: db}
}

// Open implements [driver.Driver.Open]
func (d *ClusterDriver) Open(name string) (driver.Conn, error) {
	return d.Connect(context.Background())
}

// OpenConnector implements [driver.DriverContext.OpenConnector]
func (d *ClusterDriver) OpenConnector(name string) (driver.Connector, error) {
	return d, nil
}

// Connect implements [driver.Connector.Connect]
func (d *ClusterDriver) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &dbConn{conn: conn}, nil
}

// Driver implements [driver.Connector.Driver]
func (d *ClusterDriver) Driver() driver.Driver {
	return d
}

// dbConn implements driver.Conn with both context and legacy methods
type dbConn struct {
	conn *sql.Conn
	mu   sync.Mutex
}

// Prepare implements [driver.Conn.Prepare] (legacy method)
func (c *dbConn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

// PrepareContext implements [driver.ConnPrepareContext.PrepareContext]
func (c *dbConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &dbStmt{stmt: stmt}, nil
}

// Exec implements [driver.Execer.Exec] (legacy method)
func (c *dbConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	namedArgs := make([]driver.NamedValue, len(args))
	for i, value := range args {
		namedArgs[i] = driver.NamedValue{Value: value}
	}
	return c.ExecContext(context.Background(), query, namedArgs)
}

// ExecContext implements [driver.ExecerContext.ExecContext]
func (c *dbConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert driver.NamedValue to any
	interfaceArgs := make([]any, len(args))
	for i, arg := range args {
		interfaceArgs[i] = arg.Value
	}

	return c.conn.ExecContext(ctx, query, interfaceArgs...)
}

// Query implements [driver.Queryer.Query] (legacy method)
func (c *dbConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	namedArgs := make([]driver.NamedValue, len(args))
	for i, value := range args {
		namedArgs[i] = driver.NamedValue{Value: value}
	}
	return c.QueryContext(context.Background(), query, namedArgs)
}

// QueryContext implements [driver.QueryerContext.QueryContext]
func (c *dbConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert driver.NamedValue to any
	interfaceArgs := make([]any, len(args))
	for i, arg := range args {
		interfaceArgs[i] = arg.Value
	}

	rows, err := c.conn.QueryContext(ctx, query, interfaceArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

// Begin implements [driver.Conn.Begin] (legacy method)
func (c *dbConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

// BeginTx implements [driver.ConnBeginTx.BeginTx]
func (c *dbConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	sqlOpts := &sql.TxOptions{
		Isolation: sql.IsolationLevel(opts.Isolation),
		ReadOnly:  opts.ReadOnly,
	}

	tx, err := c.conn.BeginTx(ctx, sqlOpts)
	if err != nil {
		return nil, err
	}

	return &dbTx{tx: tx}, nil
}

// Ping implements [driver.Pinger.Ping]
func (c *dbConn) Ping(ctx context.Context) error {
	return c.conn.PingContext(ctx)
}

// Close implements [driver.Conn.Close]
func (c *dbConn) Close() error {
	return c.conn.Close()
}

// IsValid implements [driver.Validator.IsValid]
func (c *dbConn) IsValid() bool {
	// Ping with a short timeout to check if the connection is still valid
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return c.conn.PingContext(ctx) == nil
}

// dbStmt implements [driver.Stmt] with both context and legacy methods
type dbStmt struct {
	stmt *sql.Stmt
	mu   sync.Mutex
}

// Close implements [driver.Stmt.Close]
func (s *dbStmt) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stmt.Close()
}

// Close implements [driver.Stmt.NumInput]
func (s *dbStmt) NumInput() int {
	return -1 // Number of parameters is unknown
}

// Exec implements [driver.Stmt.Exec] (legacy method)
func (s *dbStmt) Exec(args []driver.Value) (driver.Result, error) {
	namedArgs := make([]driver.NamedValue, len(args))
	for i, value := range args {
		namedArgs[i] = driver.NamedValue{Value: value}
	}
	return s.ExecContext(context.Background(), namedArgs)
}

// ExecContext implements [driver.StmtExecContext.ExecContext]
func (s *dbStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	interfaceArgs := make([]any, len(args))
	for i, arg := range args {
		interfaceArgs[i] = arg.Value
	}
	return s.stmt.ExecContext(ctx, interfaceArgs...)
}

// Query implements [driver.Stmt.Query] (legacy method)
func (s *dbStmt) Query(args []driver.Value) (driver.Rows, error) {
	namedArgs := make([]driver.NamedValue, len(args))
	for i, value := range args {
		namedArgs[i] = driver.NamedValue{Value: value}
	}
	return s.QueryContext(context.Background(), namedArgs)
}

// QueryContext implements [driver.StmtQueryContext.QueryContext]
func (s *dbStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	interfaceArgs := make([]any, len(args))
	for i, arg := range args {
		interfaceArgs[i] = arg.Value
	}

	rows, err := s.stmt.QueryContext(ctx, interfaceArgs...)
	if err != nil {
		return nil, err
	}

	return &dbRows{rows: rows}, nil
}

// dbRows implements [driver.Rows]
type dbRows struct {
	rows *sql.Rows
}

// Columns implements [driver.Rows.Columns]
func (r *dbRows) Columns() []string {
	cols, err := r.rows.Columns()
	if err != nil {
		// This shouldn't happen if the query was successful
		return []string{}
	}
	return cols
}

// Close implements [driver.Rows.Close]
func (r *dbRows) Close() error {
	return r.rows.Close()
}

// Next implements [driver.Rows.Next]
func (r *dbRows) Next(dest []driver.Value) error {
	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return io.EOF
	}

	// Create a slice of interfaces to scan into
	scanArgs := make([]any, len(dest))
	for i := range scanArgs {
		scanArgs[i] = &dest[i]
	}

	return r.rows.Scan(scanArgs...)
}

// dbTx implements [driver.Tx]
type dbTx struct {
	tx *sql.Tx
	mu sync.Mutex
}

// Commit implements [driver.Tx.Commit]
func (t *dbTx) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tx.Commit()
}

// Rollback implements [driver.Tx.Rollback]
func (t *dbTx) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tx.Rollback()
}
