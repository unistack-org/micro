package store

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/tracer"
)

type contextKey string

const testContextKey contextKey = "key"

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
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()
	if err := s.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := s.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
}

func TestNoopStore_ReadWriteDelete(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
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
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()
	if err := s.Exists(ctx, "nonexistent"); err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestNoopStore_List(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()
	list, err := s.List(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if list != nil {
		t.Error("expected nil list")
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

func TestNoopStore_Name(t *testing.T) {
	s := NewStore()
	if s.Name() != "" {
		t.Errorf("expected empty name, got %q", s.Name())
	}
	s2 := NewStore(Name("noop"))
	if s2.Name() != "noop" {
		t.Errorf("expected 'noop', got %q", s2.Name())
	}
}

func TestNoopStore_Watch(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	w, err := s.Watch(ctx)
	if err != nil {
		t.Fatalf("unexpected error on Watch: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
	evt, err := w.Next()
	if err != nil {
		t.Errorf("unexpected error on Next: %v", err)
	}
	if evt != nil {
		t.Errorf("expected nil event, got %v", evt)
	}
	w.Stop()
}

func TestNewNamespaceStore(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ns := NewNamespaceStore(s, "test-ns")
	if ns == nil {
		t.Fatal("expected non-nil NamespaceStore")
	}
	ctx := context.Background()
	// Test Write/Read
	if err := ns.Write(ctx, "key", "val"); err != nil {
		t.Fatal(err)
	}
	if err := ns.Read(ctx, "key", nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Test Delete
	if err := ns.Delete(ctx, "key"); err != nil {
		t.Errorf("unexpected error on Delete: %v", err)
	}
	// Test Exists
	if err := ns.Exists(ctx, "key"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	// Test List
	list, err := ns.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if list != nil {
		t.Error("expected nil list")
	}
	// Test Init
	if err := ns.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	// Test Connect/Disconnect
	if err := ns.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := ns.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
	// Test Options
	opts := ns.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context")
	}
	// Test Name
	if ns.Name() != "" {
		t.Errorf("expected empty name, got %q", ns.Name())
	}
	// Test String
	if ns.String() != "noop" {
		t.Errorf("expected 'noop', got %q", ns.String())
	}
	// Test Watch
	w, err := ns.Watch(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
	w.Stop()
	// Test Live/Ready/Health
	if !ns.Live() {
		t.Error("expected Live() true")
	}
	if !ns.Ready() {
		t.Error("expected Ready() true")
	}
	if !ns.Health() {
		t.Error("expected Health() true")
	}
}

func TestNoopStore_LazyConnect(t *testing.T) {
	s := NewStore(LazyConnect(true))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// Read should trigger connect
	if err := s.Read(ctx, "key", nil); err != nil {
		t.Errorf("unexpected error on Read with LazyConnect: %v", err)
	}
	// Write should trigger connect
	if err := s.Write(ctx, "key", "val"); err != nil {
		t.Errorf("unexpected error on Write with LazyConnect: %v", err)
	}
	// Exists should trigger connect
	if err := s.Exists(ctx, "key"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	// List should trigger connect
	list, err := s.List(ctx)
	if err != nil {
		t.Errorf("unexpected error on List with LazyConnect: %v", err)
	}
	if list != nil {
		t.Error("expected nil list")
	}
	// Delete should trigger connect
	if err := s.Delete(ctx, "key"); err != nil {
		t.Errorf("unexpected error on Delete with LazyConnect: %v", err)
	}
}

func TestStoreWatchFunction(t *testing.T) {
	w, err := Watch(context.Background())
	if w != nil || err != nil {
		t.Errorf("expected nil watcher and nil error, got %v, %v", w, err)
	}
}

func TestNoopWatcherStop(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	w, err := s.Watch(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Stop should not panic
	w.Stop()
}

// Original TestMustContext continues

func TestSetReadOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := SetReadOption(key, val)
	opts := NewReadOptions(opt)
	if opts.Context == nil {
		t.Fatal("expected non-nil context")
	}
	got := opts.Context.Value(key)
	if got != val {
		t.Errorf("expected %q, got %v", val, got)
	}
}

func TestSetWriteOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := SetWriteOption(key, val)
	opts := NewWriteOptions(opt)
	if opts.Context == nil {
		t.Fatal("expected non-nil context")
	}
	got := opts.Context.Value(key)
	if got != val {
		t.Errorf("expected %q, got %v", val, got)
	}
}

func TestSetListOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := SetListOption(key, val)
	opts := NewListOptions(opt)
	if opts.Context == nil {
		t.Fatal("expected non-nil context")
	}
	got := opts.Context.Value(key)
	if got != val {
		t.Errorf("expected %q, got %v", val, got)
	}
}

func TestSetDeleteOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := SetDeleteOption(key, val)
	opts := NewDeleteOptions(opt)
	if opts.Context == nil {
		t.Fatal("expected non-nil context")
	}
	got := opts.Context.Value(key)
	if got != val {
		t.Errorf("expected %q, got %v", val, got)
	}
}

func TestSetExistsOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := SetExistsOption(key, val)
	opts := NewExistsOptions(opt)
	if opts.Context == nil {
		t.Fatal("expected non-nil context")
	}
	got := opts.Context.Value(key)
	if got != val {
		t.Errorf("expected %q, got %v", val, got)
	}
}

func TestTLSConfigOption(t *testing.T) {
	tlsConfig := &tls.Config{}
	s := NewStore(TLSConfig(tlsConfig))
	if s.Options().TLSConfig != tlsConfig {
		t.Error("TLSConfig option not set")
	}
}

func TestContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey, "val")
	s := NewStore(Context(ctx))
	if s.Options().Context.Value(testContextKey) != "val" {
		t.Error("Context option not set")
	}
}

func TestCodecOption(t *testing.T) {
	c := codec.DefaultCodec
	s := NewStore(Codec(c))
	if s.Options().Codec != c {
		t.Error("Codec option not set")
	}
}

func TestLoggerOption(t *testing.T) {
	l := logger.DefaultLogger
	s := NewStore(Logger(l))
	if s.Options().Logger != l {
		t.Error("Logger option not set")
	}
}

func TestMeterOption(t *testing.T) {
	m := meter.DefaultMeter
	s := NewStore(Meter(m))
	if s.Options().Meter != m {
		t.Error("Meter option not set")
	}
}

func TestNameOption(t *testing.T) {
	s := NewStore(Name("test"))
	if s.Options().Name != "test" {
		t.Errorf("expected 'test', got %q", s.Options().Name)
	}
}

func TestSeparatorOption(t *testing.T) {
	s := NewStore(Separator("/"))
	if s.Options().Separator != "/" {
		t.Errorf("expected '/', got %q", s.Options().Separator)
	}
}

func TestNamespaceOption(t *testing.T) {
	s := NewStore(Namespace("ns"))
	if s.Options().Namespace != "ns" {
		t.Errorf("expected 'ns', got %q", s.Options().Namespace)
	}
}

func TestTracerOption(t *testing.T) {
	tr := tracer.DefaultTracer
	s := NewStore(Tracer(tr))
	if s.Options().Tracer != tr {
		t.Error("Tracer option not set")
	}
}

func TestTimeoutOption(t *testing.T) {
	td := time.Second
	s := NewStore(Timeout(td))
	if s.Options().Timeout != td {
		t.Errorf("expected %v, got %v", td, s.Options().Timeout)
	}
}

func TestLazyConnectOption(t *testing.T) {
	s := NewStore(LazyConnect(true))
	if !s.Options().LazyConnect {
		t.Error("LazyConnect option not set")
	}
}

func TestAddrsOption(t *testing.T) {
	s := NewStore(Addrs("addr1", "addr2"))
	opts := s.Options()
	if len(opts.Addrs) != 2 || opts.Addrs[0] != "addr1" || opts.Addrs[1] != "addr2" {
		t.Errorf("unexpected addrs: %v", opts.Addrs)
	}
}

func TestReadTimeoutOption(t *testing.T) {
	td := time.Second
	opts := NewReadOptions(ReadTimeout(td))
	if opts.Timeout != td {
		t.Errorf("expected %v, got %v", td, opts.Timeout)
	}
}

func TestReadNameOption(t *testing.T) {
	opts := NewReadOptions(ReadName("test"))
	if opts.Name != "test" {
		t.Errorf("expected 'test', got %q", opts.Name)
	}
}

func TestReadContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey, "val")
	opts := NewReadOptions(ReadContext(ctx))
	if opts.Context.Value(testContextKey) != "val" {
		t.Error("ReadContext option not set")
	}
}

func TestReadNamespaceOption(t *testing.T) {
	opts := NewReadOptions(ReadNamespace("ns"))
	if opts.Namespace != "ns" {
		t.Errorf("expected 'ns', got %q", opts.Namespace)
	}
}

func TestWriteTimeoutOption(t *testing.T) {
	td := time.Second
	opts := NewWriteOptions(WriteTimeout(td))
	if opts.Timeout != td {
		t.Errorf("expected %v, got %v", td, opts.Timeout)
	}
}

func TestWriteContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey, "val")
	opts := NewWriteOptions(WriteContext(ctx))
	if opts.Context.Value(testContextKey) != "val" {
		t.Error("WriteContext option not set")
	}
}

