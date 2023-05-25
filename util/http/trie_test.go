package http

import (
	"net/http"
	"testing"
)

func TestTrieWithDot(t *testing.T) {
	var err error
	type handler struct {
		name string
	}

	tr := NewTrie()
	if err = tr.Insert([]string{http.MethodGet}, "/*", &handler{name: "slash_astar"}); err != nil {
		t.Fatal(err)
	}
	if err = tr.Insert([]string{http.MethodGet}, "*", &handler{name: "astar"}); err != nil {
		t.Fatal(err)
	}
	if err = tr.Insert([]string{http.MethodGet}, "/.well-known/test", &handler{name: "well-known"}); err != nil {
		t.Fatal(err)
	}
	h, _, err := tr.Search(http.MethodGet, "/.well-known/test")
	if err != nil {
		t.Fatalf("unexpected error handler not found")
	}
	if h.(*handler).name != "well-known" {
		t.Fatalf("error matching route with wildcard")
	}
}

func TestTrieWildcardPathPrefix(t *testing.T) {
	var err error
	type handler struct {
		name string
	}
	tr := NewTrie()
	if err = tr.Insert([]string{http.MethodPost}, "/v1/update", &handler{name: "post_update"}); err != nil {
		t.Fatal(err)
	}
	if err = tr.Insert([]string{http.MethodPost}, "/v1/*", &handler{name: "post_create"}); err != nil {
		t.Fatal(err)
	}
	h, _, err := tr.Search(http.MethodPost, "/v1/test/one")
	if err != nil {
		t.Fatalf("unexpected error handler not found")
	}
	if h.(*handler).name != "post_create" {
		t.Fatalf("invalid handler %v", h)
	}
	h, _, err = tr.Search(http.MethodPost, "/v1/update")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if h.(*handler).name != "post_update" {
		t.Fatalf("invalid handler %v", h)
	}
	h, _, err = tr.Search(http.MethodPost, "/v1/update/some/{x}")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if h.(*handler).name != "post_create" {
		t.Fatalf("invalid handler %v", h)
	}
}

func TestTriePathPrefix(t *testing.T) {
	type handler struct {
		name string
	}
	tr := NewTrie()
	_ = tr.Insert([]string{http.MethodPost}, "/v1/create/{id}", &handler{name: "post_create"})
	_ = tr.Insert([]string{http.MethodPost}, "/v1/update/{id}", &handler{name: "post_update"})
	_ = tr.Insert([]string{http.MethodPost}, "/", &handler{name: "post_wildcard"})
	h, _, err := tr.Search(http.MethodPost, "/")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if h.(*handler).name != "post_wildcard" {
		t.Fatalf("invalid handler %v", h)
	}
}

func TestTrieFixedPattern(t *testing.T) {
	type handler struct {
		name string
	}
	tr := NewTrie()
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create/{id}", &handler{name: "pattern"})
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create/12", &handler{name: "fixed"})
	h, _, err := tr.Search(http.MethodPut, "/v1/create/12")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if h.(*handler).name != "fixed" {
		t.Fatalf("invalid handler %v", h)
	}
}

func TestTrieNoMatchMethod(t *testing.T) {
	tr := NewTrie()
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create/{id}", nil)
	_, _, err := tr.Search(http.MethodPost, "/v1/create")
	if err == nil && err != ErrNotFound {
		t.Fatalf("must be not found error")
	}
}

func TestTrieMatchRegexp(t *testing.T) {
	type handler struct{}
	tr := NewTrie()
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create/{category}/{id:[0-9]+}", &handler{})
	_, params, err := tr.Search(http.MethodPut, "/v1/create/test_cat/12345")
	switch {
	case err != nil:
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
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create/{id:[a-z]+}", &handler{})
	_, _, err := tr.Search(http.MethodPut, "/v1/create/12345")
	if err != ErrNotFound {
		t.Fatalf("route must not be not found")
	}
}

func TestTrieMatchLongest(t *testing.T) {
	type handler struct {
		name string
	}
	tr := NewTrie()
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create", &handler{name: "first"})
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create/{id:[0-9]+}", &handler{name: "second"})
	if h, _, err := tr.Search(http.MethodPut, "/v1/create/12345"); err != nil {
		t.Fatalf("route must be found")
	} else if h.(*handler).name != "second" {
		t.Fatalf("invalid handler found: %s != %s", h.(*handler).name, "second")
	}
	if h, _, err := tr.Search(http.MethodPut, "/v1/create"); err != nil {
		t.Fatalf("route must be found")
	} else if h.(*handler).name != "first" {
		t.Fatalf("invalid handler found: %s != %s", h.(*handler).name, "first")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	type handler struct{}
	tr := NewTrie()
	_ = tr.Insert([]string{http.MethodPut}, "/v1/create", &handler{})
	_, _, err := tr.Search(http.MethodPost, "/v1/create")
	if err != ErrMethodNotAllowed {
		t.Fatalf("route must be method not allowed: %v", err)
	}
	_, _, err = tr.Search(http.MethodPut, "/v1/create")
	if err != nil {
		t.Fatalf("route must be found: %v", err)
	}
}
