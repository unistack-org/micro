package store

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// NewStore returns a memory store
func NewStore(opts ...Option) Store {
	return &memoryStore{
		opts:  NewOptions(opts...),
		store: cache.New(cache.NoExpiration, 5*time.Minute),
	}
}

func (m *memoryStore) Connect(ctx context.Context) error {
	return nil
}

func (m *memoryStore) Disconnect(ctx context.Context) error {
	m.store.Flush()
	return nil
}

type memoryStore struct {
	store *cache.Cache
	opts  Options
}

func (m *memoryStore) key(prefix, key string) string {
	return filepath.Join(prefix, key)
}

func (m *memoryStore) exists(prefix, key string) error {
	key = m.key(prefix, key)

	_, found := m.store.Get(key)
	if !found {
		return ErrNotFound
	}

	return nil
}

func (m *memoryStore) get(prefix, key string, val interface{}) error {
	key = m.key(prefix, key)

	r, found := m.store.Get(key)
	if !found {
		return ErrNotFound
	}

	buf, ok := r.([]byte)
	if !ok {
		return ErrNotFound
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

func (m *memoryStore) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

func (m *memoryStore) String() string {
	return "memory"
}

func (m *memoryStore) Name() string {
	return m.opts.Name
}

func (m *memoryStore) Exists(ctx context.Context, key string, opts ...ExistsOption) error {
	options := NewExistsOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}
	return m.exists(options.Namespace, key)
}

func (m *memoryStore) Read(ctx context.Context, key string, val interface{}, opts ...ReadOption) error {
	options := NewReadOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}
	return m.get(options.Namespace, key, val)
}

func (m *memoryStore) Write(ctx context.Context, key string, val interface{}, opts ...WriteOption) error {
	options := NewWriteOptions(opts...)
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

func (m *memoryStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	options := NewDeleteOptions(opts...)
	if options.Namespace == "" {
		options.Namespace = m.opts.Namespace
	}

	m.delete(options.Namespace, key)
	return nil
}

func (m *memoryStore) Options() Options {
	return m.opts
}

func (m *memoryStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	options := NewListOptions(opts...)
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
