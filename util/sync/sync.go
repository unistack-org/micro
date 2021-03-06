// Package sync will sync multiple stores
package sync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ef-ds/deque"
	"github.com/unistack-org/micro/v3/store"
)

// Sync implements a sync in for stores
type Sync interface {
	// Implements the store interface
	store.Store
	// Force a full sync
	Sync() error
}

type syncStore struct {
	storeOpts           store.Options
	pendingWrites       []*deque.Deque
	pendingWriteTickers []*time.Ticker
	syncOpts            Options
	sync.RWMutex
}

// NewSync returns a new Sync
func NewSync(opts ...Option) Sync {
	c := &syncStore{}
	for _, o := range opts {
		o(&c.syncOpts)
	}
	if c.syncOpts.SyncInterval == 0 {
		c.syncOpts.SyncInterval = 1 * time.Minute
	}
	if c.syncOpts.SyncMultiplier == 0 {
		c.syncOpts.SyncMultiplier = 5
	}
	return c
}

func (c *syncStore) Connect(ctx context.Context) error {
	return nil
}

func (c *syncStore) Disconnect(ctx context.Context) error {
	return nil
}

func (c *syncStore) Close(ctx context.Context) error {
	return nil
}

// Init initialises the storeOptions
func (c *syncStore) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&c.storeOpts)
	}
	if len(c.syncOpts.Stores) == 0 {
		return fmt.Errorf("the sync has no stores")
	}
	for _, s := range c.syncOpts.Stores {
		if err := s.Init(); err != nil {
			return fmt.Errorf("Store %s failed to Init(): %w", s.String(), err)
		}
	}
	c.pendingWrites = make([]*deque.Deque, len(c.syncOpts.Stores)-1)
	c.pendingWriteTickers = make([]*time.Ticker, len(c.syncOpts.Stores)-1)
	for i := 0; i < len(c.pendingWrites); i++ {
		c.pendingWrites[i] = deque.New()
		c.pendingWrites[i].Init()
		c.pendingWriteTickers[i] = time.NewTicker(c.syncOpts.SyncInterval * time.Duration(intpow(c.syncOpts.SyncMultiplier, int64(i))))
	}
	go c.syncManager()
	return nil
}

// Name returns the store name
func (c *syncStore) Name() string {
	return c.storeOpts.Name
}

// Options returns the sync's store options
func (c *syncStore) Options() store.Options {
	return c.storeOpts
}

// String returns a printable string describing the sync
func (c *syncStore) String() string {
	backends := make([]string, len(c.syncOpts.Stores))
	for i, s := range c.syncOpts.Stores {
		backends[i] = s.String()
	}
	return fmt.Sprintf("sync %v", backends)
}

func (c *syncStore) List(ctx context.Context, opts ...store.ListOption) ([]string, error) {
	return c.syncOpts.Stores[0].List(ctx, opts...)
}

func (c *syncStore) Exists(ctx context.Context, key string, opts ...store.ExistsOption) error {
	return c.syncOpts.Stores[0].Exists(ctx, key)
}

func (c *syncStore) Read(ctx context.Context, key string, val interface{}, opts ...store.ReadOption) error {
	return c.syncOpts.Stores[0].Read(ctx, key, val, opts...)
}

func (c *syncStore) Write(ctx context.Context, key string, val interface{}, opts ...store.WriteOption) error {
	return c.syncOpts.Stores[0].Write(ctx, key, val, opts...)
}

// Delete removes a key from the sync
func (c *syncStore) Delete(ctx context.Context, key string, opts ...store.DeleteOption) error {
	return c.syncOpts.Stores[0].Delete(ctx, key, opts...)
}

func (c *syncStore) Sync() error {
	return nil
}

type internalRecord struct {
	expiresAt time.Time
	key       string
	value     []byte
}
