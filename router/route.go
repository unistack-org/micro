package router

import (
	"hash/fnv"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	// DefaultLink is default network link
	DefaultLink = "local"
	// DefaultLocalMetric is default route cost for a local route
	DefaultLocalMetric int64 = 1
)

// Route is network route
type Route struct {
	// Metadata for the route
	Metadata metadata.Metadata
	// Service is destination service name
	Service string
	// Gateway is route gateway
	Gateway string
	// Network is network address
	Network string
	// Router is router id
	Router string
	// Link is network link
	Link string
	// Address is service node address
	Address string
	// Metric is the route cost metric
	Metric int64
}

// Hash returns route hash sum.
func (r *Route) Hash() uint64 {
	h := fnv.New64()
	//nolint:errcheck
	h.Write([]byte(r.Service))
	//nolint:errcheck
	h.Write([]byte(r.Address))
	//nolint:errcheck
	h.Write([]byte(r.Gateway))
	//nolint:errcheck
	h.Write([]byte(r.Network))
	//nolint:errcheck
	h.Write([]byte(r.Router))
	//nolint:errcheck
	h.Write([]byte(r.Link))
	return h.Sum64()
}
