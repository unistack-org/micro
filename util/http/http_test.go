//go:build ignore
// +build ignore

package http

import (
	"io"
	"net"
	"net/http"
	"testing"

	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/register/memory"
	"go.unistack.org/micro/v4/router"
	regRouter "go.unistack.org/micro/v4/router/register"
)

func TestRoundTripper(t *testing.T) {
	m := memory.NewRegister()
	r := regRouter.NewRouter(router.Register(m))

	rt := NewRoundTripper(WithRouter(r))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	go http.Serve(l, nil)

	m.Register(&register.Service{
		Name: "example.com",
		Nodes: []*register.Node{
			{
				Id:      "1",
				Address: l.Addr().String(),
			},
		},
	})

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	w.Body.Close()

	if string(b) != "hello world" {
		t.Fatal("response is", string(b))
	}

	// test http request
	c := &http.Client{
		Transport: rt,
	}

	rsp, err := c.Get("http://example.com")
	if err != nil {
		t.Fatal(err)
	}

	b, err = io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}
	rsp.Body.Close()

	if string(b) != "hello world" {
		t.Fatal("response is", string(b))
	}
}
