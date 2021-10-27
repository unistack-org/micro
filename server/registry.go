package server

import (
	"net"
	"time"

	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/util/addr"
	"go.unistack.org/micro/v3/util/backoff"
)

var (
	// DefaultRegisterFunc uses backoff to register service
	DefaultRegisterFunc = func(svc *register.Service, config Options) error {
		var err error

		opts := []register.RegisterOption{
			register.RegisterTTL(config.RegisterTTL),
			register.RegisterDomain(config.Namespace),
		}

		for i := 0; i <= config.RegisterAttempts; i++ {
			err = config.Register.Register(config.Context, svc, opts...)
			if err == nil {
				break
			}
			// backoff then retry
			time.Sleep(backoff.Do(i + 1))
			continue
		}
		return err
	}
	// DefaultDeregisterFunc uses backoff to deregister service
	DefaultDeregisterFunc = func(svc *register.Service, config Options) error {
		var err error

		opts := []register.DeregisterOption{
			register.DeregisterDomain(config.Namespace),
		}

		for i := 0; i <= config.DeregisterAttempts; i++ {
			err = config.Register.Deregister(config.Context, svc, opts...)
			if err == nil {
				break
			}
			// backoff then retry
			time.Sleep(backoff.Do(i + 1))
			continue
		}
		return err
	}
)

// NewRegisterService returns *register.Service from Server
func NewRegisterService(s Server) (*register.Service, error) {
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

	node := &register.Node{
		ID:      opts.Name + "-" + opts.ID,
		Address: net.JoinHostPort(addr, port),
	}
	node.Metadata = metadata.Copy(opts.Metadata)

	node.Metadata["server"] = s.String()
	node.Metadata["broker"] = opts.Broker.String()
	node.Metadata["register"] = opts.Register.String()

	return &register.Service{
		Name:     opts.Name,
		Version:  opts.Version,
		Nodes:    []*register.Node{node},
		Metadata: metadata.New(0),
	}, nil
}
