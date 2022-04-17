package config

import (
	"context"
	"testing"
)

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
