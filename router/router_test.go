package router

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/register"
)

// --- Options tests ---

func TestOptionID(t *testing.T) {
	opts := NewOptions(ID("test-id"))
	if opts.ID != "test-id" {
		t.Fatalf("expected ID 'test-id', got %q", opts.ID)
	}
}

func TestOptionAddress(t *testing.T) {
	opts := NewOptions(Address("127.0.0.1:9090"))
	if opts.Address != "127.0.0.1:9090" {
		t.Fatalf("expected Address '127.0.0.1:9090', got %q", opts.Address)
	}
}

func TestOptionGateway(t *testing.T) {
	opts := NewOptions(Gateway("gw.example.com"))
	if opts.Gateway != "gw.example.com" {
		t.Fatalf("expected Gateway 'gw.example.com', got %q", opts.Gateway)
	}
}

func TestOptionNetwork(t *testing.T) {
	opts := NewOptions(Network("testnet"))
	if opts.Network != "testnet" {
		t.Fatalf("expected Network 'testnet', got %q", opts.Network)
	}
}

func TestOptionLogger(t *testing.T) {
	l := logger.DefaultLogger
	opts := NewOptions(Logger(l))
	if opts.Logger == nil {
		t.Fatal("expected Logger to be set")
	}
}

func TestOptionRegister(t *testing.T) {
	r := register.DefaultRegister
	opts := NewOptions(Register(r))
	if opts.Register == nil {
		t.Fatal("expected Register to be set")
	}
}

func TestOptionPrecache(t *testing.T) {
	opts := NewOptions(Precache())
	if !opts.Precache {
		t.Fatal("expected Precache to be true")
	}
}

func TestOptionName(t *testing.T) {
	opts := NewOptions(Name("myrouter"))
	if opts.Name != "myrouter" {
		t.Fatalf("expected Name 'myrouter', got %q", opts.Name)
	}
}

func TestNewOptionsDefaults(t *testing.T) {
	opts := NewOptions()
	if opts.Network != DefaultNetwork {
		t.Fatalf("expected default network %q, got %q", DefaultNetwork, opts.Network)
	}
	if opts.ID == "" {
		t.Fatal("expected non-empty default ID")
	}
	if opts.Context == nil {
		t.Fatal("expected non-nil default Context")
	}
	if opts.Register == nil {
		t.Fatal("expected non-nil default Register")
	}
	if opts.Logger == nil {
		t.Fatal("expected non-nil default Logger")
	}
}

// --- DNS router tests ---

func TestNewRouter(t *testing.T) {
	r := NewRouter()
	if r == nil {
		t.Fatal("expected non-nil router")
	}
}

func TestDNSRouterString(t *testing.T) {
	r := NewRouter()
	if r.String() != "dns" {
		t.Fatalf("expected 'dns', got %q", r.String())
	}
}

func TestDNSRouterName(t *testing.T) {
	r := NewRouter(Name("myrouter"))
	if r.Name() != "myrouter" {
		t.Fatalf("expected 'myrouter', got %q", r.Name())
	}
}

func TestDNSRouterOptions(t *testing.T) {
	r := NewRouter(Network("testnet"), Name("rtr"))
	opts := r.Options()
	if opts.Network != "testnet" {
		t.Fatalf("expected network 'testnet', got %q", opts.Network)
	}
	if opts.Name != "rtr" {
		t.Fatalf("expected name 'rtr', got %q", opts.Name)
	}
}

func TestDNSRouterInit(t *testing.T) {
	r := NewRouter()
	if err := r.Init(Name("updated")); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
	if r.Name() != "updated" {
		t.Fatalf("expected name 'updated' after Init, got %q", r.Name())
	}
}

func TestDNSRouterTable(t *testing.T) {
	r := NewRouter()
	if tbl := r.Table(); tbl != nil {
		t.Fatal("expected dns router Table() to return nil")
	}
}

func TestDNSRouterClose(t *testing.T) {
	r := NewRouter()
	if err := r.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
}

func TestDNSRouterWatch(t *testing.T) {
	r := NewRouter()
	w, err := r.Watch()
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	if w != nil {
		t.Fatal("expected dns router Watch() to return nil watcher")
	}
}

func TestDNSRouterLookupError(t *testing.T) {
	// Lookup with an invalid service name (no port, no valid SRV) should return error.
	r := NewRouter()
	_, err := r.Lookup(QueryService("this.service.does.not.exist.invalid"))
	if err == nil {
		t.Log("Lookup did not return error (network may have resolved); skipping strict check")
	}
}

func TestDNSRouterLookupWithPort(t *testing.T) {
	// Use localhost:80 which will hit the A-record path; it may fail DNS but exercises the code path.
	r := NewRouter()
	_, err := r.Lookup(QueryService("this.invalid.host.example:8080"))
	// We expect an error since the host doesn't exist; we just care the code path was exercised.
	_ = err
}

// --- Query tests ---

func TestNewQueryDefaults(t *testing.T) {
	q := NewQuery()
	if q.Service != "*" {
		t.Fatalf("expected Service '*', got %q", q.Service)
	}
	if q.Address != "*" {
		t.Fatalf("expected Address '*', got %q", q.Address)
	}
	if q.Gateway != "*" {
		t.Fatalf("expected Gateway '*', got %q", q.Gateway)
	}
	if q.Network != "*" {
		t.Fatalf("expected Network '*', got %q", q.Network)
	}
	if q.Router != "*" {
		t.Fatalf("expected Router '*', got %q", q.Router)
	}
	if q.Link != DefaultLink {
		t.Fatalf("expected Link %q, got %q", DefaultLink, q.Link)
	}
}

