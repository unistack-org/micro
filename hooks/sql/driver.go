package sql

import (
	"context"
	"database/sql/driver"
	"time"
)

var (
// _ driver.DriverContext = (*wrapperDriver)(nil)
// _ driver.Connector     = (*wrapperDriver)(nil)
)

/*
type conn interface {
	driver.Pinger
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Conn
	driver.ConnPrepareContext
	driver.ConnBeginTx
}
*/

// wrapperDriver defines a wrapper for driver.Driver
type wrapperDriver struct {
	driver driver.Driver
	opts   Options
	ctx    context.Context
}

// NewWrapper creates and returns a new SQL driver with passed capabilities
func NewWrapper(d driver.Driver, opts ...Option) driver.Driver {
	return &wrapperDriver{driver: d, opts: NewOptions(opts...), ctx: context.Background()}
}

type wrappedConnector struct {
	connector driver.Connector
//	name      string
	opts      Options
	ctx       context.Context
}

func NewWrapperConnector(c driver.Connector, opts ...Option) driver.Connector {
	return &wrappedConnector{connector: c, opts: NewOptions(opts...), ctx: context.Background()}
}

// Connect implements driver.Driver Connect
func (w *wrappedConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return w.connector.Connect(ctx)
}

// Driver implements driver.Driver Driver
func (w *wrappedConnector) Driver() driver.Driver {
	return w.connector.Driver()
}

/*
// Connect implements driver.Driver OpenConnector
func (w *wrapperDriver) OpenConnector(name string) (driver.Conn, error) {
	return &wrapperConnector{driver: w.driver, name: name, opts: w.opts}, nil
}
*/

// Open implements driver.Driver Open
func (w *wrapperDriver) Open(name string) (driver.Conn, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Ensure eventual timeout
	// defer cancel()

	/*
		connector, err := w.OpenConnector(name)
		if err != nil {
			return nil, err
		}
		return connector.Connect(ctx)
	*/

	ts := time.Now()
	c, err := w.driver.Open(name)
	td := time.Since(ts)
	/*
		if w.opts.LoggerEnabled {
			w.opts.Logger.Log(w.ctx, w.opts.LoggerLevel, w.opts.LoggerObserver(w.ctx, "Open", getCallerName(), td, err)...)
		}
	*/
	_ = td
	if err != nil {
		return nil, err
	}

	return wrapConn(c, w.opts), nil
}
