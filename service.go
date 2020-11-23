package micro

import (
	"fmt"
	rtime "runtime"
	"sync"

	"github.com/unistack-org/micro/v3/auth"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/registry"
	"github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/server"
	"github.com/unistack-org/micro/v3/store"
)

type service struct {
	opts Options
	sync.RWMutex
	//	once sync.Once
}

func newService(opts ...Option) Service {
	service := &service{opts: NewOptions(opts...)}
	return service
}

func (s *service) Name() string {
	return s.opts.Server.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) error {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	if s.opts.Cmd != nil {
		// set cmd name
		if len(s.opts.Cmd.App().Name) == 0 {
			s.opts.Cmd.App().Name = s.Server().Options().Name
		}

		// Initialise the command options
		if err := s.opts.Cmd.Init(); err != nil {
			return err
		}
	}

	if s.opts.Logger != nil {
		if err := s.opts.Logger.Init(
			logger.WithContext(s.opts.Context),
		); err != nil {
			return err
		}
	}

	if s.opts.Registry != nil {
		if err := s.opts.Registry.Init(
			registry.Context(s.opts.Context),
		); err != nil {
			return err
		}
	}

	if s.opts.Broker != nil {
		if err := s.opts.Broker.Init(
			broker.Context(s.opts.Context),
		); err != nil {
			return err
		}
	}

	if s.opts.Transport != nil {
		if err := s.opts.Transport.Init(
			transport.Context(s.opts.Context),
		); err != nil {
			return err
		}
	}

	if s.opts.Store != nil {
		if err := s.opts.Store.Init(
			store.Context(s.opts.Context),
		); err != nil {
			return err
		}
	}

	if s.opts.Server != nil {
		if err := s.opts.Server.Init(
			server.Context(s.opts.Context),
		); err != nil {
			return err
		}
	}

	if s.opts.Client != nil {
		if err := s.opts.Client.Init(
			client.Context(s.opts.Context),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Broker() broker.Broker {
	return s.opts.Broker
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) Store() store.Store {
	return s.opts.Store
}

func (s *service) Registry() registry.Registry {
	return s.opts.Registry
}

func (s *service) Logger() logger.Logger {
	return s.opts.Logger
}

func (s *service) Transport() transport.Transport {
	return s.opts.Transport
}

func (s *service) Config() config.Config {
	return s.opts.Config
}

func (s *service) Auth() auth.Auth {
	return s.opts.Auth
}

func (s *service) Router() router.Router {
	return s.opts.Router
}

func (s *service) String() string {
	return "micro"
}

func (s *service) Start() error {
	var err error

	s.RLock()
	config := s.opts
	s.RUnlock()

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Infof("starting [service] %s", s.Name())
	}

	for _, fn := range s.opts.BeforeStart {
		if err = fn(); err != nil {
			return err
		}
	}

	if s.opts.Server == nil {
		return fmt.Errorf("cant start nil server")
	}

	if s.opts.Registry != nil {
		if err := s.opts.Registry.Connect(s.opts.Context); err != nil {
			return err
		}
	}

	if s.opts.Broker != nil {
		if err := s.opts.Broker.Connect(s.opts.Context); err != nil {
			return err
		}
	}

	if s.opts.Store != nil {
		if err := s.opts.Store.Connect(s.opts.Context); err != nil {
			return err
		}
	}

	if err = s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err = fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	s.RLock()
	config := s.opts
	s.RUnlock()

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Infof("stoppping [service] %s", s.Name())
	}

	var err error
	for _, fn := range s.opts.BeforeStop {
		if err = fn(); err != nil {
			return err
		}
	}

	if err = s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err = fn(); err != nil {
			return err
		}
	}

	if s.opts.Registry != nil {
		if err := s.opts.Registry.Disconnect(s.opts.Context); err != nil {
			return err
		}
	}

	if s.opts.Broker != nil {
		if err := s.opts.Broker.Disconnect(s.opts.Context); err != nil {
			return err
		}
	}

	if s.opts.Store != nil {
		if err := s.opts.Store.Disconnect(s.opts.Context); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Run() error {
	// start the profiler
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

	if s.opts.Cmd != nil {
		if err := s.opts.Cmd.Run(); err != nil {
			return err
		}
	}

	if err := s.Start(); err != nil {
		return err
	}

	// wait on context cancel
	<-s.opts.Context.Done()

	return s.Stop()
}
