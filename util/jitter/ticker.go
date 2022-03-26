package jitter // import "go.unistack.org/micro/v3/util/jitter"

import (
	"context"
	"time"

	"go.unistack.org/micro/v3/util/rand"
)

// Ticker is similar to time.Ticker but ticks at random intervals between
// the min and max duration values (stored internally as int64 nanosecond
// counts).
type Ticker struct {
	ctx  context.Context
	done chan chan struct{}
	C    chan time.Time
	min  int64
	max  int64
	exp  int64
	exit bool
	rng  rand.Rand
}

// NewTickerContext returns a pointer to an initialized instance of the Ticker.
// It works like NewTicker except that it has ability to close via context.
// Also it works fine with context.WithTimeout to handle max time to run ticker.
func NewTickerContext(ctx context.Context, min, max time.Duration) *Ticker {
	ticker := &Ticker{
		C:    make(chan time.Time),
		done: make(chan chan struct{}),
		min:  min.Nanoseconds(),
		max:  max.Nanoseconds(),
		ctx:  ctx,
	}
	go ticker.run()
	return ticker
}

// NewTicker returns a pointer to an initialized instance of the Ticker.
// Min and max are durations of the shortest and longest allowed
// ticks. Ticker will run in a goroutine until explicitly stopped.
func NewTicker(min, max time.Duration) *Ticker {
	ticker := &Ticker{
		C:    make(chan time.Time),
		done: make(chan chan struct{}),
		min:  min.Nanoseconds(),
		max:  max.Nanoseconds(),
		ctx:  context.Background(),
	}
	go ticker.run()
	return ticker
}

// Stop terminates the ticker goroutine and closes the C channel.
func (ticker *Ticker) Stop() {
	if ticker.exit {
		return
	}
	c := make(chan struct{})
	ticker.done <- c
	<-c
	// close(ticker.C)
	ticker.exit = true
}

func (ticker *Ticker) run() {
	defer close(ticker.C)
	t := time.NewTimer(ticker.nextInterval())
	for {
		// either a stop signal or a timeout
		select {
		case <-ticker.ctx.Done():
			t.Stop()
		case c := <-ticker.done:
			t.Stop()
			close(c)
			return
		case <-t.C:
			select {
			case ticker.C <- time.Now():
				t.Stop()
				t = time.NewTimer(ticker.nextInterval())
			default:
				// there could be noone receiving...
			}
		}
	}
}

func (ticker *Ticker) nextInterval() time.Duration {
	return time.Duration(ticker.rng.Int63n(ticker.max-ticker.min)+ticker.min) * time.Nanosecond
}
