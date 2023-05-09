package config_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.unistack.org/micro/v4/config"
	mtime "go.unistack.org/micro/v4/util/time"
)

type cfg struct {
	StringValue    string `default:"string_value"`
	IgnoreValue    string `json:"-"`
	StructValue    *cfgStructValue
	IntValue       int             `default:"99"`
	DurationValue  time.Duration   `default:"10s"`
	MDurationValue mtime.Duration  `default:"10s"`
	MapValue       map[string]bool `default:"key1=true,key2=false"`
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
