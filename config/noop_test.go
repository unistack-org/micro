package config_test

import "testing"
import "context"
import "fmt"
import "github.com/unistack-org/micro/v3/config"

type Cfg struct {
	Value string
}

func TestNoop(t *testing.T) {
	ctx := context.Background()
	conf := &Cfg{}
	blfn := func(ctx context.Context, cfg config.Config) error {
		conf, ok := cfg.Options().Struct.(*Cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", cfg.Options())
		}
		conf.Value = "before_load"
		return nil
	}
	alfn := func(ctx context.Context, cfg config.Config) error {
		conf, ok := cfg.Options().Struct.(*Cfg)
		if !ok {
			return fmt.Errorf("failed to get Struct from options: %v", cfg.Options())
		}
		conf.Value = "after_load"
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
	if conf.Value != "before_load" {
		t.Fatal("BeforeLoad option not working")
	}

	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}

	for _, fn := range cfg.Options().AfterLoad {
		if err := fn(ctx, cfg); err != nil {
			t.Fatal(err)
		}
	}
	if conf.Value != "after_load" {
		t.Fatal("AfterLoad option not working")
	}

}
