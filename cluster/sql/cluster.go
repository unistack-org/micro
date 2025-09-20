package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"

	"golang.yandex/hasql/v2"
)

var errNoAliveNodes = errors.New("no alive nodes")

func newSQLRowError() *sql.Row {
	row := &sql.Row{}
	t := reflect.TypeOf(row).Elem()
	field, _ := t.FieldByName("err")
	rowPtr := unsafe.Pointer(row)
	errFieldPtr := unsafe.Pointer(uintptr(rowPtr) + field.Offset)
	errPtr := (*error)(errFieldPtr)
	*errPtr = errNoAliveNodes
	return row
}

type ClusterQuerier interface {
	Querier
	WaitForNodes(ctx context.Context, criterion ...hasql.NodeStateCriterion) error
}

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

	// Connection pool management
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	Stats() sql.DBStats

	Conn(ctx context.Context) (*sql.Conn, error)
}

var (
	ErrClusterChecker    = errors.New("cluster node checker required")
	ErrClusterDiscoverer = errors.New("cluster node discoverer required")
	ErrClusterPicker     = errors.New("cluster node picker required")
)

type Cluster struct {
	hasql   *hasql.Cluster[Querier]
	options ClusterOptions
}

// NewCluster returns Querier that provides cluster of nodes
func NewCluster[T Querier](opts ...ClusterOption) (ClusterQuerier, error) {
	options := ClusterOptions{Context: context.Background()}
	for _, opt := range opts {
		opt(&options)
	}
	if options.NodeChecker == nil {
		return nil, ErrClusterChecker
	}
	if options.NodeDiscoverer == nil {
		return nil, ErrClusterDiscoverer
	}
	if options.NodePicker == nil {
		return nil, ErrClusterPicker
	}

	if options.Retries < 1 {
		options.Retries = 1
	}

	if options.NodeStateCriterion == 0 {
		options.NodeStateCriterion = hasql.Primary
	}

	options.Options = append(options.Options, hasql.WithNodePicker(options.NodePicker))
	if p, ok := options.NodePicker.(*CustomPicker[Querier]); ok {
		p.opts.Priority = options.NodePriority
	}

	c, err := hasql.NewCluster(
		options.NodeDiscoverer,
		options.NodeChecker,
		options.Options...,
	)
	if err != nil {
		return nil, err
	}

	return &Cluster{hasql: c, options: options}, nil
}

// compile time guard
var _ hasql.NodePicker[Querier] = (*CustomPicker[Querier])(nil)

type nodeStateCriterionKey struct{}

// NodeStateCriterion inject hasql.NodeStateCriterion to context
func NodeStateCriterion(ctx context.Context, c hasql.NodeStateCriterion) context.Context {
	return context.WithValue(ctx, nodeStateCriterionKey{}, c)
}

// CustomPickerOptions holds options to pick nodes
type CustomPickerOptions struct {
	MaxLag   int
	Priority map[string]int32
	Retries  int
}

// CustomPickerOption func apply option to CustomPickerOptions
type CustomPickerOption func(*CustomPickerOptions)

// CustomPickerMaxLag specifies max lag for which node can be used
func CustomPickerMaxLag(n int) CustomPickerOption {
	return func(o *CustomPickerOptions) {
		o.MaxLag = n
	}
}

// NewCustomPicker creates new node picker
func NewCustomPicker[T Querier](opts ...CustomPickerOption) *CustomPicker[Querier] {
	options := CustomPickerOptions{}
	for _, o := range opts {
		o(&options)
	}
	return &CustomPicker[Querier]{opts: options}
}

// CustomPicker holds node picker options
type CustomPicker[T Querier] struct {
	opts CustomPickerOptions
}

// PickNode used to return specific node
func (p *CustomPicker[T]) PickNode(cnodes []hasql.CheckedNode[T]) hasql.CheckedNode[T] {
	for _, n := range cnodes {
		fmt.Printf("node %s\n", n.Node.String())
	}
	return cnodes[0]
}

func (p *CustomPicker[T]) getPriority(nodeName string) int32 {
	if prio, ok := p.opts.Priority[nodeName]; ok {
		return prio
	}
	return math.MaxInt32 // Default to lowest priority
}

