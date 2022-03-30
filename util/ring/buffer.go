// Package ring provides a simple ring buffer for storing local data
package ring // import "go.unistack.org/micro/v3/util/ring"

import (
	"sync"
	"time"

	"go.unistack.org/micro/v3/util/id"
)

// Buffer is ring buffer
type Buffer struct {
	sync.RWMutex
	streams map[string]*Stream
	vals    []*Entry
	size    int
}

// Entry is ring buffer data entry
type Entry struct {
	Timestamp time.Time
	Value     interface{}
}

// Stream is used to stream the buffer
type Stream struct {
	// Buffered entries
	Entries chan *Entry
	// Stop channel
	Stop chan bool
	// ID of the stream
	ID string
}

// Put adds a new value to ring buffer
func (b *Buffer) Put(v interface{}) {
	b.Lock()
	defer b.Unlock()

	// append to values
	entry := &Entry{
		Value:     v,
		Timestamp: time.Now(),
	}
	b.vals = append(b.vals, entry)

	// trim if bigger than size required
	if len(b.vals) > b.size {
		b.vals = b.vals[1:]
	}

	// send to every stream
	for _, stream := range b.streams {
		select {
		case <-stream.Stop:
			delete(b.streams, stream.ID)
			close(stream.Entries)
		case stream.Entries <- entry:
		}
	}
}

// Get returns the last n entries
func (b *Buffer) Get(n int) []*Entry {
	b.RLock()
	defer b.RUnlock()

	// reset any invalid values
	if n > len(b.vals) || n < 0 {
		n = len(b.vals)
	}

	// create a delta
	delta := len(b.vals) - n

	// return the delta set
	return b.vals[delta:]
}

// Since returns the entries since a specific time
func (b *Buffer) Since(t time.Time) []*Entry {
	b.RLock()
	defer b.RUnlock()

	// return all the values
	if t.IsZero() {
		return b.vals
	}

	// if its in the future return nothing
	if time.Since(t).Seconds() < 0.0 {
		return nil
	}

	for i, v := range b.vals {
		// find the starting point
		d := v.Timestamp.Sub(t)

		// return the values
		if d.Seconds() > 0.0 {
			return b.vals[i:]
		}
	}

	return nil
}

// Stream logs from the buffer
// Close the channel when you want to stop
func (b *Buffer) Stream() (<-chan *Entry, chan bool) {
	b.Lock()
	defer b.Unlock()

	entries := make(chan *Entry, 128)
	id := id.Must()
	stop := make(chan bool)

	b.streams[id] = &Stream{
		ID:      id,
		Entries: entries,
		Stop:    stop,
	}

	return entries, stop
}

// Size returns the size of the ring buffer
func (b *Buffer) Size() int {
	return b.size
}

// New returns a new buffer of the given size
func New(i int) *Buffer {
	return &Buffer{
		size:    i,
		streams: make(map[string]*Stream),
	}
}
