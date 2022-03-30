// Package dns resolves names to dns records
package dns // import "go.unistack.org/micro/v3/resolver/dns"

import (
	"context"
	"net"
	"sync"
	"time"

	"go.unistack.org/micro/v3/resolver"
)

// Resolver is a DNS network resolve
type Resolver struct {
	sync.RWMutex
	goresolver *net.Resolver
	Address    string
}

// Resolve tries to resolve endpoint address
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	host, port, err := net.SplitHostPort(name)
	if err != nil {
		host = name
		port = "8085"
	}

	if len(host) == 0 {
		host = "localhost"
	}

	if len(r.Address) == 0 {
		r.Address = "1.1.1.1:53"
	}

	// parsed an actual ip
	if v := net.ParseIP(host); v != nil {
		rec := &resolver.Record{Address: net.JoinHostPort(host, port)}
		return []*resolver.Record{rec}, nil
	}

	r.RLock()
	goresolver := r.goresolver
	r.RUnlock()

	if goresolver == nil {
		r.Lock()
		r.goresolver = &net.Resolver{
			Dial: func(ctx context.Context, _ string, _ string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(100),
				}
				return d.DialContext(ctx, "udp", r.Address)
			},
		}
		r.Unlock()
	}

	addrs, err := goresolver.LookupIP(context.TODO(), "ip", host)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		rec := &resolver.Record{Address: net.JoinHostPort(host, port)}
		return []*resolver.Record{rec}, nil
	}

	records := make([]*resolver.Record, 0, len(addrs))
	for _, addr := range addrs {
		records = append(records, &resolver.Record{
			Address: net.JoinHostPort(addr.String(), port),
		})
	}

	return records, nil
}
