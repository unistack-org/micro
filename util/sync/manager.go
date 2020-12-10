package sync

import (
	"fmt"
	"time"

	"github.com/unistack-org/micro/v3/store"
)

func (c *syncStore) syncManager() {
	tickerAggregator := make(chan struct{ index int })
	for i, ticker := range c.pendingWriteTickers {
		go func(index int, c chan struct{ index int }, t *time.Ticker) {
			for range t.C {
				c <- struct{ index int }{index: index}
			}
		}(i, tickerAggregator, ticker)
	}
	for {
		select {
		case i := <-tickerAggregator:
			println(i.index, "ticked")
			c.processQueue(i.index)
		}
	}
}

func (c *syncStore) processQueue(index int) {
	c.Lock()
	defer c.Unlock()
	q := c.pendingWrites[index]
	for i := 0; i < q.Len(); i++ {
		r, ok := q.PopFront()
		if !ok {
			panic(fmt.Errorf("retrieved an invalid value from the L%d sync queue", index+1))
		}
		ir, ok := r.(*internalRecord)
		if !ok {
			panic(fmt.Errorf("retrieved a non-internal record from the L%d sync queue", index+1))
		}
		if !ir.expiresAt.IsZero() && time.Now().After(ir.expiresAt) {
			continue
		}
		var opts []store.WriteOption
		if !ir.expiresAt.IsZero() {
			opts = append(opts, store.WriteTTL(ir.expiresAt.Sub(time.Now())))
		}
		// Todo = internal queue also has to hold the corresponding store.WriteOptions
		if err := c.syncOpts.Stores[index+1].Write(c.storeOpts.Context, ir.key, ir.value, opts...); err != nil {
			// some error, so queue for retry and bail
			q.PushBack(ir)
			return
		}
	}
}

func intpow(x, y int64) int64 {
	result := int64(1)
	for 0 != y {
		if 0 != (y & 1) {
			result *= x
		}
		y >>= 1
		x *= x
	}
	return result
}
