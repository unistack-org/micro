// Package jitter provides a random jitter
package jitter // import "go.unistack.org/micro/v3/util/jitter"

import (
	"time"

	"go.unistack.org/micro/v3/util/rand"
)

// Random returns a random time to jitter with max cap specified
func Random(d time.Duration) time.Duration {
	var rng rand.Rand
	v := rng.Float64() * float64(d.Nanoseconds())
	return time.Duration(v)
}
