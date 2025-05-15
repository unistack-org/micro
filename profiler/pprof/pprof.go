// Package pprof provides a pprof profiler which writes output to /tmp/[name].{cpu,mem}.pprof
package pprof

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	profile "go.unistack.org/micro/v4/profiler"
)

type profiler struct {
	exit    chan bool
	cpuFile *os.File
	memFile *os.File
	opts    profile.Options
	mu      sync.Mutex
	running bool
}

func (p *profiler) writeHeap(f *os.File) {
	defer f.Close()

	t := time.NewTicker(time.Second * 30)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			runtime.GC()
			_ = pprof.WriteHeapProfile(f)
		case <-p.exit:
			return
		}
	}
}

func (p *profiler) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	// create exit channel
	p.exit = make(chan bool)

	cpuFile := filepath.Join(string(os.PathSeparator)+"tmp", "cpu.pprof")
	memFile := filepath.Join(string(os.PathSeparator)+"tmp", "mem.pprof")

	if len(p.opts.Name) > 0 {
		cpuFile = filepath.Join(string(os.PathSeparator)+"tmp", p.opts.Name+".cpu.pprof")
		memFile = filepath.Join(string(os.PathSeparator)+"tmp", p.opts.Name+".mem.pprof")
	}

	f1, err := os.Create(cpuFile)
	if err != nil {
		return err
	}

	f2, err := os.Create(memFile)
	if err != nil {
		return err
	}

	// start cpu profiling
	if err := pprof.StartCPUProfile(f1); err != nil {
		return err
	}

	// write the heap periodically
	go p.writeHeap(f2)

	// set cpu file
	p.cpuFile = f1
	// set mem file
	p.memFile = f2

	p.running = true

	return nil
}

func (p *profiler) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case <-p.exit:
		return nil
	default:
		close(p.exit)
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		p.running = false
		p.cpuFile = nil
		p.memFile = nil
		return nil
	}
}

func (p *profiler) String() string {
	return "pprof"
}

// NewProfile create new profiler
func NewProfile(opts ...profile.Option) profile.Profiler {
	options := profile.Options{}
	for _, o := range opts {
		o(&options)
	}
	return &profiler{opts: options}
}
