package pool

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/semconv"
)

func unregisterMetrics(size int) {
	meter.DefaultMeter.Unregister(semconv.PoolGetTotal, "capacity", strconv.Itoa(size))
	meter.DefaultMeter.Unregister(semconv.PoolPutTotal, "capacity", strconv.Itoa(size))
	meter.DefaultMeter.Unregister(semconv.PoolMisTotal, "capacity", strconv.Itoa(size))
	meter.DefaultMeter.Unregister(semconv.PoolRetTotal, "capacity", strconv.Itoa(size))
}

type Stats struct {
	Get uint64
	Put uint64
	Mis uint64
	Ret uint64
}

type Pool[T any] struct {
	p   *sync.Pool
	get *atomic.Uint64
	put *atomic.Uint64
	mis *atomic.Uint64
	ret *atomic.Uint64
	c   int
}

func (p Pool[T]) Put(t T) {
	p.p.Put(t)
}

func (p Pool[T]) Get() T {
	return p.p.Get().(T)
}

func NewPool[T any](fn func() T, size int) Pool[T] {
	p := Pool[T]{
		c:   size,
		get: &atomic.Uint64{},
		put: &atomic.Uint64{},
		mis: &atomic.Uint64{},
		ret: &atomic.Uint64{},
	}

	p.p = &sync.Pool{
		New: func() any {
			p.mis.Add(1)
			return fn()
		},
	}

	meter.DefaultMeter.Gauge(semconv.PoolGetTotal, func() float64 {
		return float64(p.get.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolPutTotal, func() float64 {
		return float64(p.put.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolMisTotal, func() float64 {
		return float64(p.mis.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolRetTotal, func() float64 {
		return float64(p.ret.Load())
	}, "capacity", strconv.Itoa(p.c))

	return p
}

type BytePool struct {
	p   *sync.Pool
	get *atomic.Uint64
	put *atomic.Uint64
	mis *atomic.Uint64
	ret *atomic.Uint64
	c   int
}

func NewBytePool(size int) *BytePool {
	p := &BytePool{
		c:   size,
		get: &atomic.Uint64{},
		put: &atomic.Uint64{},
		mis: &atomic.Uint64{},
		ret: &atomic.Uint64{},
	}
	p.p = &sync.Pool{
		New: func() any {
			p.mis.Add(1)
			b := make([]byte, 0, size)
			return &b
		},
	}

	meter.DefaultMeter.Gauge(semconv.PoolGetTotal, func() float64 {
		return float64(p.get.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolPutTotal, func() float64 {
		return float64(p.put.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolMisTotal, func() float64 {
		return float64(p.mis.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolRetTotal, func() float64 {
		return float64(p.ret.Load())
	}, "capacity", strconv.Itoa(p.c))

	return p
}

func (p *BytePool) Cap() int {
	return p.c
}

func (p *BytePool) Stats() Stats {
	return Stats{
		Put: p.put.Load(),
		Get: p.get.Load(),
		Mis: p.mis.Load(),
		Ret: p.ret.Load(),
	}
}

func (p *BytePool) Get() *[]byte {
	p.get.Add(1)
	return p.p.Get().(*[]byte)
}

func (p *BytePool) Put(b *[]byte) {
	p.put.Add(1)
	if cap(*b) > p.c {
		p.ret.Add(1)
		return
	}
	*b = (*b)[:0]
	p.p.Put(b)
}

func (p *BytePool) Close() {
	unregisterMetrics(p.c)
}

type BytesPool struct {
	p   *sync.Pool
	get *atomic.Uint64
	put *atomic.Uint64
	mis *atomic.Uint64
	ret *atomic.Uint64
	c   int
}

func NewBytesPool(size int) *BytesPool {
	p := &BytesPool{
		c:   size,
		get: &atomic.Uint64{},
		put: &atomic.Uint64{},
		mis: &atomic.Uint64{},
		ret: &atomic.Uint64{},
	}
	p.p = &sync.Pool{
		New: func() any {
			p.mis.Add(1)
			b := bytes.NewBuffer(make([]byte, 0, size))
			return b
		},
	}

	meter.DefaultMeter.Gauge(semconv.PoolGetTotal, func() float64 {
		return float64(p.get.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolPutTotal, func() float64 {
		return float64(p.put.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolMisTotal, func() float64 {
		return float64(p.mis.Load())
	}, "capacity", strconv.Itoa(p.c))

	meter.DefaultMeter.Gauge(semconv.PoolRetTotal, func() float64 {
		return float64(p.ret.Load())
	}, "capacity", strconv.Itoa(p.c))

	return p
}

func (p *BytesPool) Cap() int {
	return p.c
}

func (p *BytesPool) Stats() Stats {
	return Stats{
		Put: p.put.Load(),
		Get: p.get.Load(),
		Mis: p.mis.Load(),
		Ret: p.ret.Load(),
	}
}

func (p *BytesPool) Get() *bytes.Buffer {
	return p.p.Get().(*bytes.Buffer)
}

func (p *BytesPool) Put(b *bytes.Buffer) {
	p.put.Add(1)
	if (*b).Cap() > p.c {
		p.ret.Add(1)
		return
	}
	b.Reset()
	p.p.Put(b)
}

func (p *BytesPool) Close() {
	unregisterMetrics(p.c)
}

type StringsPool struct {
	p   *sync.Pool
	get *atomic.Uint64
	put *atomic.Uint64
	mis *atomic.Uint64
	ret *atomic.Uint64
	c   int
}

func NewStringsPool(size int) *StringsPool {
	p := &StringsPool{
		c:   size,
		get: &atomic.Uint64{},
		put: &atomic.Uint64{},
		mis: &atomic.Uint64{},
		ret: &atomic.Uint64{},
	}
	p.p = &sync.Pool{
		New: func() any {
			p.mis.Add(1)
			return &strings.Builder{}
		},
	}

	return p
}

func (p *StringsPool) Cap() int {
	return p.c
}

func (p *StringsPool) Stats() Stats {
	return Stats{
		Put: p.put.Load(),
		Get: p.get.Load(),
		Mis: p.mis.Load(),
		Ret: p.ret.Load(),
	}
}

func (p *StringsPool) Get() *strings.Builder {
	p.get.Add(1)
	return p.p.Get().(*strings.Builder)
}

func (p *StringsPool) Put(b *strings.Builder) {
	p.put.Add(1)
	if b.Cap() > p.c {
		p.ret.Add(1)
		return
	}
	b.Reset()
	p.p.Put(b)
}

func (p *StringsPool) Close() {
	unregisterMetrics(p.c)
}
