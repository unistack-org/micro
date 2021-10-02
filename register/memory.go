package register

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/util/id"
)

var (
	sendEventTime = 10 * time.Millisecond
	ttlPruneTime  = time.Second
)

type node struct {
	LastSeen time.Time
	*Node
	TTL time.Duration
}

type record struct {
	Name      string
	Version   string
	Metadata  map[string]string
	Nodes     map[string]*node
	Endpoints []*Endpoint
}

type memory struct {
	records  map[string]services
	watchers map[string]*watcher
	opts     Options
	sync.RWMutex
}

// services is a KV map with service name as the key and a map of records as the value
type services map[string]map[string]*record

// NewRegister returns an initialized in-memory register
func NewRegister(opts ...Option) Register {
	r := &memory{
		opts:     NewOptions(opts...),
		records:  make(map[string]services),
		watchers: make(map[string]*watcher),
	}

	go r.ttlPrune()

	return r
}

func (m *memory) ttlPrune() {
	prune := time.NewTicker(ttlPruneTime)
	defer prune.Stop()

	for range prune.C {
		m.Lock()
		for domain, services := range m.records {
			for service, versions := range services {
				for version, record := range versions {
					for id, n := range record.Nodes {
						if n.TTL != 0 && time.Since(n.LastSeen) > n.TTL {
							if m.opts.Logger.V(logger.DebugLevel) {
								m.opts.Logger.Debugf(m.opts.Context, "Register TTL expired for node %s of service %s", n.ID, service)
							}
							delete(m.records[domain][service][version].Nodes, id)
						}
					}
				}
			}
		}
		m.Unlock()
	}
}

func (m *memory) sendEvent(r *Result) {
	m.RLock()
	watchers := make([]*watcher, 0, len(m.watchers))
	for _, w := range m.watchers {
		watchers = append(watchers, w)
	}
	m.RUnlock()

	for _, w := range watchers {
		select {
		case <-w.exit:
			m.Lock()
			delete(m.watchers, w.id)
			m.Unlock()
		default:
			select {
			case w.res <- r:
			case <-time.After(sendEventTime):
			}
		}
	}
}

func (m *memory) Connect(ctx context.Context) error {
	return nil
}

func (m *memory) Disconnect(ctx context.Context) error {
	return nil
}

func (m *memory) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.opts)
	}

	// add services
	m.Lock()
	defer m.Unlock()

	return nil
}

func (m *memory) Options() Options {
	return m.opts
}

func (m *memory) Register(ctx context.Context, s *Service, opts ...RegisterOption) error {
	m.Lock()
	defer m.Unlock()

	options := NewRegisterOptions(opts...)

	// get the services for this domain from the register
	srvs, ok := m.records[options.Domain]
	if !ok {
		srvs = make(services)
	}

	// domain is set in metadata so it can be passed to watchers
	if s.Metadata == nil {
		s.Metadata = map[string]string{"domain": options.Domain}
	} else {
		s.Metadata["domain"] = options.Domain
	}

	// ensure the service name exists
	r := serviceToRecord(s, options.TTL)
	if _, ok := srvs[s.Name]; !ok {
		srvs[s.Name] = make(map[string]*record)
	}

	if _, ok := srvs[s.Name][s.Version]; !ok {
		srvs[s.Name][s.Version] = r
		if m.opts.Logger.V(logger.DebugLevel) {
			m.opts.Logger.Debugf(m.opts.Context, "Register added new service: %s, version: %s", s.Name, s.Version)
		}
		m.records[options.Domain] = srvs
		go m.sendEvent(&Result{Action: "create", Service: s})
	}

	var addedNodes bool

	for _, n := range s.Nodes {
		// check if already exists
		if _, ok := srvs[s.Name][s.Version].Nodes[n.ID]; ok {
			continue
		}

		metadata := make(map[string]string)

		// make copy of metadata
		for k, v := range n.Metadata {
			metadata[k] = v
		}

		// set the domain
		metadata["domain"] = options.Domain

		// add the node
		srvs[s.Name][s.Version].Nodes[n.ID] = &node{
			Node: &Node{
				ID:       n.ID,
				Address:  n.Address,
				Metadata: metadata,
			},
			TTL:      options.TTL,
			LastSeen: time.Now(),
		}

		addedNodes = true
	}

	if addedNodes {
		if m.opts.Logger.V(logger.DebugLevel) {
			m.opts.Logger.Debugf(m.opts.Context, "Register added new node to service: %s, version: %s", s.Name, s.Version)
		}
		go m.sendEvent(&Result{Action: "update", Service: s})
	} else {
		// refresh TTL and timestamp
		for _, n := range s.Nodes {
			if m.opts.Logger.V(logger.DebugLevel) {
				m.opts.Logger.Debugf(m.opts.Context, "Updated registration for service: %s, version: %s", s.Name, s.Version)
			}
			srvs[s.Name][s.Version].Nodes[n.ID].TTL = options.TTL
			srvs[s.Name][s.Version].Nodes[n.ID].LastSeen = time.Now()
		}
	}

	m.records[options.Domain] = srvs
	return nil
}

