package router

import (
	"fmt"
	"net"
	"strconv"
)

// NewRouter returns an initialized dns router
func NewRouter(opts ...Option) Router {
	options := NewOptions(opts...)
	return &dns{options}
}

type dns struct {
	options Options
}

func (d *dns) Init(opts ...Option) error {
	for _, o := range opts {
		o(&d.options)
	}
	return nil
}

func (d *dns) Options() Options {
	return d.options
}

func (d *dns) Table() Table {
	return nil
}

func (d *dns) Close() error {
	return nil
}

func (d *dns) Lookup(opts ...QueryOption) ([]Route, error) {
	options := NewQuery(opts...)
	// check to see if we have the port provided in the service, e.g. go-micro-srv-foo:8000
	host, port, err := net.SplitHostPort(options.Service)
	if err == nil {
		// lookup the service using A records
		ips, err := net.LookupHost(host)
		if err != nil {
			return nil, err
		}

		p, _ := strconv.Atoi(port)

		// convert the ip addresses to routes
		result := make([]Route, len(ips))
		for i, ip := range ips {
			result[i] = Route{
				Service: options.Service,
				Address: fmt.Sprintf("%s:%d", ip, p),
			}
		}
		return result, nil
	}

	// we didn't get the port so we'll lookup the service using SRV records. If we can't lookup the
	// service using the SRV record, we return the error.
	_, nodes, err := net.LookupSRV(options.Service, "tcp", d.options.Network)
	if err != nil {
		return nil, err
	}

	// convert the nodes (net services) to routes
	result := make([]Route, len(nodes))
	for i, n := range nodes {
		result[i] = Route{
			Service: options.Service,
			Address: fmt.Sprintf("%s:%d", n.Target, n.Port),
			Network: d.options.Network,
		}
	}
	return result, nil
}

func (d *dns) Watch(opts ...WatchOption) (Watcher, error) {
	return nil, nil
}

func (d *dns) Name() string {
	return d.options.Name
}

func (d *dns) String() string {
	return "dns"
}
