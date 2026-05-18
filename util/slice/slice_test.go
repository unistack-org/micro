package slice

import "testing"

func TestContains(t *testing.T) {
	t.Run("found int", func(t *testing.T) {
		if !Contains([]int{1, 2, 3}, 2) {
			t.Fatal("expected Contains to return true")
		}
	})
	t.Run("not found int", func(t *testing.T) {
		if Contains([]int{1, 2, 3}, 5) {
			t.Fatal("expected Contains to return false")
		}
	})
	t.Run("empty slice", func(t *testing.T) {
		if Contains([]string{}, "x") {
			t.Fatal("expected Contains to return false for empty slice")
		}
	})
	t.Run("found string", func(t *testing.T) {
		if !Contains([]string{"a", "b", "c"}, "b") {
			t.Fatal("expected Contains to return true")
		}
	})
}
