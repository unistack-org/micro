// Package dnssrv resolves names to dns srv records
package dnssrv // import "go.unistack.org/micro/v4/resolver/dnssrv"

import (
	"fmt"
	"net"

	"go.unistack.org/micro/v4/resolver"
)

// Resolver is a DNS network resolve
type Resolver struct {
	Address string
}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	_, addrs, err := net.LookupSRV("network", "udp", name)
	if err != nil {
		return nil, err
	}
	records := make([]*resolver.Record, 0, len(addrs))
	for _, addr := range addrs {
		address := addr.Target
		if addr.Port > 0 {
			address = fmt.Sprintf("%s:%d", addr.Target, addr.Port)
		}
		records = append(records, &resolver.Record{
			Address: address,
		})
	}
	return records, nil
}
