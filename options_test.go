package micro

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/config"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/register"
	"go.unistack.org/micro/v5/router"
	"go.unistack.org/micro/v5/server"
	"go.unistack.org/micro/v5/store"
	"go.unistack.org/micro/v5/tracer"
)

// --- Option functions ---

func TestOptionBroker(t *testing.T) {
	b := broker.NewBroker()
	opts := NewOptions(Broker(b))
	_ = opts // Broker() calls srv/cli Init internally — just verify no panic
}

func TestOptionBrokerWithClientServer(t *testing.T) {
	b := broker.NewBroker()
	// BrokerClient / BrokerServer sub-options narrow which components get set
	opts := NewOptions(Broker(b, BrokerClient("noop"), BrokerServer("noop")))
	_ = opts
}

func TestOptionClients(t *testing.T) {
	c := client.NewClient()
	opts := NewOptions(Clients(c))
	require.Len(t, opts.Clients, 1)
	assert.Equal(t, c, opts.Clients[0])
}

func TestOptionServers(t *testing.T) {
	s := server.NewServer()
	opts := NewOptions(Servers(s))
	require.Len(t, opts.Servers, 1)
	assert.Equal(t, s, opts.Servers[0])
}

func TestOptionServer(t *testing.T) {
	s := server.NewServer()
	opts := NewOptions(Server(s))
	require.Len(t, opts.Servers, 1)
	assert.Equal(t, s, opts.Servers[0])
}

func TestOptionStore(t *testing.T) {
	s := store.NewStore()
	opts := NewOptions(Store(s))
	require.Len(t, opts.Stores, 1)
	assert.Equal(t, s, opts.Stores[0])
}

func TestOptionStores(t *testing.T) {
	s := store.NewStore()
	opts := NewOptions(Stores(s))
	require.Len(t, opts.Stores, 1)
}

func TestOptionLogger(t *testing.T) {
	l := logger.NewLogger()
	opts := NewOptions(Logger(l))
	_ = opts
}

func TestOptionMeter(t *testing.T) {
	m := meter.NewMeter()
	opts := NewOptions(Meter(m))
	require.Len(t, opts.Meters, 1)
	assert.Equal(t, m, opts.Meters[0])
}

func TestOptionMeters(t *testing.T) {
	m := meter.NewMeter()
	opts := NewOptions(Meters(m))
	require.Len(t, opts.Meters, 1)
}

func TestOptionConfig(t *testing.T) {
	c := config.NewConfig()
	opts := NewOptions(Config(c))
	require.Len(t, opts.Configs, 1)
	assert.Equal(t, c, opts.Configs[0])
}

func TestOptionConfigs(t *testing.T) {
	c := config.NewConfig()
	opts := NewOptions(Configs(c))
	require.Len(t, opts.Configs, 1)
}

func TestOptionContext(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "v")
	opts := NewOptions(Context(ctx))
	assert.Equal(t, ctx, opts.Context)
}

func TestOptionVersion(t *testing.T) {
	opts := NewOptions(Version("1.2.3"))
	assert.Equal(t, "1.2.3", opts.Version)
}

func TestOptionMetadata(t *testing.T) {
	md := metadata.Metadata{"key": []string{"val"}}
	opts := NewOptions(Metadata(md))
	assert.Equal(t, []string{"val"}, opts.Metadata["key"])
}

func TestOptionRegister(t *testing.T) {
	r := register.NewRegister()
	opts := NewOptions(Register(r))
	_ = opts
}

func TestOptionRegisterSubOptions(t *testing.T) {
	r := register.NewRegister()
	opts := NewOptions(Register(r, RegisterRouter("noop"), RegisterServer("noop"), RegisterBroker("noop")))
	_ = opts
}

func TestOptionTracer(t *testing.T) {
	tr := tracer.NewTracer()
	opts := NewOptions(Tracer(tr))
	_ = opts
}

func TestOptionTracerSubOptions(t *testing.T) {
	tr := tracer.NewTracer()
	opts := NewOptions(Tracer(tr, TracerClient("noop"), TracerServer("noop"), TracerBroker("noop"), TracerStore("noop")))
	_ = opts
}

func TestOptionRouter(t *testing.T) {
	r := router.NewRouter()
	opts := NewOptions(Router(r))
	_ = opts
}

func TestOptionRouterWithClient(t *testing.T) {
	r := router.NewRouter()
	opts := NewOptions(Router(r, RouterClient("noop")))
	_ = opts
}

func TestOptionAddress(t *testing.T) {
	opts := NewOptions(Address(":0"))
	_ = opts
}

func TestOptionAddressMultipleServers(t *testing.T) {
	s1 := server.NewServer()
	s2 := server.NewServer()
	// Multiple servers — should return error (not applied)
	opt := Address(":8080")
	o := NewOptions(Servers(s1, s2))
	err := opt(&o)
	require.Error(t, err)
}

func TestOptionRegisterTTL(t *testing.T) {
	opts := NewOptions(RegisterTTL(5 * time.Second))
	_ = opts
}

func TestOptionRegisterTTLWithServer(t *testing.T) {
	opts := NewOptions(RegisterTTL(5*time.Second, RegisterServer("noop")))
	_ = opts
}

func TestOptionRegisterInterval(t *testing.T) {
	opts := NewOptions(RegisterInterval(10 * time.Second))
	_ = opts
}

func TestOptionRegisterIntervalWithServer(t *testing.T) {
	opts := NewOptions(RegisterInterval(10*time.Second, RegisterServer("noop")))
	_ = opts
}

func TestOptionBeforeStart(t *testing.T) {
	called := false
	fn := func(_ context.Context) error { called = true; return nil }
	opts := NewOptions(BeforeStart(fn))
	require.Len(t, opts.BeforeStart, 1)
	require.NoError(t, opts.BeforeStart[0](context.Background()))
	assert.True(t, called)
}

func TestOptionBeforeStop(t *testing.T) {
	called := false
	fn := func(_ context.Context) error { called = true; return nil }
	opts := NewOptions(BeforeStop(fn))
	require.Len(t, opts.BeforeStop, 1)
	require.NoError(t, opts.BeforeStop[0](context.Background()))
	assert.True(t, called)
}

func TestOptionAfterStart(t *testing.T) {
	called := false
	fn := func(_ context.Context) error { called = true; return nil }
	opts := NewOptions(AfterStart(fn))
	require.Len(t, opts.AfterStart, 1)
	require.NoError(t, opts.AfterStart[0](context.Background()))
	assert.True(t, called)
}

func TestOptionAfterStop(t *testing.T) {
	called := false
	fn := func(_ context.Context) error { called = true; return nil }
	opts := NewOptions(AfterStop(fn))
	require.Len(t, opts.AfterStop, 1)
	require.NoError(t, opts.AfterStop[0](context.Background()))
	assert.True(t, called)
}

// --- service.Live / Ready / Health ---

func TestServiceLiveReadyHealth(t *testing.T) {
	svc := NewService()
	assert.True(t, svc.Live())
	assert.True(t, svc.Ready())
	assert.True(t, svc.Health())
}

// --- NewContext nil path ---

func TestNewContextNil(t *testing.T) {
	//nolint:staticcheck
	ctx := NewContext(nil, NewService())
	svc, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, svc)
}

func TestFromContextNil(t *testing.T) {
	//nolint:staticcheck
	svc, ok := FromContext(nil)
	assert.False(t, ok)
	assert.Nil(t, svc)
}
