// Package micro is a pluggable framework for microservices
package micro

import (
	"fmt"
	"sync"

	"go.unistack.org/micro/v3/auth"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/config"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/server"
	"go.unistack.org/micro/v3/store"
	"go.unistack.org/micro/v3/tracer"
)

// Service is an interface that wraps the lower level components.
// Its works as container with building blocks for service.
type Service interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...Option) error
	// Options returns the current options
	Options() Options
	// Auth is for handling auth
	Auth(...string) auth.Auth
	// Logger is for logs
	Logger(...string) logger.Logger
	// Config if for config
	Config(...string) config.Config
	// Client is for calling services
	Client(...string) client.Client
	// Broker is for sending and receiving events
	Broker(...string) broker.Broker
	// Server is for handling requests and events
	Server(...string) server.Server
	// Store is for key/val store
	Store(...string) store.Store
	// Register
	Register(...string) register.Register
	// Tracer
	Tracer(...string) tracer.Tracer
	// Router
	Router(...string) router.Router
	// Meter
	Meter(...string) meter.Meter

	// Runtime
	// Runtime(string) (runtime.Runtime, bool)
	// Profile
	// Profile(string) (profile.Profile, bool)
	// Run the service
	Run() error
	// The service implementation
	String() string
}

// RegisterHandler is syntactic sugar for registering a handler
func RegisterHandler(s server.Server, h interface{}, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(h, opts...))
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func RegisterSubscriber(topic string, s server.Server, h interface{}, opts ...server.SubscriberOption) error {
	return s.Subscribe(s.NewSubscriber(topic, h, opts...))
}

type service struct {
	opts Options
	sync.RWMutex
	//	once sync.Once
}

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return &service{opts: NewOptions(opts...)}
}

func (s *service) Name() string {
	return s.opts.Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
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

		if err = cfg.Load(cfg.Options().Context); err != nil {
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

func (s *service) Auth(names ...string) auth.Auth {
	idx := 0
	if len(names) == 1 {
		idx = getNameIndex(names[0], s.opts.Auths)
	}
	return s.opts.Auths[idx]
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
	return "micro"
}

//nolint:gocyclo
func (s *service) Start() error {
	var err error

	s.RLock()
	config := s.opts
	s.RUnlock()

	if config.Loggers[0].V(logger.InfoLevel) {
		config.Loggers[0].Infof(s.opts.Context, "starting [service] %s", s.Name())
	}

	for _, fn := range s.opts.BeforeStart {
		if err = fn(s.opts.Context); err != nil {
			return err
		}
	}

	for _, cfg := range s.opts.Configs {
		if cfg.Options().Struct == nil {
			// skip config as the struct not passed
			continue
		}

		if err = cfg.Load(s.opts.Context); err != nil {
			return err
		}
	}

	if len(s.opts.Servers) == 0 {
		return fmt.Errorf("cant start nil server")
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
		config.Loggers[0].Infof(s.opts.Context, "stoppping [service] %s", s.Name())
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

	// wait on context cancel
	<-s.opts.Context.Done()

	return s.Stop()
}

type nameIface interface {
	Name() string
}

func getNameIndex(n string, ifaces interface{}) int {
	values, ok := ifaces.([]interface{})
	if !ok {
		return 0
	}
	for idx, iface := range values {
		if ifc, ok := iface.(nameIface); ok && ifc.Name() == n {
			return idx
		}
	}
	return 0
}
