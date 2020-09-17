package scope

import (
	"context"
	"fmt"

	"github.com/unistack-org/micro/v3/store"
)

// Scope extends the store, applying a prefix to each request
type Scope struct {
	store.Store
	prefix string
}

// NewScope returns an initialised scope
func NewScope(s store.Store, prefix string) Scope {
	return Scope{Store: s, prefix: prefix}
}

func (s *Scope) Options() store.Options {
	o := s.Store.Options()
	o.Table = s.prefix
	return o
}

func (s *Scope) Read(ctx context.Context, key string, opts ...store.ReadOption) ([]*store.Record, error) {
	key = fmt.Sprintf("%v/%v", s.prefix, key)
	return s.Store.Read(ctx, key, opts...)
}

func (s *Scope) Write(ctx context.Context, r *store.Record, opts ...store.WriteOption) error {
	r.Key = fmt.Sprintf("%v/%v", s.prefix, r.Key)
	return s.Store.Write(ctx, r, opts...)
}

func (s *Scope) Delete(ctx context.Context, key string, opts ...store.DeleteOption) error {
	key = fmt.Sprintf("%v/%v", s.prefix, key)
	return s.Store.Delete(ctx, key, opts...)
}

func (s *Scope) List(ctx context.Context, opts ...store.ListOption) ([]string, error) {
	var lops store.ListOptions
	for _, o := range opts {
		o(&lops)
	}

	key := fmt.Sprintf("%v/%v", s.prefix, lops.Prefix)
	opts = append(opts, store.ListPrefix(key))

	return s.Store.List(ctx, opts...)
}
