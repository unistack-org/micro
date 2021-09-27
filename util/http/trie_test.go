package http

import (
	"net/http"
	"testing"
)

func TestTrieNoMatchMethod(t *testing.T) {
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create/{id}", nil)

	_, _, ok := tr.Search(http.MethodPost, "/v1/create")
	if ok {
		t.Fatalf("must be not found error")
	}
}

type handler struct{}

func TestTrieMatchRegexp(t *testing.T) {
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create/{category}/{id:[0-9]+}", &handler{})

	_, params, ok := tr.Search(http.MethodPut, "/v1/create/test_cat/12345")
	if !ok {
		t.Fatalf("route not found")
	} else if len(params) != 2 {
		t.Fatalf("param matching error %v", params)
	} else if params["category"] != "test_cat" {
		t.Fatalf("param matching error %v", params)
	}
}

func TestTrieMatchRegexpFail(t *testing.T) {
	tr := NewTrie()
	tr.Insert([]string{http.MethodPut}, "/v1/create/{id:[a-z]+}", &handler{})

	_, _, ok := tr.Search(http.MethodPut, "/v1/create/12345")
	if ok {
		t.Fatalf("route must not be not found")
	}
}