func (m *memory) Deregister(ctx context.Context, s *Service, opts ...DeregisterOption) error {
	m.Lock()
	defer m.Unlock()

	options := NewDeregisterOptions(opts...)

	// domain is set in metadata so it can be passed to watchers
	if s.Metadata == nil {
		s.Metadata = map[string]string{"domain": options.Domain}
	} else {
		s.Metadata["domain"] = options.Domain
	}

	// if the domain doesn't exist, there is nothing to deregister
	services, ok := m.records[options.Domain]
	if !ok {
		return nil
	}

	// if no services with this name and version exist, there is nothing to deregister
	versions, ok := services[s.Name]
	if !ok {
		return nil
	}

	version, ok := versions[s.Version]
	if !ok {
		return nil
	}

	// deregister all of the service nodes from this version
	for _, n := range s.Nodes {
		if _, ok := version.Nodes[n.ID]; ok {
			if m.opts.Logger.V(logger.DebugLevel) {
				m.opts.Logger.Debugf(m.opts.Context, "Register removed node from service: %s, version: %s", s.Name, s.Version)
			}
			delete(version.Nodes, n.ID)
		}
	}

	// if the nodes not empty, we replace the version in the store and exist, the rest of the logic
	// is cleanup
	if len(version.Nodes) > 0 {
		m.records[options.Domain][s.Name][s.Version] = version
		go m.sendEvent(&Result{Action: "update", Service: s})
		return nil
	}

	// if this version was the only version of the service, we can remove the whole service from the
	// register and exit
	if len(versions) == 1 {
		delete(m.records[options.Domain], s.Name)
		go m.sendEvent(&Result{Action: "delete", Service: s})

		if m.opts.Logger.V(logger.DebugLevel) {
			m.opts.Logger.Debugf(m.opts.Context, "Register removed service: %s", s.Name)
		}
		return nil
	}

	// there are other versions of the service running, so only remove this version of it
	delete(m.records[options.Domain][s.Name], s.Version)
	go m.sendEvent(&Result{Action: "delete", Service: s})
	if m.opts.Logger.V(logger.DebugLevel) {
		m.opts.Logger.Debugf(m.opts.Context, "Register removed service: %s, version: %s", s.Name, s.Version)
	}

	return nil
}

func (m *memory) LookupService(ctx context.Context, name string, opts ...LookupOption) ([]*Service, error) {
	options := NewLookupOptions(opts...)

	// if it's a wildcard domain, return from all domains
	if options.Domain == WildcardDomain {
		m.RLock()
		recs := m.records
		m.RUnlock()

		var services []*Service

		for domain := range recs {
			srvs, err := m.LookupService(ctx, name, append(opts, LookupDomain(domain))...)
			if err == ErrNotFound {
				continue
			} else if err != nil {
				return nil, err
			}
			services = append(services, srvs...)
		}

		if len(services) == 0 {
			return nil, ErrNotFound
		}
		return services, nil
	}

	m.RLock()
	defer m.RUnlock()

	// check the domain exists
	services, ok := m.records[options.Domain]
	if !ok {
		return nil, ErrNotFound
	}

	// check the service exists
	versions, ok := services[name]
	if !ok || len(versions) == 0 {
		return nil, ErrNotFound
	}

	// serialize the response
	result := make([]*Service, len(versions))

	var i int

	for _, r := range versions {
		result[i] = recordToService(r, options.Domain)
		i++
	}

	return result, nil
}

