// Package auth provides authentication and authorization capability
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

const (
	// BearerScheme used for Authorization header
	BearerScheme = "Bearer "
	// ScopePublic is the scope applied to a rule to allow access to the public
	ScopePublic = ""
	// ScopeAccount is the scope applied to a rule to limit to users with any valid account
	ScopeAccount = "*"
)

var (
	// DefaultAuth holds default auth implementation
	DefaultAuth Auth = NewAuth()
	// ErrInvalidToken is when the token provided is not valid
	ErrInvalidToken = errors.New("invalid token provided")
	// ErrForbidden is when a user does not have the necessary scope to access a resource
	ErrForbidden = errors.New("resource forbidden")
)

// Auth provides authentication and authorization
type Auth interface {
	// Init the auth
	Init(opts ...Option) error
	// Options set for auth
	Options() Options
	// Generate a new account
	Generate(id string, opts ...GenerateOption) (*Account, error)
	// Verify an account has access to a resource using the rules
	Verify(acc *Account, res *Resource, opts ...VerifyOption) error
	// Inspect a token
	Inspect(token string) (*Account, error)
	// Token generated using refresh token or credentials
	Token(opts ...TokenOption) (*Token, error)
	// Grant access to a resource
	Grant(rule *Rule) error
	// Revoke access to a resource
	Revoke(rule *Rule) error
	// Rules returns all the rules used to verify requests
	Rules(...RulesOption) ([]*Rule, error)
	// String returns the name of the implementation
	String() string
}

// Account provided by an auth provider
type Account struct {
	// Metadata any other associated metadata
	Metadata metadata.Metadata `json:"metadata"`
	// ID of the account e.g. email or uuid
	ID string `json:"id"`
	// Type of the account, e.g. service
	Type string `json:"type"`
	// Issuer of the account
	Issuer string `json:"issuer"`
	// Secret for the account, e.g. the password
	Secret string `json:"secret"`
	// Scopes the account has access to
	Scopes []string `json:"scopes"`
}

// Token can be short or long lived
type Token struct {
	// Time of token creation
	Created time.Time `json:"created"`
	// Time of token expiry
	Expiry time.Time `json:"expiry"`
	// The token to be used for accessing resources
	AccessToken string `json:"access_token"`
	// RefreshToken to be used to generate a new token
	RefreshToken string `json:"refresh_token"`
}

// Expired returns a boolean indicating if the token needs to be refreshed
func (t *Token) Expired() bool {
	return t.Expiry.Unix() < time.Now().Unix()
}

// Resource is an entity such as a user or
type Resource struct {
	// Name of the resource, e.g. go.micro.service.notes
	Name string `json:"name"`
	// Type of resource, e.g. service
	Type string `json:"type"`
	// Endpoint resource e.g NotesService.Create
	Endpoint string `json:"endpoint"`
}

// Access defines the type of access a rule grants
type Access int

const (
	// AccessGranted to a resource
	AccessGranted Access = iota
	// AccessDenied to a resource
	AccessDenied
)

// Rule is used to verify access to a resource
type Rule struct {
	// Resource that rule belongs to
	Resource *Resource
	// ID of the rule
	ID string
	// Scope of the rule
	Scope string
	// Access flag allow/deny
	Access Access
	// Priority holds the rule priority
	Priority int32
}

type accountKey struct{}

// AccountFromContext gets the account from the context, which
// is set by the auth wrapper at the start of a call. If the account
// is not set, a nil account will be returned. The error is only returned
// when there was a problem retrieving an account
func AccountFromContext(ctx context.Context) (*Account, bool) {
	if ctx == nil {
		return nil, false
	}
	acc, ok := ctx.Value(accountKey{}).(*Account)
	return acc, ok
}

// ContextWithAccount sets the account in the context
func ContextWithAccount(ctx context.Context, account *Account) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, accountKey{}, account)
}
