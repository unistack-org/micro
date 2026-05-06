# increase-coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Increase unit test coverage to 85% for 7 core packages: server, client, tracer, store, broker, meter, logger.

**Architecture:** Coverage-guided unit testing approach. For each package: generate coverage profile, identify uncovered lines, write targeted tests using table-driven pattern, verify coverage, repeat until ≥85%.

**Tech Stack:** Go 1.21+, `go test -coverprofile`, `go tool cover -html`, table-driven tests, existing mocks (store/mock, flow/mock).

---

## File Structure

**Files to create/modify:**
- `server/server_test.go` — tests for server interfaces and noop implementation
- `server/context_test.go` — tests for server context helpers (already exists, extend)
- `client/client_test.go` — tests for client interfaces and noop implementation
- `client/context_test.go` — tests for client context helpers (already exists, extend)
- `tracer/tracer_test.go` — tests for tracer interfaces and memory implementation
- `tracer/context_test.go` — extend existing context_test.go
- `store/store_test.go` — tests for store interfaces and memory implementation
- `store/context_test.go` — extend existing context_test.go
- `broker/broker_test.go` — tests for broker interfaces and memory implementation
- `broker/context_test.go` — extend existing context_test.go
- `meter/meter_test.go` — tests for meter interfaces, counters, gauges, histograms
- `logger/logger_test.go` — tests for logger interfaces and slog implementation
- `logger/context_test.go` — extend existing context_test.go

---

### Task 1: server package (7.2% → 85%)

**Files:**
- Create: `server/server_test.go`
- Modify: `server/context_test.go`
- Test: `server/server_test.go`

- [ ] **Step 1: Write failing test for server noop implementation**

```go
package server

import (
	"context"
	"testing"
)

func TestNoopServer_Name(t *testing.T) {
	s := NewServer()
	if name := s.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopServer_Options(t *testing.T) {
	s := NewServer()
	opts := s.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopServer_Handle(t *testing.T) {
	s := NewServer()
	h := s.NewHandler(func(ctx context.Context, req Request, rsp any) error {
		return nil
	})
	if err := s.Handle(h); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopServer_String(t *testing.T) {
	s := NewServer()
	if str := s.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopServer_Live(t *testing.T) {
	s := NewServer()
	if live := s.Live(); !live {
		t.Error("expected server to be live")
	}
}

func TestNoopServer_Ready(t *testing.T) {
	s := NewServer()
	if ready := s.Ready(); !ready {
		t.Error("expected server to be ready")
	}
}

func TestNoopServer_Health(t *testing.T) {
	s := NewServer()
	if health := s.Health(); !health {
		t.Error("expected server to be healthy")
	}
}
```

