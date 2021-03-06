package runtime

import (
	"context"
	"io"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
)

// Options configure runtime
type Options struct {
	Scheduler Scheduler
	Client    client.Client
	Logger    logger.Logger
	Type      string
	Source    string
	Image     string
}

// Option func signature
type Option func(o *Options)

// WithLogger sets the logger
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithSource sets the base image / repository
func WithSource(src string) Option {
	return func(o *Options) {
		o.Source = src
	}
}

// WithScheduler specifies a scheduler for updates
func WithScheduler(n Scheduler) Option {
	return func(o *Options) {
		o.Scheduler = n
	}
}

// WithType sets the service type to manage
func WithType(t string) Option {
	return func(o *Options) {
		o.Type = t
	}
}

// WithImage sets the image to use
func WithImage(t string) Option {
	return func(o *Options) {
		o.Image = t
	}
}

// WithClient sets the client to use
func WithClient(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// CreateOption func signature
type CreateOption func(o *CreateOptions)

// ReadOption func signature
type ReadOption func(o *ReadOptions)

// CreateOptions configure runtime services
type CreateOptions struct {
	Context   context.Context
	Output    io.Writer
	Resources *Resources
	Secrets   map[string]string
	Image     string
	Namespace string
	Type      string
	Command   []string
	Args      []string
	Env       []string
	Retries   int
}

// ReadOptions queries runtime services
type ReadOptions struct {
	Context   context.Context
	Service   string
	Version   string
	Type      string
	Namespace string
}

// CreateType sets the type of service to create
func CreateType(t string) CreateOption {
	return func(o *CreateOptions) {
		o.Type = t
	}
}

// CreateImage sets the image to use
func CreateImage(img string) CreateOption {
	return func(o *CreateOptions) {
		o.Image = img
	}
}

// CreateNamespace sets the namespace
func CreateNamespace(ns string) CreateOption {
	return func(o *CreateOptions) {
		o.Namespace = ns
	}
}

// CreateContext sets the context
func CreateContext(ctx context.Context) CreateOption {
	return func(o *CreateOptions) {
		o.Context = ctx
	}
}

// WithSecret sets a secret to provide the service with
func WithSecret(key, value string) CreateOption {
	return func(o *CreateOptions) {
		if o.Secrets == nil {
			o.Secrets = map[string]string{key: value}
		} else {
			o.Secrets[key] = value
		}
	}
}

// WithCommand specifies the command to execute
func WithCommand(cmd ...string) CreateOption {
	return func(o *CreateOptions) {
		// set command
		o.Command = cmd
	}
}

// WithArgs specifies the command to execute
func WithArgs(args ...string) CreateOption {
	return func(o *CreateOptions) {
		// set command
		o.Args = args
	}
}

// WithRetries sets the max retries attempts
func WithRetries(retries int) CreateOption {
	return func(o *CreateOptions) {
		o.Retries = retries
	}
}

// WithEnv sets the created service environment
func WithEnv(env []string) CreateOption {
	return func(o *CreateOptions) {
		o.Env = env
	}
}

// WithOutput sets the arg output
func WithOutput(out io.Writer) CreateOption {
	return func(o *CreateOptions) {
		o.Output = out
	}
}

// ResourceLimits sets the resources for the service to use
func ResourceLimits(r *Resources) CreateOption {
	return func(o *CreateOptions) {
		o.Resources = r
	}
}

// ReadService returns services with the given name
func ReadService(service string) ReadOption {
	return func(o *ReadOptions) {
		o.Service = service
	}
}

// ReadVersion confifgures service version
func ReadVersion(version string) ReadOption {
	return func(o *ReadOptions) {
		o.Version = version
	}
}

// ReadType returns services of the given type
func ReadType(t string) ReadOption {
	return func(o *ReadOptions) {
		o.Type = t
	}
}

// ReadNamespace sets the namespace
func ReadNamespace(ns string) ReadOption {
	return func(o *ReadOptions) {
		o.Namespace = ns
	}
}

// ReadContext sets the context
func ReadContext(ctx context.Context) ReadOption {
	return func(o *ReadOptions) {
		o.Context = ctx
	}
}

// UpdateOption func signature
type UpdateOption func(o *UpdateOptions)

// UpdateOptions struct
type UpdateOptions struct {
	Context   context.Context
	Secrets   map[string]string
	Namespace string
}

// UpdateSecret sets a secret to provide the service with
func UpdateSecret(key, value string) UpdateOption {
	return func(o *UpdateOptions) {
		if o.Secrets == nil {
			o.Secrets = map[string]string{key: value}
		} else {
			o.Secrets[key] = value
		}
	}
}

// UpdateNamespace sets the namespace
func UpdateNamespace(ns string) UpdateOption {
	return func(o *UpdateOptions) {
		o.Namespace = ns
	}
}

// UpdateContext sets the context
func UpdateContext(ctx context.Context) UpdateOption {
	return func(o *UpdateOptions) {
		o.Context = ctx
	}
}

// DeleteOption func signature
type DeleteOption func(o *DeleteOptions)

// DeleteOptions struct
type DeleteOptions struct {
	Context   context.Context
	Namespace string
}

// DeleteNamespace sets the namespace
func DeleteNamespace(ns string) DeleteOption {
	return func(o *DeleteOptions) {
		o.Namespace = ns
	}
}

// DeleteContext sets the context
func DeleteContext(ctx context.Context) DeleteOption {
	return func(o *DeleteOptions) {
		o.Context = ctx
	}
}

// LogsOption configures runtime logging
type LogsOption func(o *LogsOptions)

// LogsOptions configure runtime logging
type LogsOptions struct {
	Context   context.Context
	Namespace string
	Count     int64
	Stream    bool
}

// LogsCount confiures how many existing lines to show
func LogsCount(count int64) LogsOption {
	return func(l *LogsOptions) {
		l.Count = count
	}
}

// LogsStream configures whether to stream new lines
func LogsStream(stream bool) LogsOption {
	return func(l *LogsOptions) {
		l.Stream = stream
	}
}

// LogsNamespace sets the namespace
func LogsNamespace(ns string) LogsOption {
	return func(o *LogsOptions) {
		o.Namespace = ns
	}
}

// LogsContext sets the context
func LogsContext(ctx context.Context) LogsOption {
	return func(o *LogsOptions) {
		o.Context = ctx
	}
}
