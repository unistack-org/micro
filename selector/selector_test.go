package selector_test

import (
	"testing"

	"go.unistack.org/micro/v5/selector"
	"go.unistack.org/micro/v5/selector/random"
)

func TestErrNoneAvailable(t *testing.T) {
	if selector.ErrNoneAvailable == nil {
		t.Fatal("expected non-nil error")
	}
	if selector.ErrNoneAvailable.Error() == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestSelectorTests(t *testing.T) {
	s := random.NewSelector()
	selector.Tests(t, s)
}

func TestNewSelectOptions(t *testing.T) {
	opts := selector.NewSelectOptions()
	_ = opts

	called := false
	customOpt := func(o *selector.SelectOptions) {
		called = true
	}
	opts2 := selector.NewSelectOptions(customOpt)
	_ = opts2
	if !called {
		t.Fatal("expected custom SelectOption to be called")
	}
}
