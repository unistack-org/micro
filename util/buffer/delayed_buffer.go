package buffer

import (
	"io"
	"sync"
	"time"
)

var _ io.WriteCloser = (*DelayedBuffer)(nil)

// DelayedBuffer is the buffer that holds items until either the buffer filled or a specified time limit is reached
type DelayedBuffer struct {
	mu        sync.Mutex
	maxWait   time.Duration
	flushTime time.Time
	buffer    chan []byte
	ticker    *time.Ticker
	w         io.Writer
	err       error
}

func NewDelayedBuffer(size int, maxWait time.Duration, w io.Writer) *DelayedBuffer {
	b := &DelayedBuffer{
		buffer:    make(chan []byte, size),
		ticker:    time.NewTicker(maxWait),
		w:         w,
		flushTime: time.Now(),
		maxWait:   maxWait,
	}
	b.loop()
	return b
}

func (b *DelayedBuffer) loop() {
	go func() {
		for range b.ticker.C {
			b.mu.Lock()
			if time.Since(b.flushTime) > b.maxWait {
				b.flush()
			}
			b.mu.Unlock()
		}
	}()
}

func (b *DelayedBuffer) flush() {
	bufLen := len(b.buffer)
	if bufLen > 0 {
		tmp := make([][]byte, bufLen)
		for i := 0; i < bufLen; i++ {
			tmp[i] = <-b.buffer
		}
		for _, t := range tmp {
			_, b.err = b.w.Write(t)
		}
		b.flushTime = time.Now()
	}
}

func (b *DelayedBuffer) Put(items ...[]byte) {
	b.mu.Lock()
	for _, item := range items {
		select {
		case b.buffer <- item:
		default:
			b.flush()
			b.buffer <- item
		}
	}
	b.mu.Unlock()
}

func (b *DelayedBuffer) Close() error {
	b.mu.Lock()
	b.flush()
	close(b.buffer)
	b.ticker.Stop()
	b.mu.Unlock()
	return b.err
}

func (b *DelayedBuffer) Write(data []byte) (int, error) {
	b.Put(data)
	return len(data), b.err
}
