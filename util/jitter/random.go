// Package jitter provides a random jitter
package jitter // import "go.unistack.org/micro/v4/util/jitter"

import (
	"time"

	"go.unistack.org/micro/v4/util/rand"
)

// Random returns a random time to jitter with max cap specified
func Random(d time.Duration) time.Duration {
	var rng rand.Rand
	v := rng.Float64() * float64(d.Nanoseconds())
	return time.Duration(v)
}

func RandomInterval(min, max time.Duration) time.Duration {
	var rng rand.Rand
	return time.Duration(rng.Int63n(max.Nanoseconds()-min.Nanoseconds())+min.Nanoseconds()) * time.Nanosecond
}
