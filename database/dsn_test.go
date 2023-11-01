package database

import (
	"net/url"
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

func TestFormatDSN(t *testing.T) {
	src := "postgres://username:p@ssword#@host:12345/dbname?key1=val2&key2=val2"
	cfg, err := ParseDSN(src)
	if err != nil {
		t.Fatal(err)
	}
	dst, err := url.PathUnescape(cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	if src != dst {
		t.Fatalf("\n%s\n%s", src, dst)
	}
}
