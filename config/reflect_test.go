package config_test

import (
	"testing"

	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type Config struct {
	Value     string
	SubConfig *SubConfig
	Config    *Config
}

type SubConfig struct {
	Value string
}

func TestReflect(t *testing.T) {
	cfg1 := &Config{Value: "cfg1", Config: &Config{Value: "cfg1_1"}, SubConfig: &SubConfig{Value: "cfg1"}}
	cfg2, err := rutil.Zero(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dst: %#+v\n", cfg2)
}
