package config_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.unistack.org/micro/v5/config"
	mtime "go.unistack.org/micro/v5/util/time"
)

type cfg struct {
	MapValue    map[string]bool `default:"key1=true,key2=false"`
	StructValue *cfgStructValue

	StringValue string `default:"string_value"`
	IgnoreValue string `json:"-"`
	UUIDValue   string `default:"micro:generate uuid"`
	IDValue     string `default:"micro:generate id"`

	DurationValue  time.Duration  `default:"10s"`
	MDurationValue mtime.Duration `default:"10s"`
	IntValue       int            `default:"99"`
}

type cfgStructValue struct {
	StringValue string `default:"string_value"`
}

func (c *cfg) Validate() error {
	if c.IntValue != 10 {
		return fmt.Errorf("invalid IntValue %d != %d", 10, c.IntValue)
	}
	return nil
}

func (c *cfgStructValue) Validate() error {
	if c.StringValue != "string_value" {
		return fmt.Errorf("invalid StringValue %s != %s", "string_value", c.StringValue)
	}
	return nil
}

type testHook struct {
	f bool
}

func (t *testHook) Load(fn config.FuncLoad) config.FuncLoad {
	return func(ctx context.Context, opts ...config.LoadOption) error {
		t.f = true
		return fn(ctx, opts...)
	}
}

func TestHook(t *testing.T) {
	h := &testHook{}

	c := config.NewConfig(config.Struct(h), config.Hooks(config.HookLoad(h.Load)))

	if err := c.Init(); err != nil {
		t.Fatal(err)
	}

	if err := c.Load(context.TODO()); err != nil {
		t.Fatal(err)
	}

	if !h.f {
		t.Fatal("hook not works")
	}
}

func TestDefault(t *testing.T) {
	ctx := context.Background()
	conf := &cfg{IntValue: 10}
	blfn := func(_ context.Context, c config.Config) error {
		nconf, ok := c.Options().Struct.(*cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", c.Options())
		}
		nconf.StringValue = "before_load"
		return nil
	}
	alfn := func(_ context.Context, c config.Config) error {
		nconf, ok := c.Options().Struct.(*cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", c.Options())
		}
		nconf.StringValue = "after_load"
		return nil
	}

	cfg := config.NewConfig(config.Struct(conf), config.BeforeLoad(blfn), config.AfterLoad(alfn))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}
	if conf.StringValue != "after_load" {
		t.Fatal("AfterLoad option not working")
	}
	if len(conf.MapValue) != 2 {
		t.Fatalf("map value invalid: %#+v\n", conf.MapValue)
	}

	if conf.UUIDValue == "" {
		t.Fatalf("uuid value empty")
	} else if len(conf.UUIDValue) != 36 {
		t.Fatalf("uuid value invalid: %s", conf.UUIDValue)
	}

	if conf.IDValue == "" {
		t.Fatalf("id value empty")
	}
	_ = conf
	// t.Logf("%#+v\n", conf)
}

func TestValidate(t *testing.T) {
	ctx := context.Background()
	conf := &cfg{IntValue: 10}
	cfg := config.NewConfig(config.Struct(conf))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}

	if err := config.Validate(ctx, conf); err != nil {
		t.Fatal(err)
	}
}

func Test_SizeOf(t *testing.T) {

	tVal := reflect.TypeFor[cfg]()
	for i := 0; i < tVal.NumField(); i++ {
		field := tVal.Field(i)
		fmt.Printf("Field: %s, Offset: %d, Size: %d\n", field.Name, field.Offset, field.Type.Size())
	}
}

type numericCfg struct {
	F32  float32  `default:"1.5"`
	F64  float64  `default:"2.5"`
	I64  int64    `default:"64"`
	U64  uint64   `default:"64"`
	Strs []string `default:"a,b,c"`
}

func TestNumericDefaults(t *testing.T) {
	ctx := context.Background()
	conf := &numericCfg{}
	cfg := config.NewConfig(config.Struct(conf))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}
	if conf.F32 != 1.5 {
		t.Errorf("F32 = %v", conf.F32)
	}
	if conf.F64 != 2.5 {
		t.Errorf("F64 = %v", conf.F64)
	}
	if conf.I64 != 64 {
		t.Errorf("I64 = %v", conf.I64)
	}
	if conf.U64 != 64 {
		t.Errorf("U64 = %v", conf.U64)
	}
	if len(conf.Strs) != 3 {
		t.Errorf("Strs = %v", conf.Strs)
	}
}