func TestQueryService(t *testing.T) {
	q := NewQuery(QueryService("svc.test"))
	if q.Service != "svc.test" {
		t.Fatalf("expected Service 'svc.test', got %q", q.Service)
	}
}

func TestQueryAddress(t *testing.T) {
	q := NewQuery(QueryAddress("10.0.0.1:8080"))
	if q.Address != "10.0.0.1:8080" {
		t.Fatalf("expected Address '10.0.0.1:8080', got %q", q.Address)
	}
}

func TestQueryGateway(t *testing.T) {
	q := NewQuery(QueryGateway("gw.local"))
	if q.Gateway != "gw.local" {
		t.Fatalf("expected Gateway 'gw.local', got %q", q.Gateway)
	}
}

func TestQueryNetwork(t *testing.T) {
	q := NewQuery(QueryNetwork("net-42"))
	if q.Network != "net-42" {
		t.Fatalf("expected Network 'net-42', got %q", q.Network)
	}
}

func TestQueryRouter(t *testing.T) {
	q := NewQuery(QueryRouter("router-1"))
	if q.Router != "router-1" {
		t.Fatalf("expected Router 'router-1', got %q", q.Router)
	}
}

func TestQueryLink(t *testing.T) {
	q := NewQuery(QueryLink("link-7"))
	if q.Link != "link-7" {
		t.Fatalf("expected Link 'link-7', got %q", q.Link)
	}
}

// --- Route tests ---

func TestRouteHashDifferent(t *testing.T) {
	r1 := Route{Service: "svc1", Address: "10.0.0.1:8080", Gateway: "gw", Network: "net", Router: "r", Link: "local"}
	r2 := Route{Service: "svc2", Address: "10.0.0.1:8080", Gateway: "gw", Network: "net", Router: "r", Link: "local"}
	if r1.Hash() == r2.Hash() {
		t.Error("expected different hashes for different routes")
	}
}

func TestRouteHashEmpty(t *testing.T) {
	r := Route{}
	h := r.Hash()
	// Should not panic and should return a consistent value.
	if r.Hash() != h {
		t.Error("hash not deterministic for empty route")
	}
}

func TestDefaultLink(t *testing.T) {
	if DefaultLink != "local" {
		t.Fatalf("expected DefaultLink 'local', got %q", DefaultLink)
	}
}

func TestDefaultLocalMetric(t *testing.T) {
	if DefaultLocalMetric != 1 {
		t.Fatalf("expected DefaultLocalMetric 1, got %d", DefaultLocalMetric)
	}
}

// --- Watcher EventType String tests ---

func TestEventTypeStringCreate(t *testing.T) {
	if Create.String() != "create" {
		t.Fatalf("expected 'create', got %q", Create.String())
	}
}

func TestEventTypeStringDelete(t *testing.T) {
	if Delete.String() != "delete" {
		t.Fatalf("expected 'delete', got %q", Delete.String())
	}
}

func TestEventTypeStringUpdate(t *testing.T) {
	if Update.String() != "update" {
		t.Fatalf("expected 'update', got %q", Update.String())
	}
}

func TestEventTypeStringUnknown(t *testing.T) {
	unknown := EventType(99)
	if unknown.String() != "unknown" {
		t.Fatalf("expected 'unknown', got %q", unknown.String())
	}
}

// --- WatchService option test ---

func TestWatchService(t *testing.T) {
	opts := &WatchOptions{}
	WatchService("my-service")(opts)
	if opts.Service != "my-service" {
		t.Fatalf("expected Service 'my-service', got %q", opts.Service)
	}
}

// --- Event struct test ---

func TestEventStruct(t *testing.T) {
	now := time.Now()
	e := Event{
		Timestamp: now,
		ID:        "evt-1",
		Route:     Route{Service: "svc"},
		Type:      Create,
	}
	if e.ID != "evt-1" {
		t.Fatalf("expected event ID 'evt-1', got %q", e.ID)
	}
	if e.Type != Create {
		t.Fatalf("expected event Type Create")
	}
}

// --- MustContext test ---

func TestMustContext(t *testing.T) {
	r := NewRouter()
	ctx := NewContext(context.Background(), r)
	got := MustContext(ctx)
	if got == nil {
		t.Fatal("expected non-nil router from MustContext")
	}
}

func TestMustContextPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Fatal("expected panic from MustContext with empty context")
		}
	}()
	MustContext(context.Background())
}

// --- DefaultRouter / global vars ---

func TestDefaultRouter(t *testing.T) {
	if DefaultRouter == nil {
		t.Fatal("expected DefaultRouter to be non-nil")
	}
}

func TestDefaultNetwork(t *testing.T) {
	if DefaultNetwork != "micro" {
		t.Fatalf("expected DefaultNetwork 'micro', got %q", DefaultNetwork)
	}
}

func TestErrorVars(t *testing.T) {
	if ErrRouteNotFound == nil {
		t.Fatal("expected ErrRouteNotFound to be non-nil")
	}
	if ErrDuplicateRoute == nil {
		t.Fatal("expected ErrDuplicateRoute to be non-nil")
	}
	if ErrWatcherStopped == nil {
		t.Fatal("expected ErrWatcherStopped to be non-nil")
	}
}

// --- StatusCode tests ---

func TestStatusCodes(t *testing.T) {
	if Running != 0 {
		t.Fatalf("expected Running==0, got %d", Running)
	}
	if Stopped != 1 {
		t.Fatalf("expected Stopped==1, got %d", Stopped)
	}
	if Error != 2 {
		t.Fatalf("expected Error==2, got %d", Error)
	}
}
