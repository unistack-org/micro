package micro

import (
	"fmt"
	rtime "runtime"
	"sync"

	cmd "github.com/unistack-org/micro-config-cmd"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/registry"
	"github.com/unistack-org/micro/v3/server"
	"github.com/unistack-org/micro/v3/store"
)

type service struct {
	opts Options

	once sync.Once
}

func newService(opts ...Option) Service {
	options := newOptions(opts...)

	service := &service{opts: options}

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
		if err := s.opts.Cmd.Init(
			cmd.Auth(&s.opts.Auth),
			cmd.Broker(&s.opts.Broker),
			cmd.Registry(&s.opts.Registry),
			cmd.Runtime(&s.opts.Runtime),
			cmd.Transport(&s.opts.Transport),
			cmd.Client(&s.opts.Client),
			cmd.Config(&s.opts.Config),
			cmd.Server(&s.opts.Server),
			cmd.Store(&s.opts.Store),
			cmd.Profile(&s.opts.Profile),
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

func (s *service) String() string {
	return "micro"
}

func (s *service) Start() error {
	var err error
	if logger.V(logger.InfoLevel) {
		logger.Infof("Starting [service] %s", s.Name())
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

	/*
		if s.opts.Transport != nil {
			if err := s.opts.Transport.Connect(s.opts.Context); err != nil {
				return err
			}
		}
	*/

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
	if logger.V(logger.InfoLevel) {
		logger.Infof("Stoppping [service] %s", s.Name())
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

	select {
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}