func TestLoadWithOptions(t *testing.T) {
	ctx := context.Background()
	conf := &numericCfg{}
	cfg := config.NewConfig(config.Struct(conf))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	// LoadOverride and LoadAppend exercise additional mergo options
	if err := cfg.Load(ctx, config.LoadOverride(true), config.LoadAppend(true)); err != nil {
		t.Fatal(err)
	}
}

func TestLoadWithStruct(t *testing.T) {
	ctx := context.Background()
	conf := &numericCfg{}
	cfg := config.NewConfig(config.Struct(conf))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	alt := &numericCfg{}
	if err := cfg.Load(ctx, config.LoadStruct(alt)); err != nil {
		t.Fatal(err)
	}
}

func TestConfigSkipLoad(t *testing.T) {
	ctx := context.Background()
	conf := &numericCfg{}
	skipFn := func(context.Context, config.Config) bool { return true }
	cfg := config.NewConfig(
		config.Struct(conf),
		func(o *config.Options) { o.SkipLoad = skipFn },
	)
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}
	// values stay zero because SkipLoad returned true
	if conf.F32 != 0 {
		t.Errorf("expected zero F32 with SkipLoad, got %v", conf.F32)
	}
}

func TestConfigSaveSkip(t *testing.T) {
	ctx := context.Background()
	skipFn := func(context.Context, config.Config) bool { return true }
	cfg := config.NewConfig(
		func(o *config.Options) { o.SkipSave = skipFn },
	)
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestMDurationDefault(t *testing.T) {
	ctx := context.Background()
	conf := &cfg{IntValue: 10}
	c := config.NewConfig(config.Struct(conf))
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c.Load(ctx); err != nil {
		t.Fatal(err)
	}
	if conf.MDurationValue == 0 {
		t.Error("MDurationValue not set")
	}
}

type uintCfg struct {
	U   uint   `default:"1"`
	U8  uint8  `default:"8"`
	U16 uint16 `default:"16"`
	U32 uint32 `default:"32"`
}

func TestUintDefaults(t *testing.T) {
	ctx := context.Background()
	conf := &uintCfg{}
	c := config.NewConfig(config.Struct(conf))
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c.Load(ctx); err != nil {
		t.Fatal(err)
	}
	if conf.U != 1 {
		t.Errorf("U = %v", conf.U)
	}
	if conf.U8 != 8 {
		t.Errorf("U8 = %v", conf.U8)
	}
	if conf.U16 != 16 {
		t.Errorf("U16 = %v", conf.U16)
	}
	if conf.U32 != 32 {
		t.Errorf("U32 = %v", conf.U32)
	}
}

func TestConfigLoadMultiple(t *testing.T) {
	ctx := context.Background()
	type s1 struct{ V string `default:"hello"` }
	type s2 struct{ N int `default:"42"` }
	a, b := &s1{}, &s2{}
	cs := []config.Config{
		config.NewConfig(config.Struct(a)),
		config.NewConfig(config.Struct(b)),
	}
	if err := config.Load(ctx, cs); err != nil {
		t.Fatal(err)
	}
	if a.V != "hello" {
		t.Errorf("s1.V = %q", a.V)
	}
	if b.N != 42 {
		t.Errorf("s2.N = %d", b.N)
	}
}

func TestConfigSaveWithHooks(t *testing.T) {
	ctx := context.Background()
	called := 0
	fn := func(context.Context, config.Config) error { called++; return nil }
	type s struct{ V string `default:"x"` }
	data := &s{}
	c := config.NewConfig(
		config.Struct(data),
		config.BeforeSave(fn),
		config.AfterSave(fn),
	)
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c.Save(ctx); err != nil {
		t.Fatal(err)
	}
	if called == 0 {
		t.Error("expected hooks to be called at least once")
	}
}

func TestConfigInitWithBeforeAfter(t *testing.T) {
	ctx := context.Background()
	called := 0
	fn := func(context.Context, config.Config) error { called++; return nil }
	type s struct{ V string `default:"x"` }
	data := &s{}
	c := config.NewConfig(
		config.Struct(data),
		config.BeforeInit(fn),
		config.AfterInit(fn),
	)
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c.Load(ctx); err != nil {
		t.Fatal(err)
	}
}
