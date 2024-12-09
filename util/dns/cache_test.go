package dns

import (
	"net"
	"testing"
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
