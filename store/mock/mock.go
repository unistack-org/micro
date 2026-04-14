package mock

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.unistack.org/micro/v4/store"
)

// ExpectedWrite represents an expected Write operation
type ExpectedWrite struct {
	key       string
	value     any
	ttl       time.Duration
	metadata  map[string]string //nolint: unused
	namespace string
	times     int
	called    int
	mutex     sync.Mutex
	err       error
}

func (e *ExpectedWrite) match(key string, val any, opts ...store.WriteOption) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check key match
	if e.key != "" && e.key != key {
		return false
	}

	// Check value match
	if e.value != nil && !reflect.DeepEqual(e.value, val) {
		return false
	}

	// Check options
	options := store.NewWriteOptions(opts...)
	if e.ttl > 0 && e.ttl != options.TTL {
		return false
	}
	if e.namespace != "" && e.namespace != options.Namespace {
		return false
	}

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedRead represents an expected Read operation
type ExpectedRead struct {
	key    string
	value  any
	times  int
	called int
	mutex  sync.Mutex
	err    error
}

func (e *ExpectedRead) match(key string, opts ...store.ReadOption) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check key match
	if e.key != "" && e.key != key {
		return false
	}

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedDelete represents an expected Delete operation
type ExpectedDelete struct {
	key    string
	times  int
	called int
	mutex  sync.Mutex
	err    error
}

func (e *ExpectedDelete) match(key string, opts ...store.DeleteOption) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check key match
	if e.key != "" && e.key != key {
		return false
	}

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedExists represents an expected Exists operation
type ExpectedExists struct {
	key    string
	times  int
	called int
	mutex  sync.Mutex
	err    error
}

func (e *ExpectedExists) match(key string, opts ...store.ExistsOption) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check key match
	if e.key != "" && e.key != key {
		return false
	}

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedList represents an expected List operation
type ExpectedList struct {
	times  int
	called int
	mutex  sync.Mutex
	err    error
	keys   []string
}

func (e *ExpectedList) match(opts ...store.ListOption) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedConnect represents an expected Connect operation
type ExpectedConnect struct {
	times  int
	called int
	mutex  sync.Mutex
	err    error
}

func (e *ExpectedConnect) match() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	// For expectations without Times specified, allow only one call
	if e.times == 0 && e.called >= 1 {
		return false
	}

	e.called++
	return true
}

// ExpectedDisconnect represents an expected Disconnect operation
type ExpectedDisconnect struct {
	times  int
	called int
	mutex  sync.Mutex
	err    error
}

func (e *ExpectedDisconnect) match() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check if we've exceeded the expected times
	if e.times > 0 && e.called >= e.times {
		return false
	}

	// For expectations without Times specified, allow only one call
	if e.times == 0 && e.called >= 1 {
		return false
	}

	e.called++
	return true
}

// Store is a mock implementation of the Store interface for testing
type Store struct {
	expectedWrites      []*ExpectedWrite
	expectedReads       []*ExpectedRead
	expectedDeletes     []*ExpectedDelete
	expectedExists      []*ExpectedExists
	expectedLists       []*ExpectedList
	expectedConnects    []*ExpectedConnect
	expectedDisconnects []*ExpectedDisconnect

	data     map[string]any
	exists   map[string]bool
	ttls     map[string]time.Time // key -> expiration time
	metadata map[string]map[string]string
	err      error
	opts     store.Options
	mutex    sync.RWMutex
}

// NewStore creates a new mock store
func NewStore(opts ...store.Option) *Store {
	options := store.NewOptions(opts...)
	return &Store{
		data:     make(map[string]any),
		exists:   make(map[string]bool),
		ttls:     make(map[string]time.Time),
		metadata: make(map[string]map[string]string),
		opts:     options,
	}
}

// ExpectWrite creates an expectation for a Write operation
func (m *Store) ExpectWrite(key string) *ExpectedWrite {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedWrite{key: key}
	m.expectedWrites = append(m.expectedWrites, exp)
	return exp
}

// ExpectRead creates an expectation for a Read operation
func (m *Store) ExpectRead(key string) *ExpectedRead {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedRead{key: key}
	m.expectedReads = append(m.expectedReads, exp)
	return exp
}

// ExpectDelete creates an expectation for a Delete operation
func (m *Store) ExpectDelete(key string) *ExpectedDelete {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedDelete{key: key}
	m.expectedDeletes = append(m.expectedDeletes, exp)
	return exp
}

// ExpectExists creates an expectation for an Exists operation
func (m *Store) ExpectExists(key string) *ExpectedExists {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedExists{key: key}
	m.expectedExists = append(m.expectedExists, exp)
	return exp
}