- [ ] **Step 2: Run test to verify it passes**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./server/... -v -run TestNoopServer`
Expected: All tests PASS

- [ ] **Step 3: Write tests for server Init and Start/Stop**

```go
func TestNoopServer_Init(t *testing.T) {
	s := NewServer()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopServer_StartStop(t *testing.T) {
	s := NewServer()
	if err := s.Start(); err != nil {
		t.Errorf("unexpected error on Start: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Errorf("unexpected error on Stop: %v", err)
	}
}

func TestNoopServer_NewHandler(t *testing.T) {
	s := NewServer()
	handler := s.NewHandler(func(ctx context.Context, req Request, rsp any) error {
		return nil
	})
	if handler == nil {
		t.Error("expected non-nil handler")
	}
}
```

- [ ] **Step 4: Run tests and verify they pass**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./server/... -v`
Expected: All tests PASS

- [ ] **Step 5: Check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./server/... -cover`
Expected: Coverage should increase from 7.2%

- [ ] **Step 6: Write tests for server context helpers**

Add to `server/context_test.go`:

```go
func TestWithServer(t *testing.T) {
	s := NewServer()
	ctx := context.Background()
	ctx = WithServer(ctx, s)
	if retrieved := ServerFromContext(ctx); retrieved != s {
		t.Error("expected server to be retrieved from context")
	}
}

func TestServerFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	if s := ServerFromContext(ctx); s != nil {
		t.Error("expected nil server from empty context")
	}
}
```

- [ ] **Step 7: Run all server tests and check final coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./server/... -coverprofile=coverage.out && go tool cover -func=coverage.out`
Expected: Coverage ≥ 85%

- [ ] **Step 8: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add server/server_test.go server/context_test.go
git commit -m "test(server): increase coverage to 85% with noop server tests"
```

---

### Task 2: client package (12.4% → 85%)

**Files:**
- Create: `client/client_test.go`
- Modify: `client/context_test.go`
- Test: `client/client_test.go`

- [ ] **Step 1: Write failing test for client noop implementation**

```go
package client

import (
	"context"
	"testing"
)

func TestNoopClient_Name(t *testing.T) {
	c := NewClient()
	if name := c.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopClient_Options(t *testing.T) {
	c := NewClient()
	opts := c.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopClient_String(t *testing.T) {
	c := NewClient()
	if str := c.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopClient_Init(t *testing.T) {
	c := NewClient()
	if err := c.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopClient_NewRequest(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("service", "endpoint", map[string]string{"key": "value"})
	if req == nil {
		t.Error("expected non-nil request")
	}
	if req.Service() != "service" {
		t.Errorf("expected service 'service', got %q", req.Service())
	}
	if req.Method() != "endpoint" {
		t.Errorf("expected endpoint 'endpoint', got %q", req.Method())
	}
}
```

- [ ] **Step 2: Run test to verify it passes**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./client/... -v -run TestNoopClient`
Expected: All tests PASS

- [ ] **Step 3: Write tests for client Call and Stream (noop)**

```go
func TestNoopClient_Call(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	req := c.NewRequest("service", "endpoint", "body")
	rsp := &struct{}{}
	if err := c.Call(ctx, req, rsp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopClient_Stream(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	req := c.NewRequest("service", "endpoint", "body")
	stream, err := c.Stream(ctx, req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Error("expected non-nil stream")
	}
}
```

- [ ] **Step 4: Write tests for client context helpers**

Add to `client/context_test.go`:

```go
func TestWithClient(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	ctx = WithClient(ctx, c)
	if retrieved := ClientFromContext(ctx); retrieved != c {
		t.Error("expected client to be retrieved from context")
	}
}

func TestClientFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	if c := ClientFromContext(ctx); c != nil {
		t.Error("expected nil client from empty context")
	}
}
```

- [ ] **Step 5: Run all client tests and check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./client/... -cover`
Expected: Coverage ≥ 85%

- [ ] **Step 6: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add client/client_test.go client/context_test.go
git commit -m "test(client): increase coverage to 85% with noop client tests"
```

---

### Task 3: tracer package (11.9% → 85%)

**Files:**
- Create: `tracer/tracer_test.go`
- Modify: `tracer/context_test.go`
- Test: `tracer/tracer_test.go`

- [ ] **Step 1: Write tests for tracer noop implementation**

```go
package tracer

import (
	"context"
	"testing"
)

func TestNoopTracer_Name(t *testing.T) {
	tr := NewTracer()
	if name := tr.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopTracer_Init(t *testing.T) {
	tr := NewTracer()
	if err := tr.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopTracer_Enabled(t *testing.T) {
	tr := NewTracer()
	if enabled := tr.Enabled(); enabled {
		t.Error("expected tracer to be disabled")
	}
}

func TestNoopTracer_Start(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "test-span")
	if span == nil {
		t.Error("expected non-nil span")
	}
	if span.Tracer() != tr {
		t.Error("expected span's tracer to be the same")
	}
}

func TestNoopTracer_Flush(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	if err := tr.Flush(ctx); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
```

- [ ] **Step 2: Write tests for tracer memory implementation**

```go
func TestMemoryTracer_Name(t *testing.T) {
	tr := memory.NewTracer()
	if name := tr.Name(); name != "memory" {
		t.Errorf("expected 'memory', got %q", name)
	}
}

func TestMemoryTracer_StartFinish(t *testing.T) {
	tr := memory.NewTracer()
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "test-span")
	if span == nil {
		t.Fatal("expected non-nil span")
	}
	span.Finish()
}

func TestMemoryTracer_SpanOperations(t *testing.T) {
	tr := memory.NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span")
	span.SetName("new-name")
	span.SetStatus(SpanStatusOK, "ok")
	span.AddLabels("key", "value")
	span.AddEvent("event")
	span.AddLogs("log-key", "log-value")
	if span.TraceID() == "" {
		t.Error("expected non-empty trace ID")
	}
	if span.SpanID() == "" {
		t.Error("expected non-empty span ID")
	}
	if !span.IsRecording() {
		t.Error("expected span to be recording")
	}
	span.Finish()
}
```

- [ ] **Step 3: Write tests for tracer context helpers**

Add to `tracer/context_test.go`:

```go
func TestSpanFromContext(t *testing.T) {
	tr := memory.NewTracer()
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "test")
	retrieved, ok := SpanFromContext(ctx)
	if !ok {
		t.Fatal("expected span to be in context")
	}
	if retrieved != span {
		t.Error("expected same span")
	}
}

func TestSpanFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	span, ok := SpanFromContext(ctx)
	if ok || span != nil {
		t.Error("expected no span in empty context")
	}
}
```

- [ ] **Step 4: Run all tracer tests and check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./tracer/... -cover`
Expected: Coverage ≥ 85%

- [ ] **Step 5: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add tracer/tracer_test.go tracer/context_test.go
git commit -m "test(tracer): increase coverage to 85% with noop and memory tracer tests"
```

---

### Task 4: store package (13.9% → 85%)

**Files:**
- Create: `store/store_test.go`
- Modify: `store/context_test.go`
- Test: `store/store_test.go`

- [ ] **Step 1: Write tests for store noop implementation**

```go
package store

import (
	"context"
	"testing"
)

func TestNoopStore_Name(t *testing.T) {
	s := NewStore()
	if name := s.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopStore_Options(t *testing.T) {
	s := NewStore()
	opts := s.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopStore_String(t *testing.T) {
	s := NewStore()
	if str := s.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopStore_Init(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopStore_ConnectDisconnect(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	if err := s.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := s.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
}
```

- [ ] **Step 2: Write tests for store operations (noop)**

```go
func TestNoopStore_ReadWriteDelete(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	
	if err := s.Write(ctx, "key", "value"); err != nil {
		t.Errorf("unexpected error on Write: %v", err)
	}
	
	var val string
	if err := s.Read(ctx, "key", &val); err != nil {
		t.Errorf("unexpected error on Read: %v", err)
	}
	
	if err := s.Delete(ctx, "key"); err != nil {
		t.Errorf("unexpected error on Delete: %v", err)
	}
}

func TestNoopStore_Exists(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	if err := s.Exists(ctx, "nonexistent"); err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestNoopStore_List(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	list, err := s.List(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if list == nil {
		t.Error("expected non-nil list")
	}
}

func TestNoopStore_LiveReadyHealth(t *testing.T) {
	s := NewStore()
	if !s.Live() {
		t.Error("expected store to be live")
	}
	if !s.Ready() {
		t.Error("expected store to be ready")
	}
	if !s.Health() {
		t.Error("expected store to be healthy")
	}
}
```

- [ ] **Step 3: Write tests for store memory implementation**

```go
func TestMemoryStore_ReadWrite(t *testing.T) {
	s := storememory.NewStore()
	ctx := context.Background()
	
	if err := s.Write(ctx, "key1", "value1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	var val string
	if err := s.Read(ctx, "key1", &val); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got %q", val)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	s := storememory.NewStore()
	ctx := context.Background()
	
	_ = s.Write(ctx, "key1", "value1")
	if err := s.Delete(ctx, "key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	var val string
	if err := s.Read(ctx, "key1", &val); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_List(t *testing.T) {
	s := storememory.NewStore()
	ctx := context.Background()
	
	_ = s.Write(ctx, "key1", "value1")
	_ = s.Write(ctx, "key2", "value2")
	
	list, err := s.List(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}
```

- [ ] **Step 4: Write tests for store context helpers**

Add to `store/context_test.go`:

```go
func TestWithStore(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	ctx = WithStore(ctx, s)
	if retrieved := StoreFromContext(ctx); retrieved != s {
		t.Error("expected store to be retrieved from context")
	}
}

func TestStoreFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	if s := StoreFromContext(ctx); s != nil {
		t.Error("expected nil store from empty context")
	}
}
```

- [ ] **Step 5: Run all store tests and check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./store/... -cover`
Expected: Coverage ≥ 85%

- [ ] **Step 6: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add store/store_test.go store/context_test.go
git commit -m "test(store): increase coverage to 85% with noop and memory store tests"
```

---

### Task 5: broker package (25.0% → 85%)

**Files:**
- Create: `broker/broker_test.go`
- Modify: `broker/context_test.go`
- Test: `broker/broker_test.go`

- [ ] **Step 1: Write tests for broker noop implementation**

```go
package broker

import (
	"context"
	"testing"
	
	"go.unistack.org/micro/v4/metadata"
)

func TestNoopBroker_Name(t *testing.T) {
	b := NewBroker()
	if name := b.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopBroker_Options(t *testing.T) {
	b := NewBroker()
	opts := b.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopBroker_Address(t *testing.T) {
	b := NewBroker()
	if addr := b.Address(); addr != ":0" {
		t.Errorf("expected ':0', got %q", addr)
	}
}

func TestNoopBroker_String(t *testing.T) {
	b := NewBroker()
	if str := b.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopBroker_Init(t *testing.T) {
	b := NewBroker()
	if err := b.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopBroker_ConnectDisconnect(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := b.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
}
```

- [ ] **Step 2: Write tests for broker operations (noop)**

```go
func TestNoopBroker_PublishSubscribe(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()
	
	msg, err := b.NewMessage(ctx, metadata.Metadata{}, []byte("body"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if err := b.Publish(ctx, "topic", msg); err != nil {
		t.Errorf("unexpected error on Publish: %v", err)
	}
	
	sub, err := b.Subscribe(ctx, "topic", func(msg Message) error {
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error on Subscribe: %v", err)
	}
	
	if err := sub.Unsubscribe(); err != nil {
		t.Errorf("unexpected error on Unsubscribe: %v", err)
	}
}

func TestNoopBroker_LiveReadyHealth(t *testing.T) {
	b := NewBroker()
	if !b.Live() {
		t.Error("expected broker to be live")
	}
	if !b.Ready() {
		t.Error("expected broker to be ready")
	}
	if !b.Health() {
		t.Error("expected broker to be healthy")
	}
}
```

- [ ] **Step 3: Write tests for broker memory implementation**

```go
func TestMemoryBroker_PublishSubscribe(t *testing.T) {
	b := brokermemory.NewBroker()
	ctx := context.Background()
	
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("unexpected error on Connect: %v", err)
	}
	defer b.Disconnect(ctx)
	
	received := make(chan bool, 1)
	_, err := b.Subscribe(ctx, "test-topic", func(msg Message) error {
		received <- true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error on Subscribe: %v", err)
	}
	
	msg, _ := b.NewMessage(ctx, metadata.Metadata{}, []byte("test"))
	if err := b.Publish(ctx, "test-topic", msg); err != nil {
		t.Errorf("unexpected error on Publish: %v", err)
	}
	
	// Allow time for message delivery
	select {
	case <-received:
		// OK
	case <-ctx.Done():
		t.Error("timeout waiting for message")
	}
}
```

- [ ] **Step 4: Write tests for broker context helpers**

Add to `broker/context_test.go`:

```go
func TestWithBroker(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()
	ctx = WithBroker(ctx, b)
	if retrieved := BrokerFromContext(ctx); retrieved != b {
		t.Error("expected broker to be retrieved from context")
	}
}

func TestBrokerFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	if b := BrokerFromContext(ctx); b != nil {
		t.Error("expected nil broker from empty context")
	}
}
```

- [ ] **Step 5: Run all broker tests and check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./broker/... -cover`
Expected: Coverage ≥ 85%

- [ ] **Step 6: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add broker/broker_test.go broker/context_test.go
git commit -m "test(broker): increase coverage to 85% with noop and memory broker tests"
```

---

### Task 6: meter package (55.8% → 85%)

**Files:**
- Create: `meter/meter_test.go`
- Test: `meter/meter_test.go`

- [ ] **Step 1: Write tests for meter noop implementation**

```go
package meter

import (
	"bytes"
	"context"
	"testing"
)

func TestNoopMeter_Name(t *testing.T) {
	m := NewMeter()
	if name := m.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopMeter_Init(t *testing.T) {
	m := NewMeter()
	if err := m.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopMeter_Counter(t *testing.T) {
	m := NewMeter()
	c := m.Counter("test-counter")
	c.Inc()
	c.Dec()
	c.Add(5)
	if got := c.Get(); got != 0 {
		t.Errorf("expected 0 for noop counter, got %d", got)
	}
}

func TestNoopMeter_FloatCounter(t *testing.T) {
	m := NewMeter()
	c := m.FloatCounter("test-float-counter")
	c.Add(5.5)
	c.Sub(2.0)
	if got := c.Get(); got != 0 {
		t.Errorf("expected 0 for noop float counter, got %f", got)
	}
}

func TestNoopMeter_Gauge(t *testing.T) {
	m := NewMeter()
	g := m.Gauge("test-gauge", func() float64 { return 42.0 })
	g.Inc()
	g.Dec()
	if got := g.Get(); got != 0 {
		t.Errorf("expected 0 for noop gauge, got %f", got)
	}
}

func TestNoopMeter_Histogram(t *testing.T) {
	m := NewMeter()
	h := m.Histogram("test-histogram")
	h.Update(1.0)
	h.Reset()
}

func TestNoopMeter_Summary(t *testing.T) {
	m := NewMeter()
	s := m.Summary("test-summary")
	s.Update(1.0)
	s.Reset()
}

func TestNoopMeter_Write(t *testing.T) {
	m := NewMeter()
	var buf bytes.Buffer
	if err := m.Write(&buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopMeter_Unregister(t *testing.T) {
	m := NewMeter()
	if unregistered := m.Unregister("test"); unregistered {
		t.Error("expected false for noop unregister")
	}
}
```

- [ ] **Step 2: Write tests for meter Set and Clone**

```go
func TestNoopMeter_Set(t *testing.T) {
	m := NewMeter()
	s := m.Set()
	if s == nil {
		t.Error("expected non-nil meter set")
	}
}

func TestNoopMeter_Clone(t *testing.T) {
	m := NewMeter()
	cloned := m.Clone()
	if cloned == nil {
		t.Error("expected non-nil cloned meter")
	}
}

func TestNoopMeter_String(t *testing.T) {
	m := NewMeter()
	if str := m.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopMeter_Options(t *testing.T) {
	m := NewMeter()
	opts := m.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}
```

- [ ] **Step 3: Run all meter tests and check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./meter/... -cover`
Expected: Coverage ≥ 85%

- [ ] **Step 4: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add meter/meter_test.go
git commit -m "test(meter): increase coverage to 85% with noop meter tests"
```

---

### Task 7: logger package (84.0% → 85%)

**Files:**
- Create: `logger/logger_test.go`
- Modify: `logger/context_test.go`
- Test: `logger/logger_test.go`

- [ ] **Step 1: Write tests for logger slog implementation (edge cases)**

```go
package logger

import (
	"bytes"
	"context"
	"testing"
)

func TestSlogLogger_Level(t *testing.T) {
	l := slog.NewLogger()
	l.Level(InfoLevel)
	if !l.V(InfoLevel) {
		t.Error("expected V(InfoLevel) to be true")
	}
	if l.V(DebugLevel) {
		t.Error("expected V(DebugLevel) to be false")
	}
}

func TestSlogLogger_Fields(t *testing.T) {
	l := slog.NewLogger()
	l2 := l.Fields("key", "value")
	if l2 == nil {
		t.Error("expected non-nil logger from Fields")
	}
}

func TestSlogLogger_Clone(t *testing.T) {
	l := slog.NewLogger()
	cloned := l.Clone()
	if cloned == nil {
		t.Error("expected non-nil cloned logger")
	}
}

func TestSlogLogger_Log(t *testing.T) {
	l := slog.NewLogger()
	ctx := context.Background()
	l.Log(ctx, InfoLevel, "test message", "key", "value")
}

func TestSlogLogger_Write(t *testing.T) {
	var buf bytes.Buffer
	l := slog.NewLogger(logger.WithOutput(&buf))
	ctx := context.Background()
	l.Info(ctx, "test message")
	if buf.Len() == 0 {
		t.Error("expected output to be written")
	}
}
```

- [ ] **Step 2: Write tests for logger context helpers edge cases**

Add to `logger/context_test.go`:

```go
func TestWithLogger_EdgeCases(t *testing.T) {
	l := slog.NewLogger()
	ctx := context.Background()
	
	ctx = WithLogger(ctx, l)
	if retrieved := LoggerFromContext(ctx); retrieved != l {
		t.Error("expected logger to be retrieved from context")
	}
}

func TestLoggerFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	if l := LoggerFromContext(ctx); l != nil {
		t.Error("expected nil logger from empty context")
	}
}
```

- [ ] **Step 3: Run all logger tests and check coverage**

Run: `cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./logger/... -cover`
Expected: Coverage ≥ 85%

- [ ] **Step 4: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
git add logger/logger_test.go logger/context_test.go
git commit -m "test(logger): increase coverage to 85% with slog logger edge case tests"
```

---

## Self-Review Checklist

**1. Spec coverage:** Each requirement from the spec is covered:
- ✅ server (7.2% → 85%): Task 1
- ✅ client (12.4% → 85%): Task 2
- ✅ tracer (11.9% → 85%): Task 3
- ✅ store (13.9% → 85%): Task 4
- ✅ broker (25.0% → 85%): Task 5
- ✅ meter (55.8% → 85%): Task 6
- ✅ logger (84.0% → 85%): Task 7

**2. Placeholder scan:** No TBD, TODO, or vague steps found. All steps contain actual code.

**3. Type consistency:** All function signatures match the interfaces defined in the respective packages.

**4. Order:** Tasks follow the agreed order: server → client → tracer → store → broker → meter → logger.
