package mock

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v4/store"
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
	var value interface{}
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
}

func TestStoreConnectDisconnect(t *testing.T) {
	s := NewStore()

	err := s.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	err = s.Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}

	// Test error propagation
	s.ExpectWrite("test_key").WillReturnError(store.ErrNotConnected)
	err = s.Write(context.Background(), "test_key", "value")
	if err != store.ErrNotConnected {
		t.Fatalf("Expected store.ErrNotConnected, got %v", err)
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
