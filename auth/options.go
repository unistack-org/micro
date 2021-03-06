package auth

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/tracer"
)

// NewOptions creates Options struct from slice of options
func NewOptions(opts ...Option) Options {
	options := Options{
		Tracer: tracer.DefaultTracer,
		Logger: logger.DefaultLogger,
		Meter:  meter.DefaultMeter,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Options struct holds auth options
type Options struct {
	// Context holds the external options
	Context context.Context
	// Meter used for metrics
	Meter meter.Meter
	// Logger used for logging
	Logger logger.Logger
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Store used for stre data
	Store store.Store
	// Token is the services token used to authenticate itself
	Token *Token
	// LoginURL is the relative url path where a user can login
	LoginURL string
	// PrivateKey for encoding JWTs
	PrivateKey string
	// PublicKey for decoding JWTs
	PublicKey string
	// Secret is used to authenticate the service
	Secret string
	// ID is the services auth ID
	ID string
	// Issuer of the service's account
	Issuer string
	// Name holds the auth name
	Name string
	// Addrs sets the addresses of auth
	Addrs []string
}

// Option func
type Option func(o *Options)

// Addrs is the auth addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Issuer of the services account
func Issuer(i string) Option {
	return func(o *Options) {
		o.Issuer = i
	}
}

// Store to back auth
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// PublicKey is the JWT public key
func PublicKey(key string) Option {
	return func(o *Options) {
		o.PublicKey = key
	}
}

// PrivateKey is the JWT private key
func PrivateKey(key string) Option {
	return func(o *Options) {
		o.PrivateKey = key
	}
}

// Credentials sets the auth credentials
func Credentials(id, secret string) Option {
	return func(o *Options) {
		o.ID = id
		o.Secret = secret
	}
}

// ClientToken sets the auth token to use when making requests
func ClientToken(token *Token) Option {
	return func(o *Options) {
		o.Token = token
	}
}

// LoginURL sets the auth LoginURL
func LoginURL(url string) Option {
	return func(o *Options) {
		o.LoginURL = url
	}
}

// GenerateOptions struct
type GenerateOptions struct {
	Metadata metadata.Metadata
	Provider string
	Type     string
	Secret   string
	Issuer   string
	Scopes   []string
}

// GenerateOption func
type GenerateOption func(o *GenerateOptions)

// WithSecret for the generated account
func WithSecret(s string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Secret = s
	}
}

// WithType for the generated account
func WithType(t string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Type = t
	}
}

// WithMetadata for the generated account
func WithMetadata(md metadata.Metadata) GenerateOption {
	return func(o *GenerateOptions) {
		o.Metadata = metadata.Copy(md)
	}
}

// WithProvider for the generated account
func WithProvider(p string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Provider = p
	}
}

// WithScopes for the generated account
func WithScopes(s ...string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Scopes = s
	}
}

// WithIssuer for the generated account
func WithIssuer(i string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Issuer = i
	}
}

// NewGenerateOptions from a slice of options
func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	var options GenerateOptions
	for _, o := range opts {
		o(&options)
	}
	return options
}

// TokenOptions struct
type TokenOptions struct {
	ID           string
	Secret       string
	RefreshToken string
	Issuer       string
	Expiry       time.Duration
}

// TokenOption func
type TokenOption func(o *TokenOptions)

// WithExpiry for the token
func WithExpiry(ex time.Duration) TokenOption {
	return func(o *TokenOptions) {
		o.Expiry = ex
	}
}

// WithCredentials sets tye id and secret
func WithCredentials(id, secret string) TokenOption {
	return func(o *TokenOptions) {
		o.ID = id
		o.Secret = secret
	}
}

// WithToken sets the refresh token
func WithToken(rt string) TokenOption {
	return func(o *TokenOptions) {
		o.RefreshToken = rt
	}
}

// WithTokenIssuer sets the token issuer option
func WithTokenIssuer(iss string) TokenOption {
	return func(o *TokenOptions) {
		o.Issuer = iss
	}
}

// NewTokenOptions from a slice of options
func NewTokenOptions(opts ...TokenOption) TokenOptions {
	var options TokenOptions
	for _, o := range opts {
		o(&options)
	}

	// set default expiry of token
	if options.Expiry == 0 {
		options.Expiry = time.Minute
	}

	return options
}

// VerifyOptions struct
type VerifyOptions struct {
	Context   context.Context
	Namespace string
}

// VerifyOption func
type VerifyOption func(o *VerifyOptions)

// VerifyContext pass context to verify
func VerifyContext(ctx context.Context) VerifyOption {
	return func(o *VerifyOptions) {
		o.Context = ctx
	}
}

// VerifyNamespace sets thhe namespace for verify
func VerifyNamespace(ns string) VerifyOption {
	return func(o *VerifyOptions) {
		o.Namespace = ns
	}
}

// RulesOptions struct
type RulesOptions struct {
	Context   context.Context
	Namespace string
}

// RulesOption func
type RulesOption func(o *RulesOptions)

// RulesContext pass rules context
func RulesContext(ctx context.Context) RulesOption {
	return func(o *RulesOptions) {
		o.Context = ctx
	}
}

// RulesNamespace sets the rule namespace
func RulesNamespace(ns string) RulesOption {
	return func(o *RulesOptions) {
		o.Namespace = ns
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter sets the meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Tracer sets the meter
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}
