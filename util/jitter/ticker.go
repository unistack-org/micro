package jitter

import (
	"time"

	"go.unistack.org/micro/v3/util/rand"
)

// Ticker is similar to time.Ticker but ticks at random intervals between
// the min and max duration values (stored internally as int64 nanosecond
// counts).
type Ticker struct {
	C    chan time.Time
	done chan chan struct{}
	min  int64
	max  int64
	rng  rand.Rand
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
	}
	go ticker.run()
	return ticker
}

// Stop terminates the ticker goroutine and closes the C channel.
func (ticker *Ticker) Stop() {
	c := make(chan struct{})
	ticker.done <- c
	<-c
}

func (ticker *Ticker) run() {
	defer close(ticker.C)
	t := time.NewTimer(ticker.nextInterval())
	for {
		// either a stop signal or a timeout
		select {
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
