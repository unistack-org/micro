// Package micro is a pluggable framework for microservices
package micro

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"go.uber.org/automaxprocs/maxprocs"
	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/config"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/router"
	"go.unistack.org/micro/v4/server"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/tracer"
	utildns "go.unistack.org/micro/v4/util/dns"
)

func init() {
	_, _ = maxprocs.Set()
	_, _ = memlimit.SetGoMemLimitWithOpts(
		memlimit.WithRatio(0.9),
		memlimit.WithProvider(
			memlimit.ApplyFallback(
				memlimit.FromCgroup,
				memlimit.FromSystem,
			),
		),
	)

	net.DefaultResolver = utildns.NewNetResolver(
		utildns.Timeout(1*time.Second),
		utildns.MinCacheTTL(5*time.Second),
	)
}

// Service is an interface that wraps the lower level components.
// Its works as container with building blocks for service.
type Service interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...Option) error
	// Options returns the current options
	Options() Options
	// Logger is for output log from service components
	Logger(...string) logger.Logger
	// Config if for config handling via load/save methods and also with ability to watch changes
	Config(...string) config.Config
	// Client is for sync calling services via RPC
	Client(...string) client.Client
	// Broker is for sending and receiving async events
	Broker(...string) broker.Broker
	// Server is for handling requests and broker unmarshaled events
	Server(...string) server.Server
	// Store is for key/val store
	Store(...string) store.Store
	// Register used by client to lookup other services and server registers on it
	Register(...string) register.Register
	// Tracer
	Tracer(...string) tracer.Tracer
	// Router
	Router(...string) router.Router
	// Meter may be used internally by other component to export metrics
	Meter(...string) meter.Meter

	// Runtime
	// Runtime(string) (runtime.Runtime, bool)
	// Profile
	// Profile(string) (profile.Profile, bool)
	// Run the service and wait for stop
	Run() error
	// Start the service and not wait
	Start() error
	// Stop the service
	Stop() error
	// String service representation
	String() string
	// Live returns service liveness
	Live() bool
	// Ready returns service readiness
	Ready() bool
	// Health returns service health
	Health() bool
}

// RegisterHandler is syntactic sugar for registering a handler
func RegisterHandler(s server.Server, h interface{}, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(h, opts...))
}

type service struct {
	done chan struct{}
	opts Options
	sync.RWMutex
	stopped bool
}

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return &service{opts: NewOptions(opts...), done: make(chan struct{})}
}

func (s *service) Name() string {
	return s.opts.Name
}

// Init initialises options.
//
//nolint:gocyclo
func (s *service) Init(opts ...Option) error {
	var err error
	// process options
	for _, o := range opts {
		if err = o(&s.opts); err != nil {
			return err
		}
	}

	for _, cfg := range s.opts.Configs {
		if cfg.Options().Struct == nil {
			// skip config as the struct not passed
			continue
		}
		if err = cfg.Init(config.Context(cfg.Options().Context)); err != nil {
			return err
		}
	}

	for _, log := range s.opts.Loggers {
		if err = log.Init(logger.WithContext(log.Options().Context)); err != nil {
			return err
		}
	}

	for _, reg := range s.opts.Registers {
		if err = reg.Init(register.Context(reg.Options().Context)); err != nil {
			return err
		}
	}

	for _, brk := range s.opts.Brokers {
		if err = brk.Init(broker.Context(brk.Options().Context)); err != nil {
			return err
		}
	}

	for _, str := range s.opts.Stores {
		if err = str.Init(store.Context(str.Options().Context)); err != nil {
			return err
		}
	}

	for _, srv := range s.opts.Servers {
		if err = srv.Init(server.Context(srv.Options().Context)); err != nil {
			return err
		}
	}

	for _, cli := range s.opts.Clients {
		if err = cli.Init(client.Context(cli.Options().Context)); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Broker(names ...string) broker.Broker {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Brokers)
	}
	return s.opts.Brokers[idx]
}

func (s *service) Tracer(names ...string) tracer.Tracer {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Tracers)
	}
	return s.opts.Tracers[idx]
}

func (s *service) Config(names ...string) config.Config {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Configs)
	}
	return s.opts.Configs[idx]
}

func (s *service) Client(names ...string) client.Client {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Clients)
	}
	return s.opts.Clients[idx]
}

func (s *service) Server(names ...string) server.Server {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Servers)
	}
	return s.opts.Servers[idx]
}

func (s *service) Store(names ...string) store.Store {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Stores)
	}
	return s.opts.Stores[idx]
}

func (s *service) Register(names ...string) register.Register {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Registers)
	}
	return s.opts.Registers[idx]
}

func (s *service) Logger(names ...string) logger.Logger {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Loggers)
	}
	return s.opts.Loggers[idx]
}

