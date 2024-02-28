package options_test

import (
	"fmt"
	"testing"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/util/reflect"
)

func TestAddress(t *testing.T) {
	var err error

	type s struct {
		Address []string
	}

	src := &s{}
	var opts []options.Option
	opts = append(opts, options.Address("host:port"))

	for _, o := range opts {
		if err = o(src); err != nil {
			t.Fatal(err)
		}
	}

	if src.Address[0] != "host:port" {
		t.Fatal("failed to set Address")
	}
}

func TestCodecs(t *testing.T) {
	var err error

	type s struct {
		Codecs map[string]codec.Codec
	}

	src := &s{}
	var opts []options.Option
	c := codec.NewCodec()
	opts = append(opts, options.Codecs("text/html", c))

	for _, o := range opts {
		if err = o(src); err != nil {
			t.Fatal(err)
		}
	}

	for k, v := range src.Codecs {
		if k != "text/html" || v != c {
			continue
		}
		return
	}

	t.Fatalf("failed to set Codecs")
}

func TestLabels(t *testing.T) {
	type str1 struct {
		Labels []string
	}
	type str2 struct {
		Labels []interface{}
	}

	x1 := &str1{}

	if err := options.Labels("one", "two")(x1); err != nil {
		t.Fatal(err)
	}
	if len(x1.Labels) != 2 {
		t.Fatal("failed to set labels")
	}
	x2 := &str2{}
	if err := options.Labels("key", "val")(x2); err != nil {
		t.Fatal(err)
	}
	if len(x2.Labels) != 2 {
		t.Fatal("failed to set labels")
	}
	if x2.Labels[0] != "key" {
		t.Fatal("failed to set labels")
	}
}

func TestMetadataAny(t *testing.T) {
	type s struct {
		Metadata metadata.Metadata
	}

	testCases := []struct {
		Name     string
		Data     any
		Expected metadata.Metadata
	}{
		{
			"strings_even",
			[]string{"key1", "val1", "key2", "val2"},
			metadata.Metadata{
				"Key1": "val1",
				"Key2": "val2",
			},
		},
		{
			"strings_odd",
			[]string{"key1", "val1", "key2"},
			metadata.Metadata{
				"Key1": "val1",
			},
		},
		{
			"map",
			map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
			metadata.Metadata{
				"Key1": "val1",
				"Key2": "val2",
			},
		},
		{
			"metadata.Metadata",
			metadata.Metadata{
				"key1": "val1",
				"key2": "val2",
			},
			metadata.Metadata{
				"Key1": "val1",
				"Key2": "val2",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			src := &s{}
			var opts []options.Option
			opts = append(opts, options.MetadataAny(tt.Data))
			for _, o := range opts {
				if err := o(src); err != nil {
					t.Fatal(err)
				}
				if !reflect.Equal(tt.Expected, src.Metadata) {
					t.Fatal(fmt.Sprintf("expected: %v, actual: %v", tt.Expected, src.Metadata))
				}
			}
		})
	}
}
