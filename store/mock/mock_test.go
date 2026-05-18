package mock

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v5/store"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	s := NewStore()

	// Test Write with expectation
	s.ExpectWrite("test_key").WithValue("test_value")
	err := s.Write(ctx, "test_key", "test_value")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Test Read with expectation
	s.ExpectRead("test_key").WithValue("test_value")
	var value any
	err = s.Read(ctx, "test_key", &value)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if value != "test_value" {
		t.Fatalf("Expected 'test_value', got %v", value)
	}

	// Test Read with string
	s.ExpectRead("test_key")
	var strValue string
	err = s.Read(ctx, "test_key", &strValue)
	if err != nil {
		t.Fatalf("Read string failed: %v", err)
	}
	if strValue != "test_value" {
		t.Fatalf("Expected 'test_value', got %s", strValue)
	}

	// Test Write and Read integer with TTL
	s.ExpectWrite("int_key").WithValue(42).WithTTL(5 * time.Second)
	err = s.Write(ctx, "int_key", 42, store.WriteTTL(5*time.Second))
	if err != nil {
		t.Fatalf("Write int failed: %v", err)
	}

	s.ExpectRead("int_key")
	var intValue int
	err = s.Read(ctx, "int_key", &intValue)
	if err != nil {
		t.Fatalf("Read int failed: %v", err)
	}
	if intValue != 42 {
		t.Fatalf("Expected 42, got %d", intValue)
	}

	// Test Exists with expectation
	s.ExpectExists("test_key")
	err = s.Exists(ctx, "test_key")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}

	// Test List with expectation
	s.ExpectList().WillReturn("test_key", "another_key")
	keys, err := s.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(keys))
	}

	// Test Delete with expectation
	s.ExpectDelete("test_key")
	err = s.Delete(ctx, "test_key")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Test that deleted key doesn't exist
	s.ExpectExists("test_key").WillReturnError(store.ErrNotFound)
	err = s.Exists(ctx, "test_key")
	if err == nil {
		t.Fatalf("Expected store.ErrNotFound after delete")
	}

	// Test error handling
	s.ExpectExists("nonexistent").WillReturnError(store.ErrNotFound)
	err = s.Exists(ctx, "nonexistent")
	if err != store.ErrNotFound {
		t.Fatalf("Expected store.ErrNotFound, got %v", err)
	}

	// Verify all expectations were met
	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestStoreFastForward(t *testing.T) {
	ctx := context.Background()
	s := NewStore()

	// Write with TTL
	s.ExpectWrite("ttl_key").WithValue("ttl_value").WithTTL(100 * time.Millisecond)
	err := s.Write(ctx, "ttl_key", "ttl_value", store.WriteTTL(100*time.Millisecond))
	if err != nil {
		t.Fatalf("Write with TTL failed: %v", err)
	}

	// Check key exists before TTL expires
	s.ExpectRead("ttl_key")
	var value string
	err = s.Read(ctx, "ttl_key", &value)
	if err != nil {
		t.Fatalf("Read before TTL failed: %v", err)
	}
	if value != "ttl_value" {
		t.Fatalf("Expected 'ttl_value', got %s", value)
	}

	// Fast forward by 50ms - key should still exist
	s.FastForward(50 * time.Millisecond)

	s.ExpectRead("ttl_key")
	err = s.Read(ctx, "ttl_key", &value)
	if err != nil {
		t.Fatalf("Read after 50ms fast forward failed: %v", err)
	}
	if value != "ttl_value" {
		t.Fatalf("Expected 'ttl_value' after 50ms, got %s", value)
	}

	// Fast forward by another 60ms (total 110ms) - key should expire
	s.FastForward(60 * time.Millisecond)

	s.ExpectRead("ttl_key").WillReturnError(store.ErrNotFound)
	err = s.Read(ctx, "ttl_key", &value)
	if err != store.ErrNotFound {
		t.Fatalf("Expected store.ErrNotFound after TTL, got %v", err)
	}

	// Test FastForward on already expired keys
	s.ExpectWrite("ttl_key2").WithValue("ttl_value2").WithTTL(10 * time.Millisecond)
	err = s.Write(ctx, "ttl_key2", "ttl_value2", store.WriteTTL(10*time.Millisecond))
	if err != nil {
		t.Fatalf("Write with TTL failed: %v", err)
	}

	// Fast forward by 20ms - key should expire immediately
	s.FastForward(20 * time.Millisecond)

	s.ExpectRead("ttl_key2").WillReturnError(store.ErrNotFound)
	err = s.Read(ctx, "ttl_key2", &value)
	if err != store.ErrNotFound {
		t.Fatalf("Expected store.ErrNotFound after immediate expiration, got %v", err)
	}

	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestStoreWithOptions(t *testing.T) {
	s := NewStore(store.Name("test_mock"), store.Namespace("test_ns"))

	if s.Name() != "test_mock" {
		t.Fatalf("Expected name 'test_mock', got %s", s.Name())
	}

	opts := s.Options()
	if opts.Namespace != "test_ns" {
		t.Fatalf("Expected namespace 'test_ns', got %s", opts.Namespace)
	}
}

func TestWatcher(t *testing.T) {
	watcher := NewWatcher()

	// Test Stop
	watcher.Stop()

	// Test Next after stop
	_, err := watcher.Next()
	if err != store.ErrWatcherStopped {
		t.Fatalf("Expected store.ErrWatcherStopped, got %v", err)
	}
}

func TestStoreHealth(t *testing.T) {
	s := NewStore()

	if !s.Live() {
		t.Fatal("Expected Live() to return true")
	}

	if !s.Ready() {
		t.Fatal("Expected Ready() to return true")
	}

	if !s.Health() {
		t.Fatal("Expected Health() to return true")
	}

	// Test that Live, Ready, Health return false when error is set
	s.SetError(store.ErrNotConnected)

	if s.Live() {
		t.Fatal("Expected Live() to return false when error is set")
	}

	if s.Ready() {
		t.Fatal("Expected Ready() to return false when error is set")
	}

	if s.Health() {
		t.Fatal("Expected Health() to return false when error is set")
	}
}

func TestStoreSetError(t *testing.T) {
	s := NewStore()

	// Initially should be healthy
	if !s.Live() || !s.Ready() || !s.Health() {
		t.Fatal("Expected store to be healthy initially")
	}

	// Set error and check health
	s.SetError(store.ErrNotFound)

	if s.Live() || s.Ready() || s.Health() {
		t.Fatal("Expected store to be unhealthy when error is set")
	}

	// Reset error and check health
	s.SetError(nil)

	if !s.Live() || !s.Ready() || !s.Health() {
		t.Fatal("Expected store to be healthy after error reset")
	}
}

func TestStoreConnectDisconnect(t *testing.T) {
	s := NewStore()

	// Test Connect with expectation
	s.ExpectConnect()
	err := s.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Test Disconnect with expectation
	s.ExpectDisconnect()
	err = s.Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}

	// Test error propagation for Connect
	s.ExpectConnect().WillReturnError(store.ErrNotConnected)
	err = s.Connect(context.Background())
	if err != store.ErrNotConnected {
		t.Fatalf("Expected store.ErrNotConnected for Connect, got %v", err)
	}

	// Test error propagation for Disconnect
	s.ExpectDisconnect().WillReturnError(store.ErrNotConnected)
	err = s.Disconnect(context.Background())
	if err != store.ErrNotConnected {
		t.Fatalf("Expected store.ErrNotConnected for Disconnect, got %v", err)
	}

	// Test multiple calls with Times
	s.ExpectConnect().Times(2)
	err = s.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	err = s.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Test behavior when no expectation set but error is set globally
	s2 := NewStore()
	s2.SetError(store.ErrNotConnected)
	err = s2.Connect(context.Background())
	if err != store.ErrNotConnected {
		t.Fatalf("Expected store.ErrNotConnected for Connect when global error is set, got %v", err)
	}

	err = s2.Disconnect(context.Background())
	if err != store.ErrNotConnected {
		t.Fatalf("Expected store.ErrNotConnected for Disconnect when global error is set, got %v", err)
	}

	// Test behavior when no expectation set and no global error
	s3 := NewStore()
	err = s3.Connect(context.Background())
	if err != nil {
		t.Fatalf("Expected no error for Connect when no expectation and no global error, got %v", err)
	}

	err = s3.Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Expected no error for Disconnect when no expectation and no global error, got %v", err)
	}

	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestStoreTTL(t *testing.T) {
	ctx := context.Background()
	s := NewStore()

	// Test Write with TTL
	s.ExpectWrite("ttl_key").WithValue("ttl_value").WithTTL(100 * time.Millisecond)
	err := s.Write(ctx, "ttl_key", "ttl_value", store.WriteTTL(100*time.Millisecond))
	if err != nil {
		t.Fatalf("Write with TTL failed: %v", err)
	}

	// Read before TTL expires
	s.ExpectRead("ttl_key")
	var value string
	err = s.Read(ctx, "ttl_key", &value)
	if err != nil {
		t.Fatalf("Read before TTL failed: %v", err)
	}
	if value != "ttl_value" {
		t.Fatalf("Expected 'ttl_value', got %s", value)
	}

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Read after TTL expires should return ErrNotFound
	s.ExpectRead("ttl_key").WillReturnError(store.ErrNotFound)
	err = s.Read(ctx, "ttl_key", &value)
	if err != store.ErrNotFound {
		t.Fatalf("Expected store.ErrNotFound after TTL, got %v", err)
	}

	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestStoreExpectedOperations(t *testing.T) {
	ctx := context.Background()
	s := NewStore()

	// Test expected operations with Times
	s.ExpectWrite("once_key").Times(1)
	s.ExpectWrite("twice_key").Times(2)

	err := s.Write(ctx, "once_key", "value1")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	err = s.Write(ctx, "twice_key", "value2")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	err = s.Write(ctx, "twice_key", "value3")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestStoreInit(t *testing.T) {
	s := NewStore()
	if err := s.Init(store.Name("newname")); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if s.Name() != "newname" {
		t.Fatalf("expected name 'newname', got %q", s.Name())
	}
}

func TestStoreInitError(t *testing.T) {
	s := NewStore()
	s.SetError(store.ErrNotFound)
	if err := s.Init(); err == nil {
		t.Fatal("expected error from Init when global error is set")
	}
}

func TestStoreString(t *testing.T) {
	s := NewStore()
	if s.String() != "mock" {
		t.Fatalf("expected 'mock', got %q", s.String())
	}
}

func TestStoreWatch(t *testing.T) {
	ctx := context.Background()
	s := NewStore()

	w, err := s.Watch(ctx)
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestStoreWatchError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.SetError(store.ErrNotConnected)
	_, err := s.Watch(ctx)
	if err == nil {
		t.Fatal("expected error from Watch when global error is set")
	}
}

func TestWatcherSendEvent(t *testing.T) {
	w := NewWatcher()
	// Send an event and receive it
	ev := &mockEvent{key: "k1"}
	w.SendEvent(ev)
	got, err := w.Next()
	if err != nil {
		t.Fatalf("Next returned error: %v", err)
	}
	if got != ev {
		t.Fatal("expected the sent event back")
	}
}

// mockEvent satisfies store.Event for testing.
type mockEvent struct {
	key string
}

func (e *mockEvent) Timestamp() time.Time  { return time.Time{} }
func (e *mockEvent) Error() error          { return nil }
func (e *mockEvent) Type() store.EventType { return store.EventTypeUnknown }

func TestStoreListWithPrefixSuffixLimitOffset(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	_ = s.Write(ctx, "prefix_a", "v1")
	_ = s.Write(ctx, "prefix_b", "v2")
	_ = s.Write(ctx, "other_a_suffix", "v3")
	_ = s.Write(ctx, "other_b_suffix", "v4")
	_ = s.Write(ctx, "prefix_c_suffix", "v5")

	// Prefix filter
	keys, err := s.List(ctx, store.ListPrefix("prefix_"))
	if err != nil {
		t.Fatalf("List with prefix failed: %v", err)
	}
	for _, k := range keys {
		if len(k) < len("prefix_") || k[:len("prefix_")] != "prefix_" {
			t.Fatalf("key %q does not have expected prefix", k)
		}
	}

	// Suffix filter
	keys, err = s.List(ctx, store.ListSuffix("_suffix"))
	if err != nil {
		t.Fatalf("List with suffix failed: %v", err)
	}
	for _, k := range keys {
		if len(k) < len("_suffix") || k[len(k)-len("_suffix"):] != "_suffix" {
			t.Fatalf("key %q does not have expected suffix", k)
		}
	}

	// Limit
	keys, err = s.List(ctx, store.ListLimit(2))
	if err != nil {
		t.Fatalf("List with limit failed: %v", err)
	}
	if len(keys) > 2 {
		t.Fatalf("expected ≤2 keys with limit, got %d", len(keys))
	}

	// Offset beyond list length
	keys, err = s.List(ctx, store.ListOffset(100))
	if err != nil {
		t.Fatalf("List with large offset failed: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("expected 0 keys with offset beyond length, got %d", len(keys))
	}

	// Limit + offset within range
	keys, err = s.List(ctx, store.ListLimit(2), store.ListOffset(1))
	if err != nil {
		t.Fatalf("List with limit+offset failed: %v", err)
	}
	if len(keys) > 2 {
		t.Fatalf("expected ≤2 keys, got %d", len(keys))
	}
}

func TestStoreListError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.SetError(store.ErrNotConnected)
	_, err := s.List(ctx)
	if err == nil {
		t.Fatal("expected error from List when global error is set")
	}
}

func TestExpectedListError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectList().WillReturnError(store.ErrNotFound)
	_, err := s.List(ctx)
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestExpectedListTimes(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectList().Times(2).WillReturn("k1", "k2")
	for range 2 {
		keys, err := s.List(ctx)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(keys) != 2 {
			t.Fatalf("expected 2 keys, got %d", len(keys))
		}
	}
	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestExpectedDeleteTimes(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectDelete("dk").Times(1)
	if err := s.Delete(ctx, "dk"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestExpectedDeleteError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectDelete("ek").WillReturnError(store.ErrNotFound)
	if err := s.Delete(ctx, "ek"); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteGlobalError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.SetError(store.ErrNotConnected)
	if err := s.Delete(ctx, "k"); err == nil {
		t.Fatal("expected error from Delete when global error is set")
	}
}

func TestExpectedExistsTimes(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	_ = s.Write(ctx, "ek", "v")
	s.ExpectExists("ek").Times(1)
	if err := s.Exists(ctx, "ek"); err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestExpectedReadTimes(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	_ = s.Write(ctx, "rk", "rv")
	s.ExpectRead("rk").Times(1)
	var v string
	if err := s.Read(ctx, "rk", &v); err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if err := s.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations not met: %v", err)
	}
}

func TestExpectedReadError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectRead("missing").WillReturnError(store.ErrNotFound)
	var v string
	if err := s.Read(ctx, "missing", &v); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestReadGlobalError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.SetError(store.ErrNotConnected)
	var v string
	if err := s.Read(ctx, "k", &v); err == nil {
		t.Fatal("expected error from Read when global error is set")
	}
}

func TestExpectedWriteNamespace(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectWrite("nk").WithNamespace("ns1").WithValue("nv")
	if err := s.Write(ctx, "nk", "nv", store.WriteNamespace("ns1")); err != nil {
		t.Fatalf("Write with namespace failed: %v", err)
	}
}

func TestExpectedWriteError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.ExpectWrite("ek").WillReturnError(store.ErrNotFound)
	if err := s.Write(ctx, "ek", "v"); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestWriteGlobalError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.SetError(store.ErrNotConnected)
	if err := s.Write(ctx, "k", "v"); err == nil {
		t.Fatal("expected error from Write when global error is set")
	}
}

func TestExpectationsNotMetWrite(t *testing.T) {
	s := NewStore()
	s.ExpectWrite("k").Times(3)
	// Only call once so expectation is not met
	_ = s.Write(context.Background(), "k", "v")
	if err := s.ExpectationsWereMet(); err == nil {
		t.Fatal("expected ExpectationsWereMet to return error")
	}
}

func TestExpectationsNotMetRead(t *testing.T) {
	s := NewStore()
	s.ExpectRead("k").Times(2)
	var v string
	_ = s.Read(context.Background(), "k", &v)
	if err := s.ExpectationsWereMet(); err == nil {
		t.Fatal("expected ExpectationsWereMet to return error for read")
	}
}

func TestExpectationsNotMetDelete(t *testing.T) {
	s := NewStore()
	s.ExpectDelete("k").Times(2)
	_ = s.Delete(context.Background(), "k")
	if err := s.ExpectationsWereMet(); err == nil {
		t.Fatal("expected ExpectationsWereMet to return error for delete")
	}
}

func TestExpectationsNotMetExists(t *testing.T) {
	s := NewStore()
	s.ExpectExists("k").Times(2)
	_ = s.Exists(context.Background(), "k")
	if err := s.ExpectationsWereMet(); err == nil {
		t.Fatal("expected ExpectationsWereMet to return error for exists")
	}
}

func TestExpectationsNotMetList(t *testing.T) {
	s := NewStore()
	s.ExpectList().Times(3)
	_, _ = s.List(context.Background())
	if err := s.ExpectationsWereMet(); err == nil {
		t.Fatal("expected ExpectationsWereMet to return error for list")
	}
}

func TestExistsGlobalError(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	s.SetError(store.ErrNotConnected)
	if err := s.Exists(ctx, "k"); err == nil {
		t.Fatal("expected error from Exists when global error is set")
	}
}

func TestWatcherSendEventFull(t *testing.T) {
	w := NewWatcher()
	// Fill the channel (capacity 1)
	ev1 := &mockEvent{key: "k1"}
	ev2 := &mockEvent{key: "k2"}
	w.SendEvent(ev1)
	// This second send should be dropped silently (channel full)
	w.SendEvent(ev2)
	got, err := w.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != ev1 {
		t.Fatal("expected first event")
	}
	// Stop the watcher so Next doesn't block
	w.Stop()
}
