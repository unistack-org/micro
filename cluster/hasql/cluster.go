package sql

import (
	"context"
	"database/sql"
	"reflect"
	"unsafe"

	"golang.yandex/hasql/v2"
)

func newSQLRowError() *sql.Row {
	row := &sql.Row{}
	t := reflect.TypeFor[sql.Row]()
	field, _ := t.FieldByName("err")
	rowPtr := unsafe.Pointer(row)
	errFieldPtr := unsafe.Pointer(uintptr(rowPtr) + field.Offset)
	errPtr := (*error)(errFieldPtr)
	*errPtr = ErrorNoAliveNodes
	return row
}

type ClusterQuerier interface {
	Querier
	WaitForNodes(ctx context.Context, criterion ...hasql.NodeStateCriterion) error
}

type Cluster struct {
	hasql   *hasql.Cluster[Querier]
	options ClusterOptions
}

// NewCluster returns [Querier] that provides cluster of nodes
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
		err = ErrorNoAliveNodes
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
		err = ErrorNoAliveNodes
	}

	return conn, err
}

func (c *Cluster) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
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
		err = ErrorNoAliveNodes
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
		err = ErrorNoAliveNodes
	}

	return res, err
}

func (c *Cluster) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
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
		err = ErrorNoAliveNodes
	}

	return res, err
}

func (c *Cluster) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
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
		err = ErrorNoAliveNodes
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
