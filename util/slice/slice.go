// Package slice provides utilities for working with slices.
package slice

// Contains reports whether v is present in s.
func Contains[T comparable](s []T, v T) bool {
	for _, vs := range s {
		if vs == v {
			return true
		}
	}
	return false
}
