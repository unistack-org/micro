package time

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
)

func TestMarshalYAML(t *testing.T) {
	d := Duration(10000000)
	buf, err := yaml.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`10ms
`)) {
		t.Fatalf("invalid duration: %s != %s", buf, `10ms`)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	type str struct {
		TTL *Duration `yaml:"ttl"`
	}
	v := &str{}
	var err error

	err = yaml.Unmarshal([]byte(`{"ttl":"10ms"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if *(v.TTL) != 10000000 {
		t.Fatalf("invalid duration %v != 10000000", v.TTL)
	}

	err = yaml.Unmarshal([]byte(`{"ttl":"1d"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if *(v.TTL) != 86400000000000 {
		t.Fatalf("invalid duration %v != 86400000000000", *v.TTL)
	}
}

func TestMarshalJSON(t *testing.T) {
	d := Duration(10000000)
	buf, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`"10ms"`)) {
		t.Fatalf("invalid duration: %s != %s", buf, `"10ms"`)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type str struct {
		TTL Duration `json:"ttl"`
	}
	v := &str{}
	var err error

	err = json.Unmarshal([]byte(`{"ttl":"10ms"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if v.TTL != 10000000 {
		t.Fatalf("invalid duration %v != 10000000", v.TTL)
	}

	err = json.Unmarshal([]byte(`{"ttl":"1d"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if v.TTL != 86400000000000 {
		t.Fatalf("invalid duration %v != 86400000000000", v.TTL)
	}
}

func TestParseDuration(t *testing.T) {
	var td time.Duration
	var err error

	td, err = ParseDuration("14d4h")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td.String() != "340h0m0s" {
		t.Fatalf("ParseDuration 14d != 340h0m0s : %s", td.String())
	}
	td, err = ParseDuration("1d")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td.String() != "24h0m0s" {
		t.Fatalf("ParseDuration 1d != 24h0m0s : %s", td.String())
	}
}