func TestWriteMetadataOption(t *testing.T) {
	md := metadata.Metadata{"foo": []string{"bar"}}
	opts := NewWriteOptions(WriteMetadata(md))
	if len(opts.Metadata["foo"]) == 0 || opts.Metadata["foo"][0] != "bar" {
		t.Errorf("expected 'bar', got %v", opts.Metadata["foo"])
	}
}

func TestWriteTTLOption(t *testing.T) {
	td := time.Second
	opts := NewWriteOptions(WriteTTL(td))
	if opts.TTL != td {
		t.Errorf("expected %v, got %v", td, opts.TTL)
	}
}

func TestWriteNamespaceOption(t *testing.T) {
	opts := NewWriteOptions(WriteNamespace("ns"))
	if opts.Namespace != "ns" {
		t.Errorf("expected 'ns', got %q", opts.Namespace)
	}
}

func TestWriteNameOption(t *testing.T) {
	opts := NewWriteOptions(WriteName("test"))
	if opts.Name != "test" {
		t.Errorf("expected 'test', got %q", opts.Name)
	}
}

func TestDeleteContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey, "val")
	opts := NewDeleteOptions(DeleteContext(ctx))
	if opts.Context.Value(testContextKey) != "val" {
		t.Error("DeleteContext option not set")
	}
}

