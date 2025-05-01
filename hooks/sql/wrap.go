package sql

import (
	"database/sql/driver"
)

/*
func wrapDriver(d driver.Driver, opts Options) driver.Driver {
	if _, ok := d.(driver.DriverContext); ok {
		return &wrapperDriver{driver: d, opts: opts}
	}
	return struct{ driver.Driver }{&wrapperDriver{driver: d, opts: opts}}
}
*/

// WrapConn allows an existing driver.Conn to be wrapped.
func WrapConn(c driver.Conn, opts ...Option) driver.Conn {
	return wrapConn(c, NewOptions(opts...))
}
