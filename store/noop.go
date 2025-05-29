package store

import (
	"context"
	"sync"
	"sync/atomic"

	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/util/id"
)

var _ Store = (*noopStore)(nil)

type noopStore struct {
	watchers map[string]Watcher

	funcRead   FuncRead
	funcWrite  FuncWrite
	funcExists FuncExists
	funcList   FuncList
	funcDelete FuncDelete

	opts        Options
	isConnected atomic.Int32

	mu sync.Mutex
}

func (n *noopStore) Live() bool {
	return true
}

func (n *noopStore) Ready() bool {
	return true
}

func (n *noopStore) Health() bool {
	return true
}

func NewStore(opts ...Option) Store {
	options := NewOptions(opts...)
	return &noopStore{opts: options}
}

func (n *noopStore) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	n.funcRead = n.fnRead
	n.funcWrite = n.fnWrite
	n.funcExists = n.fnExists
	n.funcList = n.fnList
	n.funcDelete = n.fnDelete

	n.opts.Hooks.EachPrev(func(hook options.Hook) {
		switch h := hook.(type) {
		case HookRead:
			n.funcRead = h(n.funcRead)
		case HookWrite:
			n.funcWrite = h(n.funcWrite)
		case HookExists:
			n.funcExists = h(n.funcExists)
		case HookList:
			n.funcList = h(n.funcList)
		case HookDelete:
			n.funcDelete = h(n.funcDelete)
		}
	})

	return nil
}

func (n *noopStore) Connect(ctx context.Context) error {
	if n.opts.LazyConnect {
		return nil
	}
	return n.connect(ctx)
}

func (n *noopStore) Disconnect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (n *noopStore) Read(ctx context.Context, key string, val interface{}, opts ...ReadOption) error {
	if n.opts.LazyConnect {
		if err := n.connect(ctx); err != nil {
			return err
		}
	}
	return n.funcRead(ctx, key, val, opts...)
}

func (n *noopStore) fnRead(ctx context.Context, _ string, _ interface{}, _ ...ReadOption) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (n *noopStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	if n.opts.LazyConnect {
		if err := n.connect(ctx); err != nil {
			return err
		}
	}
	return n.funcDelete(ctx, key, opts...)
}

func (n *noopStore) fnDelete(ctx context.Context, _ string, _ ...DeleteOption) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (n *noopStore) Exists(ctx context.Context, key string, opts ...ExistsOption) error {
	if n.opts.LazyConnect {
		if err := n.connect(ctx); err != nil {
			return err
		}
	}
	return n.funcExists(ctx, key, opts...)
}

func (n *noopStore) fnExists(ctx context.Context, _ string, _ ...ExistsOption) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return store.ErrNotFound
}

func (n *noopStore) Write(ctx context.Context, key string, val interface{}, opts ...WriteOption) error {
	if n.opts.LazyConnect {
		if err := n.connect(ctx); err != nil {
			return err
		}
	}
	return n.funcWrite(ctx, key, val, opts...)
}

func (n *noopStore) fnWrite(ctx context.Context, _ string, _ interface{}, _ ...WriteOption) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (n *noopStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	if n.opts.LazyConnect {
		if err := n.connect(ctx); err != nil {
			return nil, err
		}
	}
	return n.funcList(ctx, opts...)
}

func (n *noopStore) fnList(ctx context.Context, opts ...ListOption) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return nil, nil
}

func (n *noopStore) Name() string {
	return n.opts.Name
}

func (n *noopStore) String() string {
	return "noop"
}

func (n *noopStore) Options() Options {
	return n.opts
}

func (n *noopStore) connect(ctx context.Context) error {
	if n.isConnected.CompareAndSwap(0, 1) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

type watcher struct {
	opts WatchOptions
	ch   chan Event
	exit chan bool
	id   string
}

func (n *noopStore) Watch(_ context.Context, opts ...WatchOption) (Watcher, error) {
	id, err := id.New()
	if err != nil {
		return nil, err
	}
	wo, err := NewWatchOptions(opts...)
	if err != nil {
		return nil, err
	}
	// construct the watcher
	w := &watcher{
		exit: make(chan bool),
		ch:   make(chan Event),
		id:   id,
		opts: wo,
	}

	n.mu.Lock()
	n.watchers[w.id] = w
	n.mu.Unlock()

	return w, nil
}

func (w *watcher) Next() (Event, error) {
	return nil, nil
}

func (w *watcher) Stop() {
}
