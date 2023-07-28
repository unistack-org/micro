package store

import (
	"context"

	"go.unistack.org/micro/v4/options"
)

// NamespaceStore wrap store with namespace
type NamespaceStore struct {
	s  Store
	ns string
}

var _ Store = &NamespaceStore{}

func NewNamespaceStore(s Store, ns string) Store {
	return &NamespaceStore{s: s, ns: ns}
}

func (w *NamespaceStore) Init(opts ...options.Option) error {
	return w.s.Init(opts...)
}

func (w *NamespaceStore) Connect(ctx context.Context) error {
	return w.s.Connect(ctx)
}

func (w *NamespaceStore) Disconnect(ctx context.Context) error {
	return w.s.Disconnect(ctx)
}

func (w *NamespaceStore) Read(ctx context.Context, key string, val interface{}, opts ...options.Option) error {
	return w.s.Read(ctx, key, val, append(opts, options.Namespace(w.ns))...)
}

func (w *NamespaceStore) Write(ctx context.Context, key string, val interface{}, opts ...options.Option) error {
	return w.s.Write(ctx, key, val, append(opts, options.Namespace(w.ns))...)
}

func (w *NamespaceStore) Delete(ctx context.Context, key string, opts ...options.Option) error {
	return w.s.Delete(ctx, key, append(opts, options.Namespace(w.ns))...)
}

func (w *NamespaceStore) Exists(ctx context.Context, key string, opts ...options.Option) error {
	return w.s.Exists(ctx, key, append(opts, options.Namespace(w.ns))...)
}

func (w *NamespaceStore) List(ctx context.Context, opts ...options.Option) ([]string, error) {
	return w.s.List(ctx, append(opts, options.Namespace(w.ns))...)
}

func (w *NamespaceStore) Options() Options {
	return w.s.Options()
}

func (w *NamespaceStore) Name() string {
	return w.s.Name()
}

func (w *NamespaceStore) String() string {
	return w.s.String()
}
