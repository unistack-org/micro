package register

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/util/id"
)

var (
	sendEventTime = 10 * time.Millisecond
	ttlPruneTime  = time.Second
)

type node struct {
	LastSeen time.Time
	*register.Node
	TTL time.Duration
}

type record struct {
	Name    string
	Version string
	Nodes   map[string]*node
}

type memory struct {
	records  map[string]services
	watchers map[string]*watcher
	opts     register.Options
	mu       sync.RWMutex
}

// services is a KV map with service name as the key and a map of records as the value
type services map[string]map[string]*record

// NewRegister returns an initialized in-memory register
func NewRegister(opts ...register.Option) register.Register {
	r := &memory{
		opts:     register.NewOptions(opts...),
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
		m.mu.Lock()
		for namespace, services := range m.records {
			for service, versions := range services {
				for version, record := range versions {
					for id, n := range record.Nodes {
						if n.TTL != 0 && time.Since(n.LastSeen) > n.TTL {
							if m.opts.Logger.V(logger.DebugLevel) {
								m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Register TTL expired for node %s of service %s", n.ID, service))
							}
							delete(m.records[namespace][service][version].Nodes, id)
						}
					}
				}
			}
		}
		m.mu.Unlock()
	}
}

func (m *memory) sendEvent(r *register.Result) {
	m.mu.RLock()
	watchers := make([]*watcher, 0, len(m.watchers))
	for _, w := range m.watchers {
		watchers = append(watchers, w)
	}
	m.mu.RUnlock()

	for _, w := range watchers {
		select {
		case <-w.exit:
			m.mu.Lock()
			delete(m.watchers, w.id)
			m.mu.Unlock()
		default:
			select {
			case w.res <- r:
			case <-time.After(sendEventTime):
			}
		}
	}
}

func (m *memory) Connect(_ context.Context) error {
	return nil
}

func (m *memory) Disconnect(_ context.Context) error {
	return nil
}

func (m *memory) Init(opts ...register.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}

	// add services
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

func (m *memory) Options() register.Options {
	return m.opts
}

func (m *memory) Register(_ context.Context, s *register.Service, opts ...register.RegisterOption) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	options := register.NewRegisterOptions(opts...)

	// get the services for this domain from the register
	srvs, ok := m.records[options.Namespace]
	if !ok {
		srvs = make(services)
	}

	s.Namespace = options.Namespace

	// ensure the service name exists
	r := serviceToRecord(s, options.TTL)
	if _, ok := srvs[s.Name]; !ok {
		srvs[s.Name] = make(map[string]*record)
	}

	if _, ok := srvs[s.Name][s.Version]; !ok {
		srvs[s.Name][s.Version] = r
		if m.opts.Logger.V(logger.DebugLevel) {
			m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Register added new service: %s, version: %s", s.Name, s.Version))
		}
		m.records[options.Namespace] = srvs
		go m.sendEvent(&register.Result{Action: register.EventCreate, Service: s})
	}

	var addedNodes bool

	for _, n := range s.Nodes {
		// check if already exists
		if _, ok := srvs[s.Name][s.Version].Nodes[n.ID]; ok {
			continue
		}

		md := metadata.Copy(n.Metadata)

		// add the node
		srvs[s.Name][s.Version].Nodes[n.ID] = &node{
			Node: &register.Node{
				ID:       n.ID,
				Address:  n.Address,
				Metadata: md,
			},
			TTL:      options.TTL,
			LastSeen: time.Now(),
		}

		addedNodes = true
	}

	if addedNodes {
		if m.opts.Logger.V(logger.DebugLevel) {
			m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Register added new node to service: %s, version: %s", s.Name, s.Version))
		}
		go m.sendEvent(&register.Result{Action: register.EventUpdate, Service: s})
	} else {
		// refresh TTL and timestamp
		for _, n := range s.Nodes {
			if m.opts.Logger.V(logger.DebugLevel) {
				m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Updated registration for service: %s, version: %s", s.Name, s.Version))
			}
			srvs[s.Name][s.Version].Nodes[n.ID].TTL = options.TTL
			srvs[s.Name][s.Version].Nodes[n.ID].LastSeen = time.Now()
		}
	}

	m.records[options.Namespace] = srvs
	return nil
}

func (m *memory) Deregister(ctx context.Context, s *register.Service, opts ...register.DeregisterOption) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	options := register.NewDeregisterOptions(opts...)

	// if the domain doesn't exist, there is nothing to deregister
	services, ok := m.records[options.Namespace]
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
				m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Register removed node from service: %s, version: %s", s.Name, s.Version))
			}
			delete(version.Nodes, n.ID)
		}
	}

	// if the nodes not empty, we replace the version in the store and exist, the rest of the logic
	// is cleanup
	if len(version.Nodes) > 0 {
		m.records[options.Namespace][s.Name][s.Version] = version
		go m.sendEvent(&register.Result{Action: register.EventUpdate, Service: s})
		return nil
	}

	// if this version was the only version of the service, we can remove the whole service from the
	// register and exit
	if len(versions) == 1 {
		delete(m.records[options.Namespace], s.Name)
		go m.sendEvent(&register.Result{Action: register.EventDelete, Service: s})

		if m.opts.Logger.V(logger.DebugLevel) {
			m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Register removed service: %s", s.Name))
		}
		return nil
	}

	// there are other versions of the service running, so only remove this version of it
	delete(m.records[options.Namespace][s.Name], s.Version)
	go m.sendEvent(&register.Result{Action: register.EventDelete, Service: s})
	if m.opts.Logger.V(logger.DebugLevel) {
		m.opts.Logger.Debug(m.opts.Context, fmt.Sprintf("Register removed service: %s, version: %s", s.Name, s.Version))
	}

	return nil
}