// ExpectList creates an expectation for a List operation
func (m *Store) ExpectList() *ExpectedList {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedList{}
	m.expectedLists = append(m.expectedLists, exp)
	return exp
}

// ExpectConnect creates an expectation for a Connect operation
func (m *Store) ExpectConnect() *ExpectedConnect {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedConnect{}
	m.expectedConnects = append(m.expectedConnects, exp)
	return exp
}

// ExpectDisconnect creates an expectation for a Disconnect operation
func (m *Store) ExpectDisconnect() *ExpectedDisconnect {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedDisconnect{}
	m.expectedDisconnects = append(m.expectedDisconnects, exp)
	return exp
}

// WithValue sets the value to return for expected operations
func (e *ExpectedWrite) WithValue(val any) *ExpectedWrite {
	e.value = val
	return e
}

// WithTTL sets the TTL for expected Write operations
func (e *ExpectedWrite) WithTTL(ttl time.Duration) *ExpectedWrite {
	e.ttl = ttl
	return e
}

// WithNamespace sets the namespace for expected operations
func (e *ExpectedWrite) WithNamespace(ns string) *ExpectedWrite {
	e.namespace = ns
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedWrite) Times(n int) *ExpectedWrite {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedWrite) WillReturnError(err error) *ExpectedWrite {
	e.err = err
	return e
}

// WithValue sets the value to return for expected Read operations
func (e *ExpectedRead) WithValue(val any) *ExpectedRead {
	e.value = val
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedRead) Times(n int) *ExpectedRead {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedRead) WillReturnError(err error) *ExpectedRead {
	e.err = err
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedDelete) Times(n int) *ExpectedDelete {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedDelete) WillReturnError(err error) *ExpectedDelete {
	e.err = err
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedExists) Times(n int) *ExpectedExists {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedExists) WillReturnError(err error) *ExpectedExists {
	e.err = err
	return e
}

// WillReturn sets the keys to return for List operations
func (e *ExpectedList) WillReturn(keys ...string) *ExpectedList {
	e.keys = keys
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedList) Times(n int) *ExpectedList {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedList) WillReturnError(err error) *ExpectedList {
	e.err = err
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedConnect) Times(n int) *ExpectedConnect {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedConnect) WillReturnError(err error) *ExpectedConnect {
	e.err = err
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedDisconnect) Times(n int) *ExpectedDisconnect {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedDisconnect) WillReturnError(err error) *ExpectedDisconnect {
	e.err = err
	return e
}

// checkTTL checks if a key has expired
func (m *Store) checkTTL(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if exp, ok := m.ttls[key]; ok {
		if time.Now().After(exp) {
			delete(m.data, key)
			delete(m.exists, key)
			delete(m.ttls, key)
			delete(m.metadata, key)
			return false
		}
	}
	return true
}

// FastForward decrements all TTLs by the given duration
func (m *Store) FastForward(d time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for key, exp := range m.ttls {
		// Calculate remaining time before fast forward
		remaining := time.Until(exp)
		if remaining <= 0 {
			// Already expired, remove it
			delete(m.data, key)
			delete(m.exists, key)
			delete(m.ttls, key)
			delete(m.metadata, key)
		} else {
			// Apply fast forward
			newRemaining := remaining - d
			if newRemaining <= 0 {
				// Would expire after fast forward, remove it
				delete(m.data, key)
				delete(m.exists, key)
				delete(m.ttls, key)
				delete(m.metadata, key)
			} else {
				// Update expiration time
				m.ttls[key] = now.Add(newRemaining)
			}
		}
	}
}

// Name returns store name
func (m *Store) Name() string {
	return m.opts.Name
}

