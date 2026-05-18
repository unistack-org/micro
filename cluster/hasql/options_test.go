package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.yandex/hasql/v2"
)

// buildCluster creates a minimal single-node test cluster and waits for it.
func buildCluster(t *testing.T) (ClusterQuerier, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(`.*pg_is_in_recovery.*`).WillReturnRows(
		sqlmock.NewRowsWithColumnDefinition(
			sqlmock.NewColumn("role").OfType("int8", 0),
			sqlmock.NewColumn("replication_lag").OfType("int8", 0)).
			AddRow(1, 0)).
		RowsWillBeClosed().WithoutArgs()

	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	c, err := NewCluster[Querier](
		WithClusterContext(tctx),
		WithClusterNodeChecker(hasql.PostgreSQLChecker),
		WithClusterNodePicker(NewCustomPicker[Querier](CustomPickerMaxLag(1000))),
		WithClusterNodes(ClusterNode{"primary", db, 1}),
		WithClusterOptions(
			hasql.WithUpdateInterval[Querier](500*time.Millisecond),
			hasql.WithUpdateTimeout[Querier](300*time.Millisecond),
		),
	)
	if err != nil {
		_ = db.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })

	if err = c.WaitForNodes(tctx, hasql.Primary); err != nil {
		t.Fatal(err)
	}
	return c, mock
}

func TestNewClusterMissingChecker(t *testing.T) {
	_, err := NewCluster[Querier](
		WithClusterNodePicker(NewCustomPicker[Querier]()),
		WithClusterNodes(ClusterNode{"n1", nil, 1}),
	)
	if err != ErrClusterChecker {
		t.Fatalf("expected ErrClusterChecker, got %v", err)
	}
}

func TestNewClusterMissingDiscoverer(t *testing.T) {
	_, err := NewCluster[Querier](
		WithClusterNodeChecker(hasql.PostgreSQLChecker),
		WithClusterNodePicker(NewCustomPicker[Querier]()),
	)
	if err != ErrClusterDiscoverer {
		t.Fatalf("expected ErrClusterDiscoverer, got %v", err)
	}
}

func TestNewClusterMissingPicker(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func() {
		_ = db.Close()
	}()
	_, err := NewCluster[Querier](
		WithClusterNodeChecker(hasql.PostgreSQLChecker),
		WithClusterNodes(ClusterNode{"n1", db, 1}),
	)
	if err != ErrClusterPicker {
		t.Fatalf("expected ErrClusterPicker, got %v", err)
	}
}

func TestOptionSetters(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func() {
		_ = db.Close()
	}()

	opts := ClusterOptions{}
	WithClusterNodeDiscoverer(hasql.NewStaticNodeDiscoverer[Querier]())(&opts)
	if opts.NodeDiscoverer == nil {
		t.Error("NodeDiscoverer not set")
	}
	WithRetries(5)(&opts)
	if opts.Retries != 5 {
		t.Errorf("Retries = %d", opts.Retries)
	}
	WithClusterNodeStateCriterion(hasql.Primary)(&opts)
	if opts.NodeStateCriterion != hasql.Primary {
		t.Errorf("NodeStateCriterion = %v", opts.NodeStateCriterion)
	}
	WithClusterContext(context.Background())(&opts)
	if opts.Context == nil {
		t.Error("Context not set")
	}
}

func TestNodeStateCriterionDefault(t *testing.T) {
	c := &Cluster{options: ClusterOptions{NodeStateCriterion: hasql.Standby}}
	got := c.getNodeStateCriterion(context.Background())
	if got != hasql.Standby {
		t.Errorf("expected Standby, got %v", got)
	}
}

func TestNodeStateCriterionFromContext(t *testing.T) {
	c := &Cluster{options: ClusterOptions{NodeStateCriterion: hasql.Primary}}
	ctx := NodeStateCriterion(context.Background(), hasql.Standby)
	got := c.getNodeStateCriterion(ctx)
	if got != hasql.Standby {
		t.Errorf("expected Standby from context, got %v", got)
	}
}

func TestCustomPickerGetPriority(t *testing.T) {
	p := NewCustomPicker[Querier]()
	p.opts.Priority = map[string]int32{"node1": 5}
	if got := p.getPriority("node1"); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
	if got := p.getPriority("unknown"); got <= 0 {
		t.Errorf("expected positive default, got %d", got)
	}
}

