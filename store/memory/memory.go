package memory

import (
	"context"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"
	"go.unistack.org/micro/v3/options"
	"go.unistack.org/micro/v3/store"
)

// NewStore returns a memory store
func NewStore(opts ...store.Option) store.Store {
	return &memoryStore{
		opts:  store.NewOptions(opts...),
		store: cache.New(cache.NoExpiration, 5*time.Minute),
	}
}

func (m *memoryStore) Connect(ctx context.Context) error {
	if m.opts.LazyConnect {
		return nil
	}
	return m.connect(ctx)
}

func (m *memoryStore) Disconnect(ctx context.Context) error {
	m.store.Flush()
	return nil
}

type memoryStore struct {
	funcRead    store.FuncRead
	funcWrite   store.FuncWrite
	funcExists  store.FuncExists
	funcList    store.FuncList
	funcDelete  store.FuncDelete
	store       *cache.Cache
	opts        store.Options
	isConnected atomic.Int32
}

func (m *memoryStore) key(prefix, key string) string {
	return prefix + m.opts.Separator + key
}

func (m *memoryStore) exists(prefix, key string) error {
	key = m.key(prefix, key)
	_, found := m.store.Get(key)
	if !found {
		return store.ErrNotFound
	}

	return nil
}

func (m *memoryStore) get(prefix, key string, val interface{}) error {
	key = m.key(prefix, key)

	r, found := m.store.Get(key)
	if !found {
		return store.ErrNotFound
	}

	buf, ok := r.([]byte)
	if !ok {
		return store.ErrNotFound
	}

	return m.opts.Codec.Unmarshal(buf, val)
}

func (m *memoryStore) delete(prefix, key string) {
	key = m.key(prefix, key)
	m.store.Delete(key)
}

func (m *memoryStore) list(prefix string, limit, offset uint) []string {
	allItems := m.store.Items()
	allKeys := make([]string, 0, len(allItems))

	for k := range allItems {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		k = strings.TrimPrefix(k, prefix)
		if k[0] == '/' {
			k = k[1:]
		}
		allKeys = append(allKeys, k)
	}

	if limit != 0 || offset != 0 {
		sort.Slice(allKeys, func(i, j int) bool { return allKeys[i] < allKeys[j] })
		sort.Slice(allKeys, func(i, j int) bool { return allKeys[i] < allKeys[j] })
		end := len(allKeys)
		if limit > 0 {
			calcLimit := int(offset + limit)
			if calcLimit < end {
				end = calcLimit
			}
		}

		if int(offset) >= end {
			return nil
		}
		return allKeys[offset:end]
	}
	return allKeys
}

func (m *memoryStore) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}

	m.funcRead = m.fnRead
	m.funcWrite = m.fnWrite
	m.funcExists = m.fnExists
	m.funcList = m.fnList
	m.funcDelete = m.fnDelete

	m.opts.Hooks.EachNext(func(hook options.Hook) {
		switch h := hook.(type) {
		case store.HookRead:
			m.funcRead = h(m.funcRead)
		case store.HookWrite:
			m.funcWrite = h(m.funcWrite)
		case store.HookExists:
			m.funcExists = h(m.funcExists)
		case store.HookList:
			m.funcList = h(m.funcList)
		case store.HookDelete:
			m.funcDelete = h(m.funcDelete)
		}
	})

	return nil
}

func (m *memoryStore) String() string {
	return "memory"
}

func (m *memoryStore) Name() string {
	return m.opts.Name
}

func (m *memoryStore) Exists(ctx context.Context, key string, opts ...store.ExistsOption) error {
	if m.opts.LazyConnect {
		if err := m.connect(ctx); err != nil {
			return err
		}
	}
	return m.funcExists(ctx, key, opts...)
}

func (m *memoryStore) fnExists(ctx context.Context, key string, opts ...store.ExistsOption) error {
	options := store.NewExistsOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}
	return m.exists(options.Namespace, key)
}

func (m *memoryStore) Read(ctx context.Context, key string, val interface{}, opts ...store.ReadOption) error {
	if m.opts.LazyConnect {
		if err := m.connect(ctx); err != nil {
			return err
		}
	}
	return m.funcRead(ctx, key, val, opts...)
}

func (m *memoryStore) fnRead(ctx context.Context, key string, val interface{}, opts ...store.ReadOption) error {
	options := store.NewReadOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}
	return m.get(options.Namespace, key, val)
}

func (m *memoryStore) Write(ctx context.Context, key string, val interface{}, opts ...store.WriteOption) error {
	if m.opts.LazyConnect {
		if err := m.connect(ctx); err != nil {
			return err
		}
	}
	return m.funcWrite(ctx, key, val, opts...)
}

func (m *memoryStore) fnWrite(ctx context.Context, key string, val interface{}, opts ...store.WriteOption) error {
	options := store.NewWriteOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}
	if options.TTL == 0 {
		options.TTL = cache.NoExpiration
	}

	key = m.key(options.Namespace, key)

	buf, err := m.opts.Codec.Marshal(val)
	if err != nil {
		return err
	}

	m.store.Set(key, buf, options.TTL)
	return nil
}

func (m *memoryStore) Delete(ctx context.Context, key string, opts ...store.DeleteOption) error {
	if m.opts.LazyConnect {
		if err := m.connect(ctx); err != nil {
			return err
		}
	}
	return m.funcDelete(ctx, key, opts...)
}

func (m *memoryStore) fnDelete(ctx context.Context, key string, opts ...store.DeleteOption) error {
	options := store.NewDeleteOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}

	m.delete(options.Namespace, key)
	return nil
}

func (m *memoryStore) Options() store.Options {
	return m.opts
}

func (m *memoryStore) List(ctx context.Context, opts ...store.ListOption) ([]string, error) {
	if m.opts.LazyConnect {
		if err := m.connect(ctx); err != nil {
			return nil, err
		}
	}
	return m.funcList(ctx, opts...)
}

func (m *memoryStore) fnList(ctx context.Context, opts ...store.ListOption) ([]string, error) {
	options := store.NewListOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}

	keys := m.list(options.Namespace, options.Limit, options.Offset)

	if len(options.Prefix) > 0 {
		var prefixKeys []string
		for _, k := range keys {
			if strings.HasPrefix(k, options.Prefix) {
				prefixKeys = append(prefixKeys, k)
			}
		}
		keys = prefixKeys
	}

	if len(options.Suffix) > 0 {
		var suffixKeys []string
		for _, k := range keys {
			if strings.HasSuffix(k, options.Suffix) {
				suffixKeys = append(suffixKeys, k)
			}
		}
		keys = suffixKeys
	}

	return keys, nil
}

func (m *memoryStore) connect(ctx context.Context) error {
	m.isConnected.CompareAndSwap(0, 1)
	return nil
}
