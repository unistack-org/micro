package auth

import (
	"github.com/google/uuid"
)

type NoopAuth struct {
	opts Options
}

// String returns the name of the implementation
func (n *NoopAuth) String() string {
	return "noop"
}

// Init the auth
func (n *NoopAuth) Init(opts ...Option) {
	for _, o := range opts {
		o(&n.opts)
	}
}

// Options set for auth
func (n *NoopAuth) Options() Options {
	return n.opts
}

// Generate a new account
func (n *NoopAuth) Generate(id string, opts ...GenerateOption) (*Account, error) {
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
func (n *NoopAuth) Grant(rule *Rule) error {
	return nil
}

// Revoke access to a resource
func (n *NoopAuth) Revoke(rule *Rule) error {
	return nil
}

// Rules used to verify requests
func (n *NoopAuth) Rules(opts ...RulesOption) ([]*Rule, error) {
	return []*Rule{}, nil
}

// Verify an account has access to a resource
func (n *NoopAuth) Verify(acc *Account, res *Resource, opts ...VerifyOption) error {
	return nil
}

// Inspect a token
func (n *NoopAuth) Inspect(token string) (*Account, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Account{ID: uid.String(), Issuer: n.Options().Issuer}, nil
}

// Token generation using an account id and secret
func (n *NoopAuth) Token(opts ...TokenOption) (*Token, error) {
	return &Token{}, nil
}
