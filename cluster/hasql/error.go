package sql

import "errors"

var (
	ErrClusterChecker    = errors.New("cluster node checker required")
	ErrClusterDiscoverer = errors.New("cluster node discoverer required")
	ErrClusterPicker     = errors.New("cluster node picker required")
	ErrorNoAliveNodes    = errors.New("cluster no alive nodes")
)
