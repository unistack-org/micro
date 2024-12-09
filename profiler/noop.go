package profiler

type noopProfiler struct{}

func (p *noopProfiler) Start() error {
	return nil
}

func (p *noopProfiler) Stop() error {
	return nil
}

func (p *noopProfiler) String() string {
	return "noop"
}

// NewProfiler returns new noop profiler
func NewProfiler(_ ...Option) Profiler {
	return &noopProfiler{}
}