// CompareNodes used to sort nodes
func (p *CustomPicker[T]) CompareNodes(a, b hasql.CheckedNode[T]) int {
	fmt.Printf("CompareNodes %s %s\n", a.Node.String(), b.Node.String())
	// Get replication lag values
	aLag := a.Info.(interface{ ReplicationLag() int }).ReplicationLag()
	bLag := b.Info.(interface{ ReplicationLag() int }).ReplicationLag()

	// First check that lag lower then MaxLag
	if aLag > p.opts.MaxLag && bLag > p.opts.MaxLag {
		fmt.Printf("CompareNodes aLag > p.opts.MaxLag && bLag > p.opts.MaxLag\n")
		return 0 // both are equal
	}

	// If one node exceeds MaxLag and the other doesn't, prefer the one that doesn't
	if aLag > p.opts.MaxLag {
		fmt.Printf("CompareNodes aLag > p.opts.MaxLag\n")
		return 1 // b is better
	}
	if bLag > p.opts.MaxLag {
		fmt.Printf("CompareNodes bLag > p.opts.MaxLag\n")
		return -1 // a is better
	}

	// Get node priorities
	aPrio := p.getPriority(a.Node.String())
	bPrio := p.getPriority(b.Node.String())

	// if both priority equals
	if aPrio == bPrio {
		fmt.Printf("CompareNodes aPrio == bPrio\n")
		// First compare by replication lag
		if aLag < bLag {
			fmt.Printf("CompareNodes aLag < bLag\n")
			return -1
		}
		if aLag > bLag {
			fmt.Printf("CompareNodes aLag > bLag\n")
			return 1
		}
		// If replication lag is equal, compare by latency
		aLatency := a.Info.(interface{ Latency() time.Duration }).Latency()
		bLatency := b.Info.(interface{ Latency() time.Duration }).Latency()

		if aLatency < bLatency {
			return -1
		}
		if aLatency > bLatency {
			return 1
		}

		// If lag and latency is equal
		return 0
	}

	// If priorities are different, prefer the node with lower priority value
	if aPrio < bPrio {
		return -1
	}

	return 1
}

// ClusterOptions contains cluster specific options
type ClusterOptions struct {
	NodeChecker        hasql.NodeChecker
	NodePicker         hasql.NodePicker[Querier]
	NodeDiscoverer     hasql.NodeDiscoverer[Querier]
	Options            []hasql.ClusterOpt[Querier]
	Context            context.Context
	Retries            int
	NodePriority       map[string]int32
	NodeStateCriterion hasql.NodeStateCriterion
}

// ClusterOption apply cluster options to ClusterOptions
type ClusterOption func(*ClusterOptions)

// WithClusterNodeChecker pass hasql.NodeChecker to cluster options
func WithClusterNodeChecker(c hasql.NodeChecker) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodeChecker = c
	}
}

// WithClusterNodePicker pass hasql.NodePicker to cluster options
func WithClusterNodePicker(p hasql.NodePicker[Querier]) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodePicker = p
	}
}

// WithClusterNodeDiscoverer pass hasql.NodeDiscoverer to cluster options
func WithClusterNodeDiscoverer(d hasql.NodeDiscoverer[Querier]) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodeDiscoverer = d
	}
}

// WithRetries retry count on other nodes in case of error
func WithRetries(n int) ClusterOption {
	return func(o *ClusterOptions) {
		o.Retries = n
	}
}

// WithClusterContext pass context.Context to cluster options and used for checks
func WithClusterContext(ctx context.Context) ClusterOption {
	return func(o *ClusterOptions) {
		o.Context = ctx
	}
}

// WithClusterOptions pass hasql.ClusterOpt
func WithClusterOptions(opts ...hasql.ClusterOpt[Querier]) ClusterOption {
	return func(o *ClusterOptions) {
		o.Options = append(o.Options, opts...)
	}
}

// WithClusterNodeStateCriterion pass default hasql.NodeStateCriterion
func WithClusterNodeStateCriterion(c hasql.NodeStateCriterion) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodeStateCriterion = c
	}
}

type ClusterNode struct {
	Name     string
	DB       Querier
	Priority int32
}

// WithClusterNodes create cluster with static NodeDiscoverer
func WithClusterNodes(cns ...ClusterNode) ClusterOption {
	return func(o *ClusterOptions) {
		nodes := make([]*hasql.Node[Querier], 0, len(cns))
		if o.NodePriority == nil {
			o.NodePriority = make(map[string]int32, len(cns))
		}
		for _, cn := range cns {
			nodes = append(nodes, hasql.NewNode(cn.Name, cn.DB))
			if cn.Priority == 0 {
				cn.Priority = math.MaxInt32
			}
			o.NodePriority[cn.Name] = cn.Priority
		}
		o.NodeDiscoverer = hasql.NewStaticNodeDiscoverer(nodes...)
	}
}

func (c *Cluster) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	var tx *sql.Tx
	var err error

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		for ; retries < c.options.Retries; retries++ {
			if tx, err = n.DB().BeginTx(ctx, opts); err != nil && retries >= c.options.Retries {
				return true
			}
		}
		return false
	})

	if tx == nil && err == nil {
		err = errNoAliveNodes
	}

	return tx, err
}

func (c *Cluster) Close() error {
	return c.hasql.Close()
}