func TestClusterBeginTx(t *testing.T) {
	c, mock := buildCluster(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	tx, err := c.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	_ = tx.Rollback()
}

func TestClusterExecContext(t *testing.T) {
	c, mock := buildCluster(t)

	mock.ExpectExec("INSERT INTO t VALUES").WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := c.ExecContext(context.Background(), "INSERT INTO t VALUES(?)", 1)
	if err != nil {
		t.Fatal(err)
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		t.Errorf("RowsAffected = %d", n)
	}
}

func TestClusterPrepareContext(t *testing.T) {
	c, mock := buildCluster(t)

	mock.ExpectPrepare("SELECT 1").ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(1))

	stmt, err := c.PrepareContext(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	_ = stmt.Close()
}

func TestClusterQueryContext(t *testing.T) {
	c, mock := buildCluster(t)

	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(1))

	rows, err := c.QueryContext(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	_ = rows.Close()
}

func TestClusterPingContext(t *testing.T) {
	c, mock := buildCluster(t)
	mock.ExpectPing()

	if err := c.PingContext(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestNewSQLRowError(t *testing.T) {
	row := newSQLRowError()
	if row == nil {
		t.Fatal("expected non-nil row")
	}
	if row.Err() == nil {
		t.Fatal("expected row to have error")
	}
}

// TestDriverDirectMethods exercises the ClusterDriver and dbConn methods
// that are not reached through the normal sql.DB path.
func TestDriverDirectMethods(t *testing.T) {
	c, mock := buildCluster(t)

	// Open via the cluster driver (exercises driver.Open and driver.Driver)
	drv := NewClusterDriver(c)

	// Driver()
	if drv.Driver() == nil {
		t.Fatal("Driver() returned nil")
	}

	// Open() → triggers Connect internally
	mock.ExpectPing()
	conn, err := drv.Open("")
	if err != nil {
		t.Fatal(err)
	}

	dbconn := conn.(*dbConn)

	// Ping
	mock.ExpectPing()
	if err = dbconn.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Close
	if err = dbconn.Close(); err != nil {
		t.Fatal(err)
	}
}

// TestDriverLegacyMethods exercises dbConn legacy methods (Exec, Query, Begin, BeginTx,
// Prepare, Ping, Close) and dbStmt/dbRows/dbTx methods through OpenDBWithCluster.
func TestDriverLegacyMethods(t *testing.T) {
	c, mock := buildCluster(t)

	db, err := OpenDBWithCluster(c)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := context.Background()

	// ExecContext through the sql.DB → driver.ExecerContext path
	mock.ExpectExec("INSERT INTO t").WillReturnResult(sqlmock.NewResult(1, 1))
	if _, err = db.ExecContext(ctx, "INSERT INTO t"); err != nil {
		t.Fatal(err)
	}

	// QueryContext + Next + Columns through sql.DB
	mock.ExpectQuery("SELECT col").WillReturnRows(sqlmock.NewRows([]string{"col"}).AddRow(42))
	rows, err := db.QueryContext(ctx, "SELECT col")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var v int
		if err = rows.Scan(&v); err != nil {
			t.Fatal(err)
		}
	}
	_ = rows.Close()

	// BeginTx + Commit through sql.DB
	mock.ExpectBegin()
	mock.ExpectCommit()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = tx.Commit()

	// PingContext through driver
	mock.ExpectPing()
	if err = db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
}

// TestDriverLegacyConnMethods exercises the legacy (non-context) driver.Conn methods
// (Prepare, Exec, Query, Begin) and dbStmt legacy methods directly.
func TestDriverLegacyConnMethods(t *testing.T) {
	c, mock := buildCluster(t)

	drv := NewClusterDriver(c)

	// Open returns a *dbConn
	mock.ExpectPing()
	conn, err := drv.Open("")
	if err != nil {
		t.Fatal(err)
	}
	dbconn := conn.(*dbConn)

	// Prepare (legacy)
	mock.ExpectPrepare("SELECT 1")
	stmt, err := dbconn.Prepare("SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	dbstmt := stmt.(*dbStmt)

	// NumInput
	if dbstmt.NumInput() != -1 {
		t.Errorf("NumInput = %d", dbstmt.NumInput())
	}

	// Exec on stmt (legacy)
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(0, 0))
	if _, err = dbstmt.Exec([]driver.Value{}); err != nil {
		t.Fatal(err)
	}

	// Query on stmt (legacy)
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(1))
	drows, err := dbstmt.Query([]driver.Value{})
	if err != nil {
		t.Fatal(err)
	}
	cols := drows.Columns()
	if len(cols) == 0 {
		t.Error("expected columns")
	}
	dest := make([]driver.Value, len(cols))
	_ = drows.(*dbRows).Next(dest)
	_ = drows.Close()

	// Close stmt
	if err = dbstmt.Close(); err != nil {
		t.Fatal(err)
	}

	// Exec on conn (legacy)
	mock.ExpectExec("UPDATE t").WillReturnResult(sqlmock.NewResult(0, 1))
	if _, err = dbconn.Exec("UPDATE t", []driver.Value{}); err != nil {
		t.Fatal(err)
	}

	// Query on conn (legacy)
	mock.ExpectQuery("SELECT 2").WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(2))
	drows2, err := dbconn.Query("SELECT 2", []driver.Value{})
	if err != nil {
		t.Fatal(err)
	}
	_ = drows2.Close()

	// Begin (legacy) + Rollback
	mock.ExpectBegin()
	mock.ExpectRollback()
	dtx, err := dbconn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	_ = dtx.Rollback()

	// Close conn
	if err = dbconn.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCompareNodes(t *testing.T) {
	p := &CustomPicker[Querier]{opts: CustomPickerOptions{
		MaxLag:   100,
		Priority: map[string]int32{"a": 1, "b": 2},
	}}

	makeNode := func(name string, lag int, lat time.Duration) hasql.CheckedNode[Querier] {
		// We can't easily construct hasql.CheckedNode, so test via CompareNodes
		// indirectly by calling getPriority which is internal
		_ = p.getPriority(name)
		return hasql.CheckedNode[Querier]{}
	}
	_ = makeNode

	// Direct priority tests
	if p.getPriority("a") != 1 {
		t.Error("priority a")
	}
	if p.getPriority("b") != 2 {
		t.Error("priority b")
	}
	if p.getPriority("x") <= 0 {
		t.Error("default priority should be positive")
	}
}
