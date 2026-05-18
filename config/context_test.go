package config

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/tracer"
)

func TestFromNilContext(t *testing.T) {
	// nolint: staticcheck
	c, ok := FromContext(nil)
	if ok || c != nil {
		t.Fatal("FromContext not works")
	}
}

func TestNewNilContext(t *testing.T) {
	// nolint: staticcheck
	ctx := NewContext(nil, NewConfig())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), configKey{}, NewConfig())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), NewConfig())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestMustContext(t *testing.T) {
	ctx := NewContext(context.TODO(), NewConfig())
	c := MustContext(ctx)
	if c == nil {
		t.Fatal("MustContext returned nil")
	}
}

func TestMustContextPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from MustContext on empty context")
		}
	}()
	MustContext(context.TODO())
}

func TestSetOption(t *testing.T) {
	type key struct{}
	o := SetOption(key{}, "test")
	opts := &Options{}
	o(opts)

	if v, ok := opts.Context.Value(key{}).(string); !ok || v == "" {
		t.Fatal("SetOption not works")
	}
}

func TestSetSaveOption(t *testing.T) {
	type key struct{}
	o := SetSaveOption(key{}, "test")
	opts := &SaveOptions{}
	o(opts)

	if v, ok := opts.Context.Value(key{}).(string); !ok || v == "" {
		t.Fatal("SetSaveOption not works")
	}
}

func TestSetLoadOption(t *testing.T) {
	type key struct{}
	o := SetLoadOption(key{}, "test")
	opts := &LoadOptions{}
	o(opts)

	if v, ok := opts.Context.Value(key{}).(string); !ok || v == "" {
		t.Fatal("SetLoadOption not works")
	}
}

func TestSetWatchOption(t *testing.T) {
	type key struct{}
	o := SetWatchOption(key{}, "test")
	opts := &WatchOptions{}
	o(opts)

	if v, ok := opts.Context.Value(key{}).(string); !ok || v == "" {
		t.Fatal("SetWatchOption not works")
	}
}

func TestOptions(t *testing.T) {
	ctx := context.Background()
	c := codec.NewCodec()
	l := logger.DefaultLogger
	tr := tracer.DefaultTracer

	opts := NewOptions(
		AllowFail(true),
		Context(ctx),
		Codec(c),
		Logger(l),
		Tracer(tr),
		Struct(struct{}{}),
		StructTag("json"),
		Name("myconfig"),
		BeforeInit(func(context.Context, Config) error { return nil }),
		AfterInit(func(context.Context, Config) error { return nil }),
		BeforeLoad(func(context.Context, Config) error { return nil }),
		AfterLoad(func(context.Context, Config) error { return nil }),
		BeforeSave(func(context.Context, Config) error { return nil }),
		AfterSave(func(context.Context, Config) error { return nil }),
	)

	if !opts.AllowFail {
		t.Error("AllowFail not set")
	}
	if opts.Codec == nil {
		t.Error("Codec not set")
	}
	if opts.Logger == nil {
		t.Error("Logger not set")
	}
	if opts.Tracer == nil {
		t.Error("Tracer not set")
	}
	if opts.StructTag != "json" {
		t.Errorf("StructTag = %q", opts.StructTag)
	}
	if opts.Name != "myconfig" {
		t.Errorf("Name = %q", opts.Name)
	}
	if len(opts.BeforeInit) == 0 {
		t.Error("BeforeInit not set")
	}
	if len(opts.AfterInit) == 0 {
		t.Error("AfterInit not set")
	}
	if len(opts.BeforeLoad) == 0 {
		t.Error("BeforeLoad not set")
	}
	if len(opts.AfterLoad) == 0 {
		t.Error("AfterLoad not set")
	}
	if len(opts.BeforeSave) == 0 {
		t.Error("BeforeSave not set")
	}
	if len(opts.AfterSave) == 0 {
		t.Error("AfterSave not set")
	}
}

func TestLoadOptions(t *testing.T) {
	src := struct{ V string }{V: "x"}
	opts := NewLoadOptions(
		LoadOverride(true),
		LoadAppend(true),
		LoadStruct(&src),
	)
	if !opts.Override {
		t.Error("Override not set")
	}
	if !opts.Append {
		t.Error("Append not set")
	}
	if opts.Struct == nil {
		t.Error("Struct not set")
	}
}

func TestSaveOptions(t *testing.T) {
	src := struct{ V string }{V: "x"}
	opts := NewSaveOptions(SaveStruct(&src))
	if opts.Struct == nil {
		t.Error("Struct not set")
	}
}

func TestWatchOptions(t *testing.T) {
	opts := NewWatchOptions(
		WatchContext(context.Background()),
		WatchCoalesce(true),
		WatchInterval(1*time.Second, 2*time.Second),
		WatchStruct(struct{}{}),
	)
	if !opts.Coalesce {
		t.Error("Coalesce not set")
	}
	if opts.MinInterval != 1*time.Second {
		t.Errorf("MinInterval = %v", opts.MinInterval)
	}
	if opts.MaxInterval != 2*time.Second {
		t.Errorf("MaxInterval = %v", opts.MaxInterval)
	}
	if opts.Struct == nil {
		t.Error("Struct not set")
	}
}

func TestConfigStringAndName(t *testing.T) {
	cfg := NewConfig(Name("testcfg"))
	if cfg.String() != "default" {
		t.Errorf("String() = %q", cfg.String())
	}
	if cfg.Name() != "testcfg" {
		t.Errorf("Name() = %q", cfg.Name())
	}
}

func TestConfigSave(t *testing.T) {
	cfg := NewConfig()
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConfigWatch(t *testing.T) {
	cfg := NewConfig()
	_, err := cfg.Watch(context.Background())
	if err != ErrWatcherNotImplemented {
		t.Fatalf("expected ErrWatcherNotImplemented, got %v", err)
	}
}

func TestConfigLoad(t *testing.T) {
	var s struct{ V string }
	cs := []Config{NewConfig(Struct(&s))}
	if err := Load(context.Background(), cs); err != nil {
		t.Fatal(err)
	}
}

func TestValidateNil(t *testing.T) {
	if err := Validate(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestValidateNonStruct(t *testing.T) {
	v := "hello"
	if err := Validate(context.Background(), v); err != nil {
		t.Fatal(err)
	}
}
