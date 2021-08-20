package auth

import (
	"github.com/unistack-org/micro/v3/util/id"
)

type noopAuth struct {
	opts Options
}

// String returns the name of the implementation
func (n *noopAuth) String() string {
	return "noop"
}

// Init the auth
func (n *noopAuth) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

// Options set for auth
func (n *noopAuth) Options() Options {
	return n.opts
}

// Generate a new account
func (n *noopAuth) Generate(id string, opts ...GenerateOption) (*Account, error) {
	options := NewGenerateOptions(opts...)

	return &Account{
		ID:       id,
		Secret:   options.Secret,
		Metadata: options.Metadata,
		Scopes:   options.Scopes,
		Issuer:   n.Options().Issuer,
	}, nil
}

// Grant access to a resource
func (n *noopAuth) Grant(rule *Rule) error {
	return nil
}

// Revoke access to a resource
func (n *noopAuth) Revoke(rule *Rule) error {
	return nil
}

// Rules used to verify requests
func (n *noopAuth) Rules(opts ...RulesOption) ([]*Rule, error) {
	return []*Rule{}, nil
}

// Verify an account has access to a resource
func (n *noopAuth) Verify(acc *Account, res *Resource, opts ...VerifyOption) error {
	return nil
}

// Inspect a token
func (n *noopAuth) Inspect(token string) (*Account, error) {
	id, err := id.New()
	if err != nil {
		return nil, err
	}
	return &Account{ID: id, Issuer: n.Options().Issuer}, nil
}

// Token generation using an account id and secret
func (n *noopAuth) Token(opts ...TokenOption) (*Token, error) {
	return &Token{}, nil
}

// NewAuth returns new noop auth
func NewAuth(opts ...Option) Auth {
	return &noopAuth{opts: NewOptions(opts...)}
}
