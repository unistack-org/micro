package config_test

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v3/config"
)

type Cfg struct {
	StringValue string `default:"string_value"`
	IgnoreValue string `json:"-"`
	StructValue struct {
		StringValue string `default:"string_value"`
	}
	IntValue int `default:"99"`
}

func TestDefault(t *testing.T) {
	ctx := context.Background()
	conf := &Cfg{IntValue: 10}
	blfn := func(ctx context.Context, cfg config.Config) error {
		nconf, ok := cfg.Options().Struct.(*Cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", cfg.Options())
		}
		nconf.StringValue = "before_load"
		return nil
	}
	alfn := func(ctx context.Context, cfg config.Config) error {
		nconf, ok := cfg.Options().Struct.(*Cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", cfg.Options())
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
	_ = conf
	// t.Logf("%#+v\n", conf)
}
