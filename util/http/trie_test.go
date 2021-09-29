package http

import (
	"net/http"
	"testing"
)

func TestTrieContentType(t *testing.T) {
	type handler struct {
		name string
	}
	tr := NewTrie()
	tr.Insert([]string{"application/json"}, "/v1/create/{id}", &handler{name: "test"})

	h, _, ok := tr.Search("application/json", "/v1/create/12")
	if !ok {
		t.Fatalf("must be found error")
	}
	if h.(*handler).name != "test" {
		t.Fatalf("invalid handler %v", h)
	}
	_, _, ok = tr.Search("text/xml", "/v1/create/12")
	if ok {
		t.Fatalf("must be not found error")
	}
}

func TestTrieNoMatchMethod(t *testing.T) {
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create/{id}", nil)

	_, _, ok := tr.Search(http.MethodPost, "/v1/create")
	if ok {
		t.Fatalf("must be not found error")
	}
}

func TestTrieMatchRegexp(t *testing.T) {
	type handler struct{}
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create/{category}/{id:[0-9]+}", &handler{})

	_, params, ok := tr.Search(http.MethodPut, "/v1/create/test_cat/12345")
	switch {
	case !ok:
		t.Fatalf("route not found")
	case len(params) != 2:
		t.Fatalf("param matching error %v", params)
	case params["category"] != "test_cat":
		t.Fatalf("param matching error %v", params)
	}
}

func TestTrieMatchRegexpFail(t *testing.T) {
	type handler struct{}
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create/{id:[a-z]+}", &handler{})

	_, _, ok := tr.Search(http.MethodPut, "/v1/create/12345")
	if ok {
		t.Fatalf("route must not be not found")
	}
}

func TestTrieMatchLongest(t *testing.T) {
	type handler struct {
		name string
	}
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create", &handler{name: "first"})
	tr.Insert([]string{http.MethodPut}, "/v1/create/{id:[0-9]+}", &handler{name: "second"})
	if h, _, ok := tr.Search(http.MethodPut, "/v1/create/12345"); !ok {
		t.Fatalf("route must be found")
	} else if h.(*handler).name != "second" {
		t.Fatalf("invalid handler found: %s != %s", h.(*handler).name, "second")
	}
	if h, _, ok := tr.Search(http.MethodPut, "/v1/create"); !ok {
		t.Fatalf("route must be found")
	} else if h.(*handler).name != "first" {
		t.Fatalf("invalid handler found: %s != %s", h.(*handler).name, "first")
	}
}
