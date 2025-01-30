package options_test

import (
	"crypto/tls"
	"sync"
	"testing"

	"go.unistack.org/micro/v4/options"
)

type codec interface {
	Marshal(v interface{}, opts ...options.Option) ([]byte, error)
	Unmarshal(b []byte, v interface{}, opts ...options.Option) error
	String() string
}

func TestCodecs(t *testing.T) {
	type s struct {
		Codecs map[string]codec
	}

	wg := &sync.WaitGroup{}
	tc := &tls.Config{InsecureSkipVerify: true}
	opts := []options.Option{
		options.NewOption("Codecs")(wg),
		options.NewOption("TLSConfig")(tc),
	}

	src := &s{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}
}

func TestSpecial(t *testing.T) {
	type s struct {
		Wait      *sync.WaitGroup
		TLSConfig *tls.Config
	}

	wg := &sync.WaitGroup{}
	tc := &tls.Config{InsecureSkipVerify: true}
	opts := []options.Option{
		options.NewOption("Wait")(wg),
		options.NewOption("TLSConfig")(tc),
	}

	src := &s{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}

	if src.Wait == nil {
		t.Fatalf("failed to set Wait %#+v", src)
	}

	if src.TLSConfig == nil {
		t.Fatalf("failed to set TLSConfig %#+v", src)
	}

	if src.TLSConfig.InsecureSkipVerify != true {
		t.Fatalf("failed to set TLSConfig %#+v", src)
	}
}

func TestNested(t *testing.T) {
	type server struct {
		Address []string
	}
	type ownserver struct {
		server
		OwnField string
	}

	opts := []options.Option{
		options.Address("host:port"),
		options.NewOption("OwnField")("fieldval"),
	}

	src := &ownserver{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}

	if src.Address[0] != "host:port" {
		t.Fatalf("failed to set Address %#+v", src)
	}

	if src.OwnField != "fieldval" {
		t.Fatalf("failed to set OwnField %#+v", src)
	}
}

func TestAddress(t *testing.T) {
	type s struct {
		Address []string
	}

	opts := []options.Option{options.Address("host:port")}

	src := &s{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}

	if src.Address[0] != "host:port" {
		t.Fatalf("failed to set Address %#+v", src)
	}
}

func TestNewOption(t *testing.T) {
	type s struct {
		Address []string
	}

	opts := []options.Option{options.NewOption("Address")("host1:port1", "host2:port2")}

	src := &s{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}

	if src.Address[0] != "host1:port1" {
		t.Fatalf("failed to set Address %#+v", src)
	}
	if src.Address[1] != "host2:port2" {
		t.Fatalf("failed to set Address %#+v", src)
	}
}

func TestArray(t *testing.T) {
	type s struct {
		Address [1]string
	}

	opts := []options.Option{options.NewOption("Address")("host:port", "host1:port1")}

	src := &s{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}

	if src.Address[0] != "host:port" {
		t.Fatalf("failed to set Address %#+v", src)
	}
}

func TestMap(t *testing.T) {
	type s struct {
		Metadata map[string]string
	}

	opts := []options.Option{
		options.NewOption("Metadata")("key1", "val1"),
		options.NewOption("Metadata")(map[string]string{"key2": "val2"}),
	}

	src := &s{}

	if err := options.Apply(src, opts...); err != nil {
		t.Fatal(err)
	}

	if len(src.Metadata) != 2 {
		t.Fatalf("failed to set Metadata %#+v", src)
	}

	if src.Metadata["key1"] != "val1" {
		t.Fatalf("failed to set Metadata %#+v", src)
	}

	if src.Metadata["key2"] != "val2" {
		t.Fatalf("failed to set Metadata %#+v", src)
	}
}
