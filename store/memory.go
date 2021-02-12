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
	opts  Options
	store *cache.Cache
}

func (m *memoryStore) key(prefix, key string) string {
	return filepath.Join(prefix, key)
}

func (m *memoryStore) prefix(database, table string) string {
	if len(database) == 0 {
		database = m.opts.Database
	}
	if len(table) == 0 {
		table = m.opts.Table
	}
	return filepath.Join(database, table)
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

	if err := m.opts.Codec.Unmarshal(buf, val); err != nil {
		return err
	}

	return nil
}

func (m *memoryStore) delete(prefix, key string) {
	key = m.key(prefix, key)
	m.store.Delete(key)
}

func (m *memoryStore) list(prefix string, limit, offset uint) []string {
	allItems := m.store.Items()
	allKeys := make([]string, len(allItems))
	i := 0

	for k := range allItems {
		if !strings.HasPrefix(k, prefix+"/") {
			continue
		}
		allKeys[i] = strings.TrimPrefix(k, prefix+"/")
		i++
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
	prefix := m.prefix(m.opts.Database, m.opts.Table)
	return m.exists(prefix, key)
}

func (m *memoryStore) Read(ctx context.Context, key string, val interface{}, opts ...ReadOption) error {
	readOpts := NewReadOptions(opts...)
	prefix := m.prefix(readOpts.Database, readOpts.Table)
	return m.get(prefix, key, val)
}

func (m *memoryStore) Write(ctx context.Context, key string, val interface{}, opts ...WriteOption) error {
	writeOpts := NewWriteOptions(opts...)

	prefix := m.prefix(writeOpts.Database, writeOpts.Table)

	key = m.key(prefix, key)

	buf, err := m.opts.Codec.Marshal(val)
	if err != nil {
		return err
	}

	m.store.Set(key, buf, writeOpts.TTL)
	return nil
}

func (m *memoryStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	deleteOptions := NewDeleteOptions(opts...)

	prefix := m.prefix(deleteOptions.Database, deleteOptions.Table)
	m.delete(prefix, key)
	return nil
}

func (m *memoryStore) Options() Options {
	return m.opts
}

func (m *memoryStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	listOptions := NewListOptions(opts...)

	prefix := m.prefix(listOptions.Database, listOptions.Table)
	keys := m.list(prefix, listOptions.Limit, listOptions.Offset)

	if len(listOptions.Prefix) > 0 {
		var prefixKeys []string
		for _, k := range keys {
			if strings.HasPrefix(k, listOptions.Prefix) {
				prefixKeys = append(prefixKeys, k)
			}
		}
		keys = prefixKeys
	}

	if len(listOptions.Suffix) > 0 {
		var suffixKeys []string
		for _, k := range keys {
			if strings.HasSuffix(k, listOptions.Suffix) {
				suffixKeys = append(suffixKeys, k)
			}
		}
		keys = suffixKeys
	}

	return keys, nil
}
