package net

import (
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	fn := func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	}

	// try to create a number of listeners
	for range 10 {
		l, err := Listen("localhost:10000-11000", fn)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = l.Close()
		}()
	}

	// TODO nats case test
	// natsAddr := "_INBOX.bID2CMRvlNp0vt4tgNBHWf"
	// Expect addr DO NOT has extra ":" at the end!
}

func TestListen_SinglePort(t *testing.T) {
	fn := func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	}
	l, err := Listen("localhost:0", fn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
}

func TestListen_InvalidAddr(t *testing.T) {
	fn := func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	}
	_, err := Listen("not-valid-addr:5000-6000", fn)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestListen_RangeExhausted(t *testing.T) {
	called := 0
	failFn := func(addr string) (net.Listener, error) {
		called++
		return nil, &net.OpError{Op: "listen", Err: net.UnknownNetworkError("already in use")}
	}
	_, err := Listen("localhost:19999-19999", failFn)
	if err == nil {
		t.Fatal("expected error when range is exhausted")
	}
	if called == 0 {
		t.Fatal("expected fn to be called at least once")
	}
}

func TestHostPort_IPv4(t *testing.T) {
	got := HostPort("192.168.1.1", 8080)
	want := "192.168.1.1:8080"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestHostPort_IPv6(t *testing.T) {
	got := HostPort("::1", 8080)
	want := "[::1]:8080"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestHostPort_EmptyStringPort(t *testing.T) {
	// empty string port => return host only
	got := HostPort("myqueue", "")
	if got != "myqueue" {
		t.Errorf("want %q, got %q", "myqueue", got)
	}
}

func TestHostPort_ZeroIntPort_NonIP(t *testing.T) {
	// zero int port with non-IP host => return host only
	got := HostPort("_INBOX.foo", 0)
	if got != "_INBOX.foo" {
		t.Errorf("want %q, got %q", "_INBOX.foo", got)
	}
}

func TestHostPort_ZeroIntPort_IP(t *testing.T) {
	// zero int port with valid IP => should still append port
	got := HostPort("127.0.0.1", 0)
	want := "127.0.0.1:0"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
