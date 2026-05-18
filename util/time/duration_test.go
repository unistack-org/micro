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

func TestParseDurationEmpty(t *testing.T) {
	_, err := ParseDuration("")
	if err == nil {
		t.Fatal("expected error for empty duration string, got nil")
	}
}

func TestParseDurationMinutes(t *testing.T) {
	td, err := ParseDuration("30m")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td != 30*time.Minute {
		t.Fatalf("expected 30m got %s", td.String())
	}
}

func TestParseDurationSeconds(t *testing.T) {
	td, err := ParseDuration("45s")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td != 45*time.Second {
		t.Fatalf("expected 45s got %s", td.String())
	}
}

func TestParseDurationHours(t *testing.T) {
	td, err := ParseDuration("2h30m")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td != 2*time.Hour+30*time.Minute {
		t.Fatalf("expected 2h30m got %s", td.String())
	}
}

func TestParseDurationYear(t *testing.T) {
	td, err := ParseDuration("1y")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td < 365*24*time.Hour {
		t.Fatalf("expected at least 8760h got %s", td.String())
	}
}

func TestUnmarshalYAMLFloat(t *testing.T) {
	type str struct {
		TTL *Duration `yaml:"ttl"`
	}
	v := &str{}
	// Use a float (with decimal) so goccy/go-yaml decodes it as float64
	err := yaml.Unmarshal([]byte(`{"ttl": 10000000.0}`), v)
	if err != nil {
		t.Fatal(err)
	}
	if v.TTL == nil || *(v.TTL) != 10000000 {
		t.Fatalf("invalid duration %v != 10000000", v.TTL)
	}
}

func TestUnmarshalJSONFloat(t *testing.T) {
	type str struct {
		TTL Duration `json:"ttl"`
	}
	v := &str{}
	err := json.Unmarshal([]byte(`{"ttl":10000000}`), v)
	if err != nil {
		t.Fatal(err)
	}
	if v.TTL != 10000000 {
		t.Fatalf("invalid duration %v != 10000000", v.TTL)
	}
}

func TestUnmarshalJSONInvalid(t *testing.T) {
	type str struct {
		TTL Duration `json:"ttl"`
	}
	v := &str{}
	// boolean is not a valid duration type
	err := json.Unmarshal([]byte(`{"ttl":true}`), v)
	if err == nil {
		t.Fatal("expected error for bool duration, got nil")
	}
}

func TestUnmarshalYAMLInvalid(t *testing.T) {
	type str struct {
		TTL *Duration `yaml:"ttl"`
	}
	v := &str{}
	// Use an invalid duration string to trigger ParseDuration error
	err := yaml.Unmarshal([]byte(`{"ttl": "not_a_duration"}`), v)
	if err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
}
