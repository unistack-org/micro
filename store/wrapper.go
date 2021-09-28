package store

import (
	"context"
)

// LogfFunc function used for Logf method
// type LogfFunc func(ctx context.Context, level Level, msg string, args ...interface{})
// type Wrapper interface {
//   Logf logs message with needed level
//   Logf(LogfFunc) LogfFunc
// }

// NamespaceStore wrap store with namespace
type NamespaceStore struct {
	s  Store
	ns string
}

var _ Store = &NamespaceStore{}

func NewNamespaceStore(s Store, ns string) Store {
	return &NamespaceStore{s: s, ns: ns}
}

func (w *NamespaceStore) Init(opts ...Option) error {
	return w.s.Init(opts...)
}

func (w *NamespaceStore) Connect(ctx context.Context) error {
	return w.s.Connect(ctx)
}

func (w *NamespaceStore) Disconnect(ctx context.Context) error {
	return w.s.Disconnect(ctx)
}

func (w *NamespaceStore) Read(ctx context.Context, key string, val interface{}, opts ...ReadOption) error {
	return w.s.Read(ctx, key, val, append(opts, ReadNamespace(w.ns))...)
}

func (w *NamespaceStore) Write(ctx context.Context, key string, val interface{}, opts ...WriteOption) error {
	return w.s.Write(ctx, key, val, append(opts, WriteNamespace(w.ns))...)
}

func (w *NamespaceStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	return w.s.Delete(ctx, key, append(opts, DeleteNamespace(w.ns))...)
}

func (w *NamespaceStore) Exists(ctx context.Context, key string, opts ...ExistsOption) error {
	return w.s.Exists(ctx, key, append(opts, ExistsNamespace(w.ns))...)
}

func (w *NamespaceStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	return w.s.List(ctx, append(opts, ListNamespace(w.ns))...)
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

// type NamespaceWrapper struct{}

//func NewNamespaceWrapper() Wrapper {
//	return &NamespaceWrapper{}
//}

/*
func (w *OmitWrapper) Logf(fn LogfFunc) LogfFunc {
	return func(ctx context.Context, level Level, msg string, args ...interface{}) {
		fn(ctx, level, msg, getArgs(args)...)
	}
}
*/
