package dns

import (
	"net"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	net.DefaultResolver = NewNetResolver(PreferIPV4(true))

	_, err := net.LookupHost("unistack.org")
	if err != nil {
		t.Fatal(err)
	}

	_, err = net.LookupHost("unistack.org")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOptions(t *testing.T) {
	opts := Options{}
	MaxCacheEntries(10)(&opts)
	if opts.MaxCacheEntries != 10 {
		t.Errorf("MaxCacheEntries: want 10, got %d", opts.MaxCacheEntries)
	}
	MaxCacheTTL(5 * time.Second)(&opts)
	if opts.MaxCacheTTL != 5*time.Second {
		t.Errorf("MaxCacheTTL: want 5s, got %v", opts.MaxCacheTTL)
	}
	MinCacheTTL(1 * time.Second)(&opts)
	if opts.MinCacheTTL != 1*time.Second {
		t.Errorf("MinCacheTTL: want 1s, got %v", opts.MinCacheTTL)
	}
	NegativeCache(true)(&opts)
	if !opts.NegativeCache {
		t.Error("NegativeCache: want true, got false")
	}
	Timeout(2 * time.Second)(&opts)
	if opts.Timeout != 2*time.Second {
		t.Errorf("Timeout: want 2s, got %v", opts.Timeout)
	}
	PreferIPV6(true)(&opts)
	if !opts.PreferIPV6 {
		t.Error("PreferIPV6: want true, got false")
	}
	PreferIPV4(true)(&opts)
	if !opts.PreferIPV4 {
		t.Error("PreferIPV4: want true, got false")
	}
}

func TestNewNetResolver_PreferIPV4AndIPV6(t *testing.T) {
	// Both prefer flags set — should fall through to "udp"
	r := NewNetResolver(PreferIPV4(true), PreferIPV6(true))
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNewNetDialer_MaxCacheEntries(t *testing.T) {
	// Zero MaxCacheEntries should default to DefaultMaxCacheEntries
	dial := NewNetDialer(nil, MaxCacheEntries(0))
	if dial == nil {
		t.Fatal("expected non-nil dialer")
	}
}

func TestCache_GetInvalidMsg(t *testing.T) {
	c := &cache{}
	// messages shorter than 12 bytes are invalid
	if got := c.get("short"); got != "" {
		t.Errorf("expected empty result for short message, got %q", got)
	}
}

func TestCache_GetHighBit(t *testing.T) {
	c := &cache{}
	// byte at index 2 >= 0x7f means it's a response, not a query — skip cache
	msg := "\x00\x01\xff\x00\x00\x00\x00\x00\x00\x00\x00\x00"
	if got := c.get(msg); got != "" {
		t.Errorf("expected empty result, got %q", got)
	}
}

func TestCache_GetNilEntries(t *testing.T) {
	c := &cache{}
	// valid-looking 12-byte query with low byte[2]
	msg := "\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
	if got := c.get(msg); got != "" {
		t.Errorf("expected empty result when entries nil, got %q", got)
	}
}
