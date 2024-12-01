package store

import (
	"context"
	"sync"
	"sync/atomic"

	"go.unistack.org/micro/v3/options"
	"go.unistack.org/micro/v3/util/id"
)

var _ Store = (*noopStore)(nil)

type noopStore struct {
	mu          sync.Mutex
	watchers    map[string]Watcher
	funcRead    FuncRead
	funcWrite   FuncWrite
	funcExists  FuncExists
	funcList    FuncList
	funcDelete  FuncDelete
	opts        Options
	isConnected atomic.Int32
}

func NewStore(opts ...Option) *noopStore {
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

	n.opts.Hooks.EachNext(func(hook options.Hook) {
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

func (n *noopStore) fnRead(ctx context.Context, key string, val interface{}, opts ...ReadOption) error {
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

func (n *noopStore) fnDelete(ctx context.Context, key string, opts ...DeleteOption) error {
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

func (n *noopStore) fnExists(ctx context.Context, key string, opts ...ExistsOption) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (n *noopStore) Write(ctx context.Context, key string, val interface{}, opts ...WriteOption) error {
	if n.opts.LazyConnect {
		if err := n.connect(ctx); err != nil {
			return err
		}
	}
	return n.funcWrite(ctx, key, val, opts...)
}

func (n *noopStore) fnWrite(ctx context.Context, key string, val interface{}, opts ...WriteOption) error {
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
	exit chan bool
	id   string
	ch   chan Event
	opts WatchOptions
}

func (m *noopStore) Watch(ctx context.Context, opts ...WatchOption) (Watcher, error) {
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

	m.mu.Lock()
	m.watchers[w.id] = w
	m.mu.Unlock()

	return w, nil
}

func (w *watcher) Next() (Event, error) {
	return nil, nil
}

func (w *watcher) Stop() {
}
