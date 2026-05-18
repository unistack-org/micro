package sync

import (
	"testing"
	"time"
)

func TestNewSync(t *testing.T) {
	s := NewSync()
	if s == nil {
		t.Fatal("NewSync returned nil")
	}
}

func TestNewSyncWithOptions(t *testing.T) {
	s := NewSync(
		Nodes("localhost:1234"),
		Prefix("test"),
	)
	if s == nil {
		t.Fatal("NewSync returned nil")
	}
}

func TestString(t *testing.T) {
	s := NewSync()
	if s.String() != "memory" {
		t.Fatalf("expected 'memory', got %q", s.String())
	}
}

func TestInit(t *testing.T) {
	s := NewSync()
	if err := s.Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}
}

func TestInitWithOptions(t *testing.T) {
	s := NewSync()
	if err := s.Init(Nodes("localhost:9999"), Prefix("pfx")); err != nil {
		t.Fatalf("Init() with options error: %v", err)
	}
}

func TestOptions(t *testing.T) {
	s := NewSync(Nodes("a", "b"))
	opts := s.Options()
	if len(opts.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(opts.Nodes))
	}
}

func TestNewOptions(t *testing.T) {
	opts := NewOptions(Nodes("x"), Prefix("p"))
	if len(opts.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(opts.Nodes))
	}
	if opts.Prefix != "p" {
		t.Fatalf("expected prefix 'p', got %q", opts.Prefix)
	}
}

func TestLockUnlock(t *testing.T) {
	s := NewSync()
	if err := s.Lock("key1"); err != nil {
		t.Fatalf("Lock error: %v", err)
	}
	if err := s.Unlock("key1"); err != nil {
		t.Fatalf("Unlock error: %v", err)
	}
}

func TestUnlockNonExistent(t *testing.T) {
	s := NewSync()
	// Unlock on a key that was never locked should not error
	if err := s.Unlock("nonexistent"); err != nil {
		t.Fatalf("Unlock nonexistent error: %v", err)
	}
}

func TestLockTTLOption(t *testing.T) {
	s := NewSync()
	if err := s.Lock("key-ttl", LockTTL(500*time.Millisecond)); err != nil {
		t.Fatalf("Lock with TTL error: %v", err)
	}
	if err := s.Unlock("key-ttl"); err != nil {
		t.Fatalf("Unlock error: %v", err)
	}
}

func TestLockWaitTimeout(t *testing.T) {
	s := NewSync()
	// Acquire first lock
	if err := s.Lock("key-wait"); err != nil {
		t.Fatalf("first Lock error: %v", err)
	}
	// Second lock attempt should time out
	err := s.Lock("key-wait", LockWait(50*time.Millisecond))
	if err != ErrLockTimeout {
		t.Fatalf("expected ErrLockTimeout, got %v", err)
	}
	// Clean up
	if err := s.Unlock("key-wait"); err != nil {
		t.Fatalf("Unlock error: %v", err)
	}
}

func TestLockReleasedByUnlock(t *testing.T) {
	s := NewSync()
	if err := s.Lock("key-release"); err != nil {
		t.Fatalf("Lock error: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		// This should block until key-release is unlocked
		done <- s.Lock("key-release", LockWait(2*time.Second))
	}()

	// Give the goroutine time to start waiting
	time.Sleep(20 * time.Millisecond)

	// Release the first lock
	if err := s.Unlock("key-release"); err != nil {
		t.Fatalf("Unlock error: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("second Lock error after release: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second lock to acquire")
	}

	// Clean up the second acquired lock
	if err := s.Unlock("key-release"); err != nil {
		t.Fatalf("final Unlock error: %v", err)
	}
}

func TestLockTTLExpiry(t *testing.T) {
	s := NewSync()
	// Lock with a very short TTL
	if err := s.Lock("key-exp", LockTTL(30*time.Millisecond)); err != nil {
		t.Fatalf("Lock error: %v", err)
	}

	// Wait for TTL to expire
	time.Sleep(60 * time.Millisecond)

	// A second lock attempt should succeed because TTL expired
	done := make(chan error, 1)
	go func() {
		done <- s.Lock("key-exp", LockWait(500*time.Millisecond))
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Lock after TTL expiry error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for lock after TTL expiry")
	}

	_ = s.Unlock("key-exp")
}

func TestLeaderResign(t *testing.T) {
	s := NewSync()
	leader, err := s.Leader("election1")
	if err != nil {
		t.Fatalf("Leader() error: %v", err)
	}
	if leader == nil {
		t.Fatal("Leader() returned nil")
	}
	if err := leader.Resign(); err != nil {
		t.Fatalf("Resign() error: %v", err)
	}
}

func TestLeaderStatus(t *testing.T) {
	s := NewSync()
	leader, err := s.Leader("election2")
	if err != nil {
		t.Fatalf("Leader() error: %v", err)
	}
	ch := leader.Status()
	if ch == nil {
		t.Fatal("Status() returned nil channel")
	}
	if err := leader.Resign(); err != nil {
		t.Fatalf("Resign() error: %v", err)
	}
}

func TestLeaderResignIdempotent(t *testing.T) {
	s := NewSync()
	leader, err := s.Leader("election3")
	if err != nil {
		t.Fatalf("Leader() error: %v", err)
	}
	// Resign twice — second call should be a no-op (sync.Once)
	if err := leader.Resign(); err != nil {
		t.Fatalf("first Resign() error: %v", err)
	}
	if err := leader.Resign(); err != nil {
		t.Fatalf("second Resign() error: %v", err)
	}
}

func TestLeaderExclusion(t *testing.T) {
	s := NewSync()
	leader1, err := s.Leader("election4")
	if err != nil {
		t.Fatalf("first Leader() error: %v", err)
	}

	// Attempt to get leadership while it's held — should block then timeout
	done := make(chan error, 1)
	go func() {
		_, e := s.(interface {
			Leader(string, ...LeaderOption) (Leader, error)
		}).Leader("election4")
		done <- e
	}()

	// resign so second candidate can take over
	time.Sleep(20 * time.Millisecond)
	if err := leader1.Resign(); err != nil {
		t.Fatalf("Resign() error: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("second Leader() error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second leader election")
	}
}

func TestLockWaitOption(t *testing.T) {
	fn := LockWait(100 * time.Millisecond)
	opts := &LockOptions{}
	fn(opts)
	if opts.Wait != 100*time.Millisecond {
		t.Fatalf("expected Wait=100ms, got %v", opts.Wait)
	}
}

func TestLockTTLOptionOnly(t *testing.T) {
	fn := LockTTL(200 * time.Millisecond)
	opts := &LockOptions{}
	fn(opts)
	if opts.TTL != 200*time.Millisecond {
		t.Fatalf("expected TTL=200ms, got %v", opts.TTL)
	}
}
