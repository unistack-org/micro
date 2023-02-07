package time

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	var td time.Duration
	var err error

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
