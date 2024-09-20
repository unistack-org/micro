package rand

import (
	crand "crypto/rand"
	"encoding/binary"
)

// Rand is a wrapper around crypto/rand that adds some convenience functions known from math/rand
type Rand struct {
	buf [8]byte
}

// Int31 function implementation
func (r *Rand) Int31() int32 {
	_, _ = crand.Read(r.buf[:4])
	return int32(binary.BigEndian.Uint32(r.buf[:4]) & ^uint32(1<<31))
}

// Int function implementation
func (r *Rand) Int() int {
	u := uint(r.Int63())
	return int(u << 1 >> 1) // clear sign bit if int == int32
}

// Float64 function implementation
func (r *Rand) Float64() float64 {
again:
	f := float64(r.Int63()) / (1 << 63)
	if f == 1 {
		goto again // resample; this branch is taken O(never)
	}
	return f
}

// Float32 function implementation
func (r *Rand) Float32() float32 {
again:
	f := float32(r.Float64())
	if f == 1 {
		goto again // resample; this branch is taken O(very rarely)
	}
	return f
}

// Uint32 function implementation
func (r *Rand) Uint32() uint32 {
	return uint32(r.Int63() >> 31)
}

// Uint64 function implementation
func (r *Rand) Uint64() uint64 {
	return uint64(r.Int63())>>31 | uint64(r.Int63())<<32
}

// Intn function implementation
func (r *Rand) Intn(n int) int {
	if n <= 1<<31-1 {
		return int(r.Int31n(int32(n)))
	}
	return int(r.Int63n(int64(n)))
}

// Int63 function implementation
func (r *Rand) Int63() int64 {
	_, _ = crand.Read(r.buf[:])
	return int64(binary.BigEndian.Uint64(r.buf[:]) & ^uint64(1<<63))
}

// Int31n function implementation copied from the standard library math/rand implementation of Int31n
func (r *Rand) Int31n(n int32) int32 {
	if n&(n-1) == 0 { // n is power of two, can mask
		return r.Int31() & (n - 1)
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := r.Int31()
	for v > max {
		v = r.Int31()
	}
	return v % n
}

// Int63n function implementation copied from the standard library math/rand implementation of Int63n
func (r *Rand) Int63n(n int64) int64 {
	if n&(n-1) == 0 { // n is power of two, can mask
		return r.Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := r.Int63()
	for v > max {
		v = r.Int63()
	}
	return v % n
}

// Shuffle function implementation copied from the standard library math/rand implementation of Shuffle
func (r *Rand) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(r.Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(r.Int31n(int32(i + 1)))
		swap(i, j)
	}
}
