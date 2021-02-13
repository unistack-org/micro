package selector

import (
	"testing"
)

// Tests runs all the tests against a selector to ensure the implementations are consistent
func Tests(t *testing.T, s Selector) {
	r1 := "127.0.0.1:8000"
	r2 := "127.0.0.1:8001"

	t.Run("Select", func(t *testing.T) {
		t.Run("NoRoutes", func(t *testing.T) {
			_, err := s.Select([]string{})
			if err != ErrNoneAvailable {
				t.Fatal("Expected error to be none available")
			}
		})

		t.Run("OneRoute", func(t *testing.T) {
			next, err := s.Select([]string{r1})
			if err != nil {
				t.Fatal("Error should be nil")
			}
			srv := next()
			if r1 != srv {
				t.Fatal("Expected the route to be returned")
			}
		})

		t.Run("MultipleRoutes", func(t *testing.T) {
			next, err := s.Select([]string{r1, r2})
			if err != nil {
				t.Fatal("Error should be nil")
			}
			srv := next()
			if srv != r1 && srv != r2 {
				t.Errorf("Expected the route to be one of the inputs")
			}
		})
	})

	t.Run("Record", func(t *testing.T) {
		if err := s.Record(r1, nil); err != nil {
			t.Fatal("Expected the error to be nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		if s.String() == "" {
			t.Fatal("String returned a blank string")
		}
	})
}
