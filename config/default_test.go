package config_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/unistack-org/micro/v3/config"
)

type Cfg struct {
	StringValue string `default:"string_value"`
	IntValue    int    `default:"99"`
}

func TestDefault(t *testing.T) {
	ctx := context.Background()
	conf := &Cfg{}
	blfn := func(ctx context.Context, cfg config.Config) error {
		conf, ok := cfg.Options().Struct.(*Cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", cfg.Options())
		}
		conf.StringValue = "before_load"
		return nil
	}
	alfn := func(ctx context.Context, cfg config.Config) error {
		conf, ok := cfg.Options().Struct.(*Cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", cfg.Options())
		}
		conf.StringValue = "after_load"
		return nil
	}

	cfg := config.NewConfig(config.Struct(conf), config.BeforeLoad(blfn), config.AfterLoad(alfn))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	for _, fn := range cfg.Options().BeforeLoad {
		if err := fn(ctx, cfg); err != nil {
			t.Fatal(err)
		}
	}
	if conf.StringValue != "before_load" {
		t.Fatal("BeforeLoad option not working")
	}

	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}

	if conf.StringValue != "string_value" || conf.IntValue != 99 {
		t.Fatalf("load failed: %#+v", conf)
	}

	for _, fn := range cfg.Options().AfterLoad {
		if err := fn(ctx, cfg); err != nil {
			t.Fatal(err)
		}
	}
	if conf.StringValue != "after_load" {
		t.Fatal("AfterLoad option not working")
	}

}
