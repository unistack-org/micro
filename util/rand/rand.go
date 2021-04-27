package rand

import (
	"crypto/rand"
	"encoding/binary"
)

// Rand is a wrapper around crypto/rand that adds some convenience functions known from math/rand.
type Rand struct {
	buf [8]byte
}

func (r *Rand) Int31() int32 {
	_, _ = rand.Read(r.buf[:4])
	return int32(binary.BigEndian.Uint32(r.buf[:4]) & ^uint32(1<<31))
}

func (r *Rand) Int() int {
	u := uint(r.Int63())
	return int(u << 1 >> 1) // clear sign bit if int == int32
}

func (r *Rand) Float64() float64 {
again:
	f := float64(r.Int63()) / (1 << 63)
	if f == 1 {
		goto again // resample; this branch is taken O(never)
	}
	return f
}

func (r *Rand) Float32() float32 {
again:
	f := float32(r.Float64())
	if f == 1 {
		goto again // resample; this branch is taken O(very rarely)
	}
	return f
}

func (r *Rand) Uint32() uint32 {
	return uint32(r.Int63() >> 31)
}

func (r *Rand) Uint64() uint64 {
	return uint64(r.Int63())>>31 | uint64(r.Int63())<<32
}

func (r *Rand) Intn(n int) int {
	if n <= 1<<31-1 {
		return int(r.Int31n(int32(n)))
	}
	return int(r.Int63n(int64(n)))
}

func (r *Rand) Int63() int64 {
	_, _ = rand.Read(r.buf[:])
	return int64(binary.BigEndian.Uint64(r.buf[:]) & ^uint64(1<<63))
}

// copied from the standard library math/rand implementation of Int63n
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