// Init initializes the mock store
func (m *Store) Init(opts ...store.Option) error {
	if m.err != nil {
		return m.err
	}
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

// Connect is used when store needs to be connected
func (m *Store) Connect(ctx context.Context) error {
	// Find matching expectation
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, exp := range m.expectedConnects {
		if exp.match() {
			// Remove the expectation if it has been fully used
			if exp.times > 0 && exp.called >= exp.times {
				// Remove this expectation from the list
				m.expectedConnects = append(m.expectedConnects[:i], m.expectedConnects[i+1:]...)
			}
			return exp.err
		}
	}

	// If no expectation matched, use default behavior
	return m.err
}

// Options returns the current options
func (m *Store) Options() store.Options {
	return m.opts
}

// Exists checks that key exists in store
func (m *Store) Exists(ctx context.Context, key string, opts ...store.ExistsOption) error {
	if m.err != nil {
		return m.err
	}

	// Check TTL first
	if !m.checkTTL(key) {
		return store.ErrNotFound
	}

	// Find matching expectation
	m.mutex.Lock()
	for _, exp := range m.expectedExists {
		if exp.match(key, opts...) {
			m.mutex.Unlock()
			if exp.err != nil {
				return exp.err
			}
			if !m.exists[key] {
				return store.ErrNotFound
			}
			return nil
		}
	}
	m.mutex.Unlock()

	// If no expectation matched, use default behavior
	if !m.exists[key] {
		return store.ErrNotFound
	}
	return nil
}

// Read reads a single key name to provided value with optional ReadOptions
func (m *Store) Read(ctx context.Context, key string, val any, opts ...store.ReadOption) error {
	if m.err != nil {
		return m.err
	}

	// Check TTL first
	if !m.checkTTL(key) {
		return store.ErrNotFound
	}

	// Find matching expectation
	m.mutex.Lock()
	for _, exp := range m.expectedReads {
		if exp.match(key, opts...) {
			m.mutex.Unlock()
			if exp.err != nil {
				return exp.err
			}
			if !m.exists[key] {
				return store.ErrNotFound
			}

			// Copy the value from expected or actual data
			data := exp.value
			if data == nil {
				data = m.data[key]
			}

			if data != nil {
				// Simple type conversion for testing
				if target, ok := val.(*any); ok {
					*target = data
				} else if target, ok := val.(*string); ok {
					if s, ok := data.(string); ok {
						*target = s
					} else {
						*target = fmt.Sprintf("%v", data)
					}
				} else if target, ok := val.(*int); ok {
					if i, ok := data.(int); ok {
						*target = i
					}
				}
			}
			return nil
		}
	}
	m.mutex.Unlock()

	// If no expectation matched, use default behavior
	if !m.exists[key] {
		return store.ErrNotFound
	}

	if data, ok := m.data[key]; ok {
		if target, ok := val.(*any); ok {
			*target = data
		} else if target, ok := val.(*string); ok {
			if s, ok := data.(string); ok {
				*target = s
			} else {
				*target = fmt.Sprintf("%v", data)
			}
		} else if target, ok := val.(*int); ok {
			if i, ok := data.(int); ok {
				*target = i
			}
		}
	}
	return nil
}

// Write writes a value to key name to the store with optional WriteOption
func (m *Store) Write(ctx context.Context, key string, val any, opts ...store.WriteOption) error {
	if m.err != nil {
		return m.err
	}

	// Find matching expectation
	m.mutex.Lock()
	for _, exp := range m.expectedWrites {
		if exp.match(key, val, opts...) {
			m.mutex.Unlock()
			if exp.err != nil {
				return exp.err
			}

			// Apply the write operation
			m.mutex.Lock()
			m.data[key] = val
			m.exists[key] = true

			// Handle TTL
			options := store.NewWriteOptions(opts...)
			if options.TTL > 0 {
				m.ttls[key] = time.Now().Add(options.TTL)
			} else {
				delete(m.ttls, key) // Remove TTL if not set
			}

			// Handle metadata
			if options.Metadata != nil {
				m.metadata[key] = make(map[string]string)
				for k, v := range options.Metadata {
					// Convert []string to string by joining with comma
					if len(v) > 0 {
						m.metadata[key][k] = strings.Join(v, ",")
					} else {
						m.metadata[key][k] = ""
					}
				}
			}
			m.mutex.Unlock()
			return nil
		}
	}
	m.mutex.Unlock()

	// If no expectation matched, use default behavior
	m.mutex.Lock()
	m.data[key] = val
	m.exists[key] = true

	options := store.NewWriteOptions(opts...)
	if options.TTL > 0 {
		m.ttls[key] = time.Now().Add(options.TTL)
	} else {
		delete(m.ttls, key)
	}

	if options.Metadata != nil {
		m.metadata[key] = make(map[string]string)
		for k, v := range options.Metadata {
			// Convert []string to string by joining with comma
			if len(v) > 0 {
				m.metadata[key][k] = strings.Join(v, ",")
			} else {
				m.metadata[key][k] = ""
			}
		}
	}
	m.mutex.Unlock()
	return nil
}

// Delete removes the record with the corresponding key from the store
func (m *Store) Delete(ctx context.Context, key string, opts ...store.DeleteOption) error {
	if m.err != nil {
		return m.err
	}

	// Find matching expectation
	m.mutex.Lock()
	for _, exp := range m.expectedDeletes {
		if exp.match(key, opts...) {
			m.mutex.Unlock()
			if exp.err != nil {
				return exp.err
			}

			m.mutex.Lock()
			delete(m.data, key)
			delete(m.exists, key)
			delete(m.ttls, key)
			delete(m.metadata, key)
			m.mutex.Unlock()
			return nil
		}
	}
	m.mutex.Unlock()

	// If no expectation matched, use default behavior
	m.mutex.Lock()
	delete(m.data, key)
	delete(m.exists, key)
	delete(m.ttls, key)
	delete(m.metadata, key)
	m.mutex.Unlock()
	return nil
}

// List returns any keys that match, or an empty list with no error if none matched
func (m *Store) List(ctx context.Context, opts ...store.ListOption) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Find matching expectation
	m.mutex.Lock()
	for _, exp := range m.expectedLists {
		if exp.match(opts...) {
			m.mutex.Unlock()
			if exp.err != nil {
				return nil, exp.err
			}
			return exp.keys, nil
		}
	}
	m.mutex.Unlock()

	// If no expectation matched, return actual keys
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var keys []string
	for key := range m.data {
		// Check TTL
		if exp, ok := m.ttls[key]; ok {
			if time.Now().After(exp) {
				continue // Skip expired keys
			}
		}
		keys = append(keys, key)
	}

	// Apply list options filtering
	options := store.NewListOptions(opts...)
	if options.Prefix != "" {
		var filtered []string
		for _, key := range keys {
			if len(key) >= len(options.Prefix) && key[:len(options.Prefix)] == options.Prefix {
				filtered = append(filtered, key)
			}
		}
		keys = filtered
	}

	if options.Suffix != "" {
		var filtered []string
		for _, key := range keys {
			if len(key) >= len(options.Suffix) && key[len(key)-len(options.Suffix):] == options.Suffix {
				filtered = append(filtered, key)
			}
		}
		keys = filtered
	}

	// Apply limit and offset
	if options.Limit > 0 && int(options.Limit) < len(keys) {
		end := min(int(options.Offset)+int(options.Limit), len(keys))
		if int(options.Offset) < len(keys) {
			keys = keys[options.Offset:end]
		} else {
			keys = []string{}
		}
	} else if options.Offset > 0 && int(options.Offset) < len(keys) {
		keys = keys[options.Offset:]
	} else if options.Offset >= uint(len(keys)) {
		keys = []string{}
	}

	return keys, nil
}

