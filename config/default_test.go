package config_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.unistack.org/micro/v4/config"
	mtime "go.unistack.org/micro/v4/util/time"
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
	st := cfg{}

	tVal := reflect.TypeOf(st)
	for i := 0; i < tVal.NumField(); i++ {
		field := tVal.Field(i)
		fmt.Printf("Field: %s, Offset: %d, Size: %d\n", field.Name, field.Offset, field.Type.Size())
	}
}
