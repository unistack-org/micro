package dns

import (
	"net"
	"testing"
)

func TestCache(t *testing.T) {
	net.DefaultResolver = NewNetResolver(PreferIPV4(true))

	addrs, err := net.LookupHost("unistack.org")
	if err != nil {
		t.Fatal(err)
	}

	addrs, err = net.LookupHost("unistack.org")
	if err != nil {
		t.Fatal(err)
	}

	_ = addrs
}