func TestDeleteNamespaceOption(t *testing.T) {
	opts := NewDeleteOptions(DeleteNamespace("ns"))
	if opts.Namespace != "ns" {
		t.Errorf("expected 'ns', got %q", opts.Namespace)
	}
}

func TestDeleteNameOption(t *testing.T) {
	opts := NewDeleteOptions(DeleteName("test"))
	if opts.Name != "test" {
		t.Errorf("expected 'test', got %q", opts.Name)
	}
}

func TestDeleteTimeoutOption(t *testing.T) {
	td := time.Second
	opts := NewDeleteOptions(DeleteTimeout(td))
	if opts.Timeout != td {
		t.Errorf("expected %v, got %v", td, opts.Timeout)
	}
}

func TestListContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey, "val")
	opts := NewListOptions(ListContext(ctx))
	if opts.Context.Value(testContextKey) != "val" {
		t.Error("ListContext option not set")
	}
}

func TestListPrefixOption(t *testing.T) {
	opts := NewListOptions(ListPrefix("pre"))
	if opts.Prefix != "pre" {
		t.Errorf("expected 'pre', got %q", opts.Prefix)
	}
}

func TestListSuffixOption(t *testing.T) {
	opts := NewListOptions(ListSuffix("suf"))
	if opts.Suffix != "suf" {
		t.Errorf("expected 'suf', got %q", opts.Suffix)
	}
}

func TestListLimitOption(t *testing.T) {
	opts := NewListOptions(ListLimit(10))
	if opts.Limit != 10 {
		t.Errorf("expected 10, got %d", opts.Limit)
	}
}

func TestListOffsetOption(t *testing.T) {
	opts := NewListOptions(ListOffset(5))
	if opts.Offset != 5 {
		t.Errorf("expected 5, got %d", opts.Offset)
	}
}

func TestListNamespaceOption(t *testing.T) {
	opts := NewListOptions(ListNamespace("ns"))
	if opts.Namespace != "ns" {
		t.Errorf("expected 'ns', got %q", opts.Namespace)
	}
}

func TestListTimeoutOption(t *testing.T) {
	td := time.Second
	opts := NewListOptions(ListTimeout(td))
	if opts.Timeout != td {
		t.Errorf("expected %v, got %v", td, opts.Timeout)
	}
}

func TestExistsContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey, "val")
	opts := NewExistsOptions(ExistsContext(ctx))
	if opts.Context.Value(testContextKey) != "val" {
		t.Error("ExistsContext option not set")
	}
}

func TestExistsNamespaceOption(t *testing.T) {
	opts := NewExistsOptions(ExistsNamespace("ns"))
	if opts.Namespace != "ns" {
		t.Errorf("expected 'ns', got %q", opts.Namespace)
	}
}

func TestExistsNameOption(t *testing.T) {
	opts := NewExistsOptions(ExistsName("test"))
	if opts.Name != "test" {
		t.Errorf("expected 'test', got %q", opts.Name)
	}
}

func TestExistsTimeoutOption(t *testing.T) {
	td := time.Second
	opts := NewExistsOptions(ExistsTimeout(td))
	if opts.Timeout != td {
		t.Errorf("expected %v, got %v", td, opts.Timeout)
	}
}
