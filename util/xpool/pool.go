package pool

import (
	"bytes"
	"sync"
)

type Pool[T any] struct {
	p *sync.Pool
}

func NewPool[T any](fn func() T) Pool[T] {
	return Pool[T]{
		p: &sync.Pool{
			New: func() interface{} {
				return fn()
			},
		},
	}
}

func (p Pool[T]) Get() T {
	return p.p.Get().(T)
}

func (p Pool[T]) Put(t T) {
	p.p.Put(t)
}

func NewBytePool(size int) Pool[T] {
	return NewPool(func() []byte { return make([]byte, size) })
}

func NewBytesPool() Pool[T] {
	return NewPool(func() *bytes.Buffer { return bytes.NewBuffer(nil) })
}
