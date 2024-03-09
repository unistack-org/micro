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
	var err error

	err = json.Unmarshal([]byte(`{"ttl":"10ms"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if v.TTL != 10000000 {
		t.Fatalf("invalid duration %v != 10000000", v.TTL)
	}

	err = json.Unmarshal([]byte(`{"ttl":"1y"}`), v)
	if err != nil {
		t.Fatal(err)
	} else if v.TTL != 31622400000000000 {
		t.Fatalf("invalid duration %v != 31536000000000000", v.TTL)
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
	td, err = ParseDuration("1y")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if td.String() != "8784h0m0s" {
		t.Fatalf("ParseDuration 1y != 8760h0m0s : %s", td.String())
	}
}
