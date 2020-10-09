package registry

import (
	"net"

	"github.com/unistack-org/micro/v3/registry"
	"github.com/unistack-org/micro/v3/server"
	"github.com/unistack-org/micro/v3/util/addr"
)

func addNodes(old, neu []*registry.Node) []*registry.Node {
	nodes := make([]*registry.Node, len(neu))
	// add all new nodes
	for i, n := range neu {
		node := *n
		nodes[i] = &node
	}

	// look at old nodes
	for _, o := range old {
		var exists bool

		// check against new nodes
		for _, n := range nodes {
			// ids match then skip
			if o.Id == n.Id {
				exists = true
				break
			}
		}

		// keep old node
		if !exists {
			node := *o
			nodes = append(nodes, &node)
		}
	}

	return nodes
}

func delNodes(old, del []*registry.Node) []*registry.Node {
	var nodes []*registry.Node
	for _, o := range old {
		var rem bool
		for _, n := range del {
			if o.Id == n.Id {
				rem = true
				break
			}
		}
		if !rem {
			nodes = append(nodes, o)
		}
	}
	return nodes
}

// CopyService make a copy of service
func CopyService(service *registry.Service) *registry.Service {
	// copy service
	s := new(registry.Service)
	*s = *service

	// copy nodes
	nodes := make([]*registry.Node, len(service.Nodes))
	for j, node := range service.Nodes {
		n := new(registry.Node)
		*n = *node
		nodes[j] = n
	}
	s.Nodes = nodes

	// copy endpoints
	eps := make([]*registry.Endpoint, len(service.Endpoints))
	for j, ep := range service.Endpoints {
		e := new(registry.Endpoint)
		*e = *ep
		eps[j] = e
	}
	s.Endpoints = eps
	return s
}

// Copy makes a copy of services
func Copy(current []*registry.Service) []*registry.Service {
	services := make([]*registry.Service, len(current))
	for i, service := range current {
		services[i] = CopyService(service)
	}
	return services
}

// Merge merges two lists of services and returns a new copy
func Merge(olist []*registry.Service, nlist []*registry.Service) []*registry.Service {
	var srv []*registry.Service

	for _, n := range nlist {
		var seen bool
		for _, o := range olist {
			if o.Version == n.Version {
				sp := new(registry.Service)
				// make copy
				*sp = *o
				// set nodes
				sp.Nodes = addNodes(o.Nodes, n.Nodes)

				// mark as seen
				seen = true
				srv = append(srv, sp)
				break
			} else {
				sp := new(registry.Service)
				// make copy
				*sp = *o
				srv = append(srv, sp)
			}
		}
		if !seen {
			srv = append(srv, Copy([]*registry.Service{n})...)
		}
	}
	return srv
}

// Remove removes services and returns a new copy
func Remove(old, del []*registry.Service) []*registry.Service {
	var services []*registry.Service

	for _, o := range old {
		srv := new(registry.Service)
		*srv = *o

		var rem bool

		for _, s := range del {
			if srv.Version == s.Version {
				srv.Nodes = delNodes(srv.Nodes, s.Nodes)

				if len(srv.Nodes) == 0 {
					rem = true
				}
			}
		}

		if !rem {
			services = append(services, srv)
		}
	}

	return services
}

func NewService(s server.Server) (*registry.Service, error) {
	opts := s.Options()

	advt := opts.Address
	if len(opts.Advertise) > 0 {
		advt = opts.Advertise
	}

	host, port, err := net.SplitHostPort(advt)
	if err != nil {
		return nil, err
	}

	addr, err := addr.Extract(host)
	if err != nil {
		addr = host
	}

	node := &registry.Node{
		Id:       opts.Name + "-" + opts.Id,
		Address:  net.JoinHostPort(addr, port),
		Metadata: opts.Metadata,
	}

	if node.Metadata == nil {
		node.Metadata = make(map[string]string, 3)
	}

	node.Metadata["server"] = s.String()
	node.Metadata["broker"] = opts.Broker.String()
	node.Metadata["registry"] = opts.Registry.String()

	return &registry.Service{
		Name:    opts.Name,
		Version: opts.Version,
		Nodes:   []*registry.Node{node},
	}, nil
}