func (c *Cluster) Conn(ctx context.Context) (*sql.Conn, error) {
	var conn *sql.Conn
	var err error

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		for ; retries < c.options.Retries; retries++ {
			if conn, err = n.DB().Conn(ctx); err != nil && retries >= c.options.Retries {
				return true
			}
		}
		return false
	})

	if conn == nil && err == nil {
		err = errNoAliveNodes
	}

	return conn, err
}

func (c *Cluster) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var res sql.Result
	var err error

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		for ; retries < c.options.Retries; retries++ {
			if res, err = n.DB().ExecContext(ctx, query, args...); err != nil && retries >= c.options.Retries {
				return true
			}
		}
		return false
	})

	if res == nil && err == nil {
		err = errNoAliveNodes
	}

	return res, err
}

func (c *Cluster) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	var res *sql.Stmt
	var err error

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		for ; retries < c.options.Retries; retries++ {
			if res, err = n.DB().PrepareContext(ctx, query); err != nil && retries >= c.options.Retries {
				return true
			}
		}
		return false
	})

	if res == nil && err == nil {
		err = errNoAliveNodes
	}

	return res, err
}

func (c *Cluster) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var res *sql.Rows
	var err error

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		for ; retries < c.options.Retries; retries++ {
			if res, err = n.DB().QueryContext(ctx, query); err != nil && err != sql.ErrNoRows && retries >= c.options.Retries {
				return true
			}
		}
		return false
	})

	if res == nil && err == nil {
		err = errNoAliveNodes
	}

	return res, err
}

func (c *Cluster) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	var res *sql.Row

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		for ; retries < c.options.Retries; retries++ {
			res = n.DB().QueryRowContext(ctx, query, args...)
			if res.Err() == nil {
				return false
			} else if res.Err() != nil && retries >= c.options.Retries {
				return false
			}
		}
		return true
	})

	if res == nil {
		res = newSQLRowError()
	}

	return res
}

func (c *Cluster) PingContext(ctx context.Context) error {
	var err error
	var ok bool

	retries := 0
	c.hasql.NodesIter(c.getNodeStateCriterion(ctx))(func(n *hasql.Node[Querier]) bool {
		ok = true
		for ; retries < c.options.Retries; retries++ {
			if err = n.DB().PingContext(ctx); err != nil && retries >= c.options.Retries {
				return true
			}
		}
		return false
	})

	if !ok {
		err = errNoAliveNodes
	}

	return err
}

func (c *Cluster) WaitForNodes(ctx context.Context, criterions ...hasql.NodeStateCriterion) error {
	for _, criterion := range criterions {
		if _, err := c.hasql.WaitForNode(ctx, criterion); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) SetConnMaxLifetime(td time.Duration) {
	c.hasql.NodesIter(hasql.NodeStateCriterion(hasql.Alive))(func(n *hasql.Node[Querier]) bool {
		n.DB().SetConnMaxIdleTime(td)
		return false
	})
}

func (c *Cluster) SetConnMaxIdleTime(td time.Duration) {
	c.hasql.NodesIter(hasql.NodeStateCriterion(hasql.Alive))(func(n *hasql.Node[Querier]) bool {
		n.DB().SetConnMaxIdleTime(td)
		return false
	})
}

func (c *Cluster) SetMaxOpenConns(nc int) {
	c.hasql.NodesIter(hasql.NodeStateCriterion(hasql.Alive))(func(n *hasql.Node[Querier]) bool {
		n.DB().SetMaxOpenConns(nc)
		return false
	})
}

func (c *Cluster) SetMaxIdleConns(nc int) {
	c.hasql.NodesIter(hasql.NodeStateCriterion(hasql.Alive))(func(n *hasql.Node[Querier]) bool {
		n.DB().SetMaxIdleConns(nc)
		return false
	})
}

func (c *Cluster) Stats() sql.DBStats {
	s := sql.DBStats{}
	c.hasql.NodesIter(hasql.NodeStateCriterion(hasql.Alive))(func(n *hasql.Node[Querier]) bool {
		st := n.DB().Stats()
		s.Idle += st.Idle
		s.InUse += st.InUse
		s.MaxIdleClosed += st.MaxIdleClosed
		s.MaxIdleTimeClosed += st.MaxIdleTimeClosed
		s.MaxOpenConnections += st.MaxOpenConnections
		s.OpenConnections += st.OpenConnections
		s.WaitCount += st.WaitCount
		s.WaitDuration += st.WaitDuration
		return false
	})
	return s
}

func (c *Cluster) getNodeStateCriterion(ctx context.Context) hasql.NodeStateCriterion {
	if v, ok := ctx.Value(nodeStateCriterionKey{}).(hasql.NodeStateCriterion); ok {
		return v
	}
	return c.options.NodeStateCriterion
}
