package sync

import (
	"context"
	"sync"
)

type WaitGroup struct {
	wg *sync.WaitGroup
	c  int
	mu sync.Mutex
}

func WrapWaitGroup(wg *sync.WaitGroup) *WaitGroup {
	g := &WaitGroup{
		wg: wg,
	}
	return g
}

func NewWaitGroup() *WaitGroup {
	var wg sync.WaitGroup
	return WrapWaitGroup(&wg)
}

func (g *WaitGroup) Add(n int) {
	g.mu.Lock()
	g.c += n
	g.wg.Add(n)
	g.mu.Unlock()
}

func (g *WaitGroup) Done() {
	g.mu.Lock()
	g.c += -1
	g.wg.Add(-1)
	g.mu.Unlock()
}

func (g *WaitGroup) Wait() {
	g.wg.Wait()
}

func (g *WaitGroup) WaitContext(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		g.mu.Lock()
		g.wg.Add(-g.c)
		<-done
		g.wg.Add(g.c)
		g.mu.Unlock()
		return
	case <-done:
		return
	}
}

func (g *WaitGroup) Waiters() int {
	g.mu.Lock()
	c := g.c
	g.mu.Unlock()
	return c
}
