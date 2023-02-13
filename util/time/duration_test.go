package time

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

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

	err := json.Unmarshal([]byte(`{"ttl":"10ms"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if v.TTL != 10000000 {
		t.Fatalf("invalid duration %v != 10000000", v.TTL)
	}
}

func TestParseDuration(t *testing.T) {
	var td time.Duration
	var err error
	t.Skip()
	td, err = ParseDuration("14d4h")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td.String() != "336h0m0s" {
		t.Fatalf("ParseDuration 14d != 336h0m0s : %s", td.String())
	}

	td, err = ParseDuration("1y")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td.String() != "8760h0m0s" {
		t.Fatalf("ParseDuration 1y != 8760h0m0s : %s", td.String())
	}
}
