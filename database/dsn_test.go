package database

import (
	"testing"
)

func TestParseDSN(t *testing.T) {
	cfg, err := ParseDSN("postgres://username:p@ssword#@host:12345/dbname?key1=val2&key2=val2")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Password != "p@ssword#" {
		t.Fatalf("parsing error")
	}
}