// Disconnect disconnects the mock store
func (m *Store) Disconnect(ctx context.Context) error {
	// Find matching expectation
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, exp := range m.expectedDisconnects {
		if exp.match() {
			// Remove the expectation if it has been fully used
			if exp.times > 0 && exp.called >= exp.times {
				// Remove this expectation from the list
				m.expectedDisconnects = append(m.expectedDisconnects[:i], m.expectedDisconnects[i+1:]...)
			}
			return exp.err
		}
	}

	// If no expectation matched, use default behavior
	return m.err
}

// String returns the name of the implementation
func (m *Store) String() string {
	return "mock"
}

// Watch returns events watcher
func (m *Store) Watch(ctx context.Context, opts ...store.WatchOption) (store.Watcher, error) {
	if m.err != nil {
		return nil, m.err
	}
	return NewWatcher(), nil
}

// Live returns store liveness
func (m *Store) Live() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.err == nil
}

// Ready returns store readiness
func (m *Store) Ready() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.err == nil
}

// Health returns store health
func (m *Store) Health() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.err == nil
}

// SetError sets an error that will be returned by all operations
func (m *Store) SetError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.err = err
}

// ExpectationsWereMet checks that all expected operations were called the expected number of times
func (m *Store) ExpectationsWereMet() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, exp := range m.expectedWrites {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected write for key %s to be called %d times, but was called %d times", exp.key, exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedReads {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected read for key %s to be called %d times, but was called %d times", exp.key, exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedDeletes {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected delete for key %s to be called %d times, but was called %d times", exp.key, exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedExists {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected exists for key %s to be called %d times, but was called %d times", exp.key, exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedLists {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected list to be called %d times, but was called %d times", exp.times, exp.called)
		}
	}

	return nil
}

// Watcher is a mock implementation of the Watcher interface
type Watcher struct {
	events chan store.Event
	stop   chan bool
}

// NewWatcher creates a new mock watcher
func NewWatcher() *Watcher {
	return &Watcher{
		events: make(chan store.Event, 1),
		stop:   make(chan bool, 1),
	}
}

// Next is a blocking call that returns the next event
func (mw *Watcher) Next() (store.Event, error) {
	select {
	case event := <-mw.events:
		return event, nil
	case <-mw.stop:
		return nil, store.ErrWatcherStopped
	}
}

// Stop stops the watcher
func (mw *Watcher) Stop() {
	select {
	case mw.stop <- true:
	default:
	}
}

// SendEvent sends an event to the watcher (for testing purposes)
func (mw *Watcher) SendEvent(event store.Event) {
	select {
	case mw.events <- event:
	default:
		// If channel is full, drop the event
	}
}
