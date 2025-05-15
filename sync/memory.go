package sync

import (
	gosync "sync"
	"time"
)

type memorySync struct {
	locks   map[string]*memoryLock
	options Options

	mu gosync.RWMutex
}

type memoryLock struct {
	time    time.Time
	release chan bool
	id      string
	ttl     time.Duration
}

type memoryLeader struct {
	opts   LeaderOptions
	resign func(id string) error
	status chan bool
	id     string
}

func (m *memoryLeader) Resign() error {
	return m.resign(m.id)
}

func (m *memoryLeader) Status() chan bool {
	return m.status
}

func (m *memorySync) Leader(id string, opts ...LeaderOption) (Leader, error) {
	var once gosync.Once
	var options LeaderOptions
	for _, o := range opts {
		o(&options)
	}

	// acquire a lock for the id
	if err := m.Lock(id); err != nil {
		return nil, err
	}

	// return the leader
	return &memoryLeader{
		opts: options,
		id:   id,
		resign: func(id string) error {
			once.Do(func() {
				_ = m.Unlock(id)
			})
			return nil
		},
		// TODO: signal when Unlock is called
		status: make(chan bool, 1),
	}, nil
}

func (m *memorySync) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.options)
	}
	return nil
}

func (m *memorySync) Options() Options {
	return m.options
}

func (m *memorySync) Lock(id string, opts ...LockOption) error {
	// lock our access
	m.mu.Lock()

	var options LockOptions
	for _, o := range opts {
		o(&options)
	}

	lk, ok := m.locks[id]
	if !ok {
		m.locks[id] = &memoryLock{
			id:      id,
			time:    time.Now(),
			ttl:     options.TTL,
			release: make(chan bool),
		}
		// unlock
		m.mu.Unlock()
		return nil
	}

	m.mu.Unlock()

	// set wait time
	var wait <-chan time.Time
	var ttl <-chan time.Time

	// decide if we should wait
	if options.Wait > time.Duration(0) {
		wait = time.After(options.Wait)
	}

	// check the ttl of the lock
	if lk.ttl > time.Duration(0) {
		// time lived for the lock
		live := time.Since(lk.time)

		// set a timer for the leftover ttl
		if live > lk.ttl {
			// release the lock if it expired
			_ = m.Unlock(id)
		} else {
			ttl = time.After(live)
		}
	}

lockLoop:
	for {
		// wait for the lock to be released
		select {
		case <-lk.release:
			m.mu.Lock()

			// someone locked before us
			lk, ok = m.locks[id]
			if ok {
				m.mu.Unlock()
				continue
			}

			// got chance to lock
			m.locks[id] = &memoryLock{
				id:      id,
				time:    time.Now(),
				ttl:     options.TTL,
				release: make(chan bool),
			}

			m.mu.Unlock()

			break lockLoop
		case <-ttl:
			// ttl exceeded
			_ = m.Unlock(id)
			// TODO: check the ttl again above
			ttl = nil
			// try acquire
			continue
		case <-wait:
			return ErrLockTimeout
		}
	}

	return nil
}

func (m *memorySync) Unlock(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lk, ok := m.locks[id]
	// no lock exists
	if !ok {
		return nil
	}

	// delete the lock
	delete(m.locks, id)

	select {
	case <-lk.release:
		return nil
	default:
		close(lk.release)
	}

	return nil
}

func (m *memorySync) String() string {
	return "memory"
}

// NewSync return new memory sync
func NewSync(opts ...Option) Sync {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	return &memorySync{
		options: options,
		locks:   make(map[string]*memoryLock),
	}
}