func (m *memory) LookupService(ctx context.Context, name string, opts ...register.LookupOption) ([]*register.Service, error) {
	options := register.NewLookupOptions(opts...)

	// if it's a wildcard domain, return from all domains
	if options.Namespace == register.WildcardNamespace {
		m.mu.RLock()
		recs := m.records
		m.mu.RUnlock()

		var services []*register.Service

		for namespace := range recs {
			srvs, err := m.LookupService(ctx, name, append(opts, register.LookupNamespace(namespace))...)
			if err == register.ErrNotFound {
				continue
			} else if err != nil {
				return nil, err
			}
			services = append(services, srvs...)
		}

		if len(services) == 0 {
			return nil, register.ErrNotFound
		}
		return services, nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// check the domain exists
	services, ok := m.records[options.Namespace]
	if !ok {
		return nil, register.ErrNotFound
	}

	// check the service exists
	versions, ok := services[name]
	if !ok || len(versions) == 0 {
		return nil, register.ErrNotFound
	}

	// serialize the response
	result := make([]*register.Service, len(versions))

	var i int

	for _, r := range versions {
		result[i] = recordToService(r, options.Namespace)
		i++
	}

	return result, nil
}

func (m *memory) ListServices(ctx context.Context, opts ...register.ListOption) ([]*register.Service, error) {
	options := register.NewListOptions(opts...)

	// if it's a wildcard domain, list from all domains
	if options.Namespace == register.WildcardNamespace {
		m.mu.RLock()
		recs := m.records
		m.mu.RUnlock()

		var services []*register.Service

		for namespace := range recs {
			srvs, err := m.ListServices(ctx, append(opts, register.ListNamespace(namespace))...)
			if err != nil {
				return nil, err
			}
			services = append(services, srvs...)
		}

		return services, nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// ensure the domain exists
	services, ok := m.records[options.Namespace]
	if !ok {
		return make([]*register.Service, 0), nil
	}

	// serialize the result, each version counts as an individual service
	var result []*register.Service

	for _, service := range services {
		for _, version := range service {
			result = append(result, recordToService(version, options.Namespace))
		}
	}

	return result, nil
}

func (m *memory) Watch(ctx context.Context, opts ...register.WatchOption) (register.Watcher, error) {
	id, err := id.New()
	if err != nil {
		return nil, err
	}
	wo := register.NewWatchOptions(opts...)
	// construct the watcher
	w := &watcher{
		exit: make(chan bool),
		res:  make(chan *register.Result),
		id:   id,
		wo:   wo,
	}

	m.mu.Lock()
	m.watchers[w.id] = w
	m.mu.Unlock()

	return w, nil
}

func (m *memory) Name() string {
	return m.opts.Name
}

func (m *memory) String() string {
	return "memory"
}

type watcher struct {
	res  chan *register.Result
	exit chan bool
	wo   register.WatchOptions
	id   string
}

func (m *watcher) Next() (*register.Result, error) {
	for {
		select {
		case r := <-m.res:
			if r.Service == nil {
				continue
			}

			if len(m.wo.Service) > 0 && m.wo.Service != r.Service.Name {
				continue
			}

			namespace := register.DefaultNamespace
			if r.Service.Namespace != "" {
				namespace = r.Service.Namespace
			}

			// only send the event if watching the wildcard or this specific domain
			if m.wo.Namespace == register.WildcardNamespace || m.wo.Namespace == namespace {
				return r, nil
			}
		case <-m.exit:
			return nil, register.ErrWatcherStopped
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

func serviceToRecord(s *register.Service, ttl time.Duration) *record {
	nodes := make(map[string]*node, len(s.Nodes))
	for _, n := range s.Nodes {
		nodes[n.ID] = &node{
			Node:     n,
			TTL:      ttl,
			LastSeen: time.Now(),
		}
	}

	return &record{
		Name:    s.Name,
		Version: s.Version,
		Nodes:   nodes,
	}
}

func recordToService(r *record, namespace string) *register.Service {
	nodes := make([]*register.Node, len(r.Nodes))
	i := 0
	for _, n := range r.Nodes {
		nmd := metadata.Copy(n.Metadata)

		nodes[i] = &register.Node{
			ID:       n.ID,
			Address:  n.Address,
			Metadata: nmd,
		}
		i++
	}

	return &register.Service{
		Name:      r.Name,
		Version:   r.Version,
		Nodes:     nodes,
		Namespace: namespace,
	}
}