func (m *memory) ListServices(ctx context.Context, opts ...ListOption) ([]*Service, error) {
	options := NewListOptions(opts...)

	// if it's a wildcard domain, list from all domains
	if options.Domain == WildcardDomain {
		m.RLock()
		recs := m.records
		m.RUnlock()

		var services []*Service

		for domain := range recs {
			srvs, err := m.ListServices(ctx, append(opts, ListDomain(domain))...)
			if err != nil {
				return nil, err
			}
			services = append(services, srvs...)
		}

		return services, nil
	}

	m.RLock()
	defer m.RUnlock()

	// ensure the domain exists
	services, ok := m.records[options.Domain]
	if !ok {
		return make([]*Service, 0), nil
	}

	// serialize the result, each version counts as an individual service
	var result []*Service

	for _, service := range services {
		for _, version := range service {
			result = append(result, recordToService(version, options.Domain))
		}
	}

	return result, nil
}

func (m *memory) Watch(ctx context.Context, opts ...WatchOption) (Watcher, error) {
	id, err := id.New()
	if err != nil {
		return nil, err
	}
	wo := NewWatchOptions(opts...)
	// construct the watcher
	w := &watcher{
		exit: make(chan bool),
		res:  make(chan *Result),
		id:   id,
		wo:   wo,
	}

	m.Lock()
	m.watchers[w.id] = w
	m.Unlock()

	return w, nil
}

func (m *memory) Name() string {
	return m.opts.Name
}

func (m *memory) String() string {
	return "memory"
}

type watcher struct {
	res  chan *Result
	exit chan bool
	wo   WatchOptions
	id   string
}

func (m *watcher) Next() (*Result, error) {
	for {
		select {
		case r := <-m.res:
			if r.Service == nil {
				continue
			}

			if len(m.wo.Service) > 0 && m.wo.Service != r.Service.Name {
				continue
			}

			// extract domain from service metadata
			var domain string
			if r.Service.Metadata != nil && len(r.Service.Metadata["domain"]) > 0 {
				domain = r.Service.Metadata["domain"]
			} else {
				domain = DefaultDomain
			}

			// only send the event if watching the wildcard or this specific domain
			if m.wo.Domain == WildcardDomain || m.wo.Domain == domain {
				return r, nil
			}
		case <-m.exit:
			return nil, errors.New("watcher stopped")
		}
	}
}

func (m *watcher) Stop() {
	select {
	case <-m.exit:
		return
	default:
		close(m.exit)
	}
}

func serviceToRecord(s *Service, ttl time.Duration) *record {
	metadata := make(map[string]string, len(s.Metadata))
	for k, v := range s.Metadata {
		metadata[k] = v
	}

	nodes := make(map[string]*node, len(s.Nodes))
	for _, n := range s.Nodes {
		nodes[n.ID] = &node{
			Node:     n,
			TTL:      ttl,
			LastSeen: time.Now(),
		}
	}

	endpoints := make([]*Endpoint, len(s.Endpoints))
	for i, e := range s.Endpoints {
		endpoints[i] = e
	}

	return &record{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  metadata,
		Nodes:     nodes,
		Endpoints: endpoints,
	}
}

func recordToService(r *record, domain string) *Service {
	metadata := make(map[string]string, len(r.Metadata))
	for k, v := range r.Metadata {
		metadata[k] = v
	}

	// set the domain in metadata so it can be determined when a wildcard query is performed
	metadata["domain"] = domain

	endpoints := make([]*Endpoint, len(r.Endpoints))
	for i, e := range r.Endpoints {
		md := make(map[string]string, len(e.Metadata))
		for k, v := range e.Metadata {
			md[k] = v
		}

		endpoints[i] = &Endpoint{
			Name:     e.Name,
			Request:  e.Request,
			Response: e.Response,
			Metadata: md,
		}
	}

	nodes := make([]*Node, len(r.Nodes))
	i := 0
	for _, n := range r.Nodes {
		md := make(map[string]string, len(n.Metadata))
		for k, v := range n.Metadata {
			md[k] = v
		}

		nodes[i] = &Node{
			ID:       n.ID,
			Address:  n.Address,
			Metadata: md,
		}
		i++
	}

	return &Service{
		Name:      r.Name,
		Version:   r.Version,
		Metadata:  metadata,
		Endpoints: endpoints,
		Nodes:     nodes,
	}
}
