package pool

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/semconv"
)

var (
	pools   = make([]Statser, 0)
	poolsMu sync.Mutex
)

// Stats struct
type Stats struct {
	Get uint64
	Put uint64
	Mis uint64
	Ret uint64
}

// Statser provides buffer pool stats
type Statser interface {
	Stats() Stats
	Cap() int
}

func init() {
	go newStatsMeter()
}

func newStatsMeter() {
	ticker := time.NewTicker(meter.DefaultMeterStatsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			poolsMu.Lock()
			for _, st := range pools {
				stats := st.Stats()
				meter.DefaultMeter.Counter(semconv.PoolGetTotal, "capacity", strconv.Itoa(st.Cap())).Set(stats.Get)
				meter.DefaultMeter.Counter(semconv.PoolPutTotal, "capacity", strconv.Itoa(st.Cap())).Set(stats.Put)
				meter.DefaultMeter.Counter(semconv.PoolMisTotal, "capacity", strconv.Itoa(st.Cap())).Set(stats.Mis)
				meter.DefaultMeter.Counter(semconv.PoolRetTotal, "capacity", strconv.Itoa(st.Cap())).Set(stats.Ret)
			}
			poolsMu.Unlock()
		}
	}
}

var (
	_ Statser = (*BytePool)(nil)
	_ Statser = (*BytesPool)(nil)
	_ Statser = (*StringsPool)(nil)
)

type Pool[T any] struct {
	p *sync.Pool
}

func (p Pool[T]) Put(t T) {
	p.p.Put(t)
}

func (p Pool[T]) Get() T {
	return p.p.Get().(T)
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

type BytePool struct {
	p   *sync.Pool
	get uint64
	put uint64
	mis uint64
	ret uint64
	c   int
}

func NewBytePool(size int) *BytePool {
	p := &BytePool{c: size}
	p.p = &sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&p.mis, 1)
			b := make([]byte, 0, size)
			return &b
		},
	}
	poolsMu.Lock()
	pools = append(pools, p)
	poolsMu.Unlock()
	return p
}

func (p *BytePool) Cap() int {
	return p.c
}

func (p *BytePool) Stats() Stats {
	return Stats{
		Put: atomic.LoadUint64(&p.put),
		Get: atomic.LoadUint64(&p.get),
		Mis: atomic.LoadUint64(&p.mis),
		Ret: atomic.LoadUint64(&p.ret),
	}
}

func (p *BytePool) Get() *[]byte {
	atomic.AddUint64(&p.get, 1)
	return p.p.Get().(*[]byte)
}

func (p *BytePool) Put(b *[]byte) {
	atomic.AddUint64(&p.put, 1)
	if cap(*b) > p.c {
		atomic.AddUint64(&p.ret, 1)
		return
	}
	*b = (*b)[:0]
	p.p.Put(b)
}

type BytesPool struct {
	p   *sync.Pool
	get uint64
	put uint64
	mis uint64
	ret uint64
	c   int
}

func NewBytesPool(size int) *BytesPool {
	p := &BytesPool{c: size}
	p.p = &sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&p.mis, 1)
			b := bytes.NewBuffer(make([]byte, 0, size))
			return b
		},
	}
	poolsMu.Lock()
	pools = append(pools, p)
	poolsMu.Unlock()
	return p
}

func (p *BytesPool) Cap() int {
	return p.c
}

func (p *BytesPool) Stats() Stats {
	return Stats{
		Put: atomic.LoadUint64(&p.put),
		Get: atomic.LoadUint64(&p.get),
		Mis: atomic.LoadUint64(&p.mis),
		Ret: atomic.LoadUint64(&p.ret),
	}
}

func (p *BytesPool) Get() *bytes.Buffer {
	return p.p.Get().(*bytes.Buffer)
}

func (p *BytesPool) Put(b *bytes.Buffer) {
	if (*b).Cap() > p.c {
		atomic.AddUint64(&p.ret, 1)
		return
	}
	b.Reset()
	p.p.Put(b)
}

type StringsPool struct {
	p   *sync.Pool
	get uint64
	put uint64
	mis uint64
	ret uint64
	c   int
}

func NewStringsPool(size int) *StringsPool {
	p := &StringsPool{c: size}
	p.p = &sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&p.mis, 1)
			return &strings.Builder{}
		},
	}
	poolsMu.Lock()
	pools = append(pools, p)
	poolsMu.Unlock()
	return p
}

func (p *StringsPool) Cap() int {
	return p.c
}

func (p *StringsPool) Stats() Stats {
	return Stats{
		Put: atomic.LoadUint64(&p.put),
		Get: atomic.LoadUint64(&p.get),
		Mis: atomic.LoadUint64(&p.mis),
		Ret: atomic.LoadUint64(&p.ret),
	}
}

func (p *StringsPool) Get() *strings.Builder {
	atomic.AddUint64(&p.get, 1)
	return p.p.Get().(*strings.Builder)
}

func (p *StringsPool) Put(b *strings.Builder) {
	atomic.AddUint64(&p.put, 1)
	if b.Cap() > p.c {
		atomic.AddUint64(&p.ret, 1)
		return
	}
	b.Reset()
	p.p.Put(b)
}