func (s *service) Router(names ...string) router.Router {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Routers)
	}
	return s.opts.Routers[idx]
}

func (s *service) Meter(names ...string) meter.Meter {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Meters)
	}
	return s.opts.Meters[idx]
}

func (s *service) String() string {
	return s.opts.Name
}

func (s *service) Live() bool {
	for _, v := range s.opts.Brokers {
		if !v.Live() {
			return false
		}
	}
	for _, v := range s.opts.Servers {
		if !v.Live() {
			return false
		}
	}
	for _, v := range s.opts.Stores {
		if !v.Live() {
			return false
		}
	}
	return true
}

func (s *service) Ready() bool {
	for _, v := range s.opts.Brokers {
		if !v.Ready() {
			return false
		}
	}
	for _, v := range s.opts.Servers {
		if !v.Ready() {
			return false
		}
	}
	for _, v := range s.opts.Stores {
		if !v.Ready() {
			return false
		}
	}
	return true
}

func (s *service) Health() bool {
	for _, v := range s.opts.Brokers {
		if !v.Health() {
			return false
		}
	}
	for _, v := range s.opts.Servers {
		if !v.Health() {
			return false
		}
	}
	for _, v := range s.opts.Stores {
		if !v.Health() {
			return false
		}
	}
	return true
}

//nolint:gocyclo
func (s *service) Start() error {
	var err error

	s.RLock()
	config := s.opts
	s.RUnlock()

	for _, cfg := range s.opts.Configs {
		if cfg.Options().Struct == nil {
			// skip config as the struct not passed
			continue
		}

		if err = cfg.Load(s.opts.Context); err != nil {
			return err
		}
	}

	for _, fn := range s.opts.BeforeStart {
		if err = fn(s.opts.Context); err != nil {
			return err
		}
	}

	if config.Loggers[0].V(logger.InfoLevel) {
		config.Loggers[0].Info(s.opts.Context, fmt.Sprintf("starting [service] %s version %s", s.Options().Name, s.Options().Version))
	}

	for _, reg := range s.opts.Registers {
		if err = reg.Connect(s.opts.Context); err != nil {
			return err
		}
	}

	for _, brk := range s.opts.Brokers {
		if err = brk.Connect(s.opts.Context); err != nil {
			return err
		}
	}

	for _, str := range s.opts.Stores {
		if err = str.Connect(s.opts.Context); err != nil {
			return err
		}
	}

	for _, srv := range s.opts.Servers {
		if err = srv.Start(); err != nil {
			return err
		}
	}

	for _, fn := range s.opts.AfterStart {
		if err = fn(s.opts.Context); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	s.RLock()
	config := s.opts
	s.RUnlock()

	if config.Loggers[0].V(logger.InfoLevel) {
		config.Loggers[0].Info(s.opts.Context, fmt.Sprintf("stoppping [service] %s", s.Name()))
	}

	var err error
	for _, fn := range s.opts.BeforeStop {
		if err = fn(s.opts.Context); err != nil {
			return err
		}
	}

	for _, srv := range s.opts.Servers {
		if err = srv.Stop(); err != nil {
			return err
		}
	}

	for _, fn := range s.opts.AfterStop {
		if err = fn(s.opts.Context); err != nil {
			return err
		}
	}

	for _, reg := range s.opts.Registers {
		if err = reg.Disconnect(s.opts.Context); err != nil {
			return err
		}
	}

	for _, brk := range s.opts.Brokers {
		if err = brk.Disconnect(s.opts.Context); err != nil {
			return err
		}
	}

	for _, str := range s.opts.Stores {
		if err = str.Disconnect(s.opts.Context); err != nil {
			return err
		}
	}

	s.notifyShutdown()

	return nil
}

func (s *service) Run() error {
	// start the profiler
	/*
		if s.opts.Profile != nil {
			// to view mutex contention
			rtime.SetMutexProfileFraction(5)
			// to view blocking profile
			rtime.SetBlockProfileRate(1)

			if err := s.opts.Profile.Start(); err != nil {
				return err
			}
			defer s.opts.Profile.Stop()
		}
	*/
	if err := s.Start(); err != nil {
		return err
	}

	<-s.done

	return nil
}

// notifyShutdown marks the service as stopped and closes the done channel.
// It ensures the channel is closed only once, preventing multiple closures.
func (s *service) notifyShutdown() {
	s.Lock()
	if s.stopped {
		s.Unlock()
		return
	}
	s.stopped = true
	s.Unlock()

	close(s.done)
}

type Namer interface {
	Name() string
}

func getNameIndex[T Namer](n string, ifaces []T) int {
	for idx, iface := range ifaces {
		if iface.Name() == n {
			return idx
		}
	}
	return 0
}
