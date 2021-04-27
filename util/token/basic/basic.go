package basic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/auth"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/util/token"
)

// Basic implementation of token provider, backed by the store
type Basic struct {
	store store.Store
}

// StorePrefix to isolate tokens
var StorePrefix = "tokens/"

// NewTokenProvider returns an initialized basic provider
func NewTokenProvider(opts ...token.Option) token.Provider {
	options := token.NewOptions(opts...)

	if options.Store == nil {
		options.Store = store.DefaultStore
	}

	return &Basic{
		store: options.Store,
	}
}

// Generate a token for an account
func (b *Basic) Generate(acc *auth.Account, opts ...token.GenerateOption) (*token.Token, error) {
	options := token.NewGenerateOptions(opts...)

	// marshal the account to bytes
	bytes, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	// write to the store
	key := uuid.New().String()
	err = b.store.Write(context.Background(), fmt.Sprintf("%v%v", StorePrefix, key), bytes, store.WriteTTL(options.Expiry))
	if err != nil {
		return nil, err
	}

	// return the token
	return &token.Token{
		Token:   key,
		Created: time.Now(),
		Expiry:  time.Now().Add(options.Expiry),
	}, nil
}

// Inspect a token
func (b *Basic) Inspect(t string) (*auth.Account, error) {
	// lookup the token in the store
	var val []byte
	err := b.store.Read(context.Background(), StorePrefix+t, val)
	if err == store.ErrNotFound {
		return nil, token.ErrInvalidToken
	} else if err != nil {
		return nil, err
	}
	// unmarshal the bytes
	var acc *auth.Account
	if err := json.Unmarshal(val, &acc); err != nil {
		return nil, err
	}

	return acc, nil
}

// String returns basic
func (b *Basic) String() string {
	return "basic"
}
