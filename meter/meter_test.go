package meter

import (
	"bytes"
	"context"
	"testing"
)

func TestNoopMeter_Name(t *testing.T) {
	m := NewMeter(Name("noop"))
	if name := m.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopMeter_Init(t *testing.T) {
	m := NewMeter()
	if err := m.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopMeter_Counter(t *testing.T) {
	m := NewMeter()
	c := m.Counter("test-counter")
	c.Inc()
	c.Dec()
	c.Add(5)
	if got := c.Get(); got != 0 {
		t.Errorf("expected 0 for noop counter, got %d", got)
	}
}

func TestNoopMeter_FloatCounter(t *testing.T) {
	m := NewMeter()
	c := m.FloatCounter("test-float-counter")
	c.Add(5.5)
	c.Sub(2.0)
	if got := c.Get(); got != 0 {
		t.Errorf("expected 0 for noop float counter, got %f", got)
	}
}

func TestNoopMeter_Gauge(t *testing.T) {
	m := NewMeter()
	g := m.Gauge("test-gauge", func() float64 { return 42.0 })
	g.Inc()
	g.Dec()
	if got := g.Get(); got != 0 {
		t.Errorf("expected 0 for noop gauge, got %f", got)
	}
}

func TestNoopMeter_Histogram(t *testing.T) {
	m := NewMeter()
	h := m.Histogram("test-histogram")
	h.Update(1.0)
	h.Reset()
}

func TestNoopMeter_Summary(t *testing.T) {
	m := NewMeter()
	s := m.Summary("test-summary")
	s.Update(1.0)
}

func TestNoopMeter_Write(t *testing.T) {
	m := NewMeter()
	var buf bytes.Buffer
	if err := m.Write(&buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopMeter_Unregister(t *testing.T) {
	m := NewMeter()
	if unregistered := m.Unregister("test"); !unregistered {
		t.Error("expected true for noop unregister")
	}
}

func TestNoopMeter_Set(t *testing.T) {
	m := NewMeter()
	s := m.Set()
	if s == nil {
		t.Error("expected non-nil meter set")
	}
}

func TestNoopMeter_Clone(t *testing.T) {
	m := NewMeter()
	cloned := m.Clone()
	if cloned == nil {
		t.Error("expected non-nil cloned meter")
	}
}

func TestNoopMeter_String(t *testing.T) {
	m := NewMeter()
	if str := m.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopMeter_Options(t *testing.T) {
	m := NewMeter()
	opts := m.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestBuildLabels(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		expected []string
	}{
		{
			name:     "empty",
			labels:   []string{},
			expected: []string{},
		},
		{
			name:     "odd number of labels",
			labels:   []string{"a", "b", "c"},
			expected: []string{"a", "b"},
		},
		{
			name:     "sorted labels",
			labels:   []string{"b", "2", "a", "1"},
			expected: []string{"a", "1", "b", "2"},
		},
		{
			name:     "duplicate labels",
			labels:   []string{"a", "1", "a", "2"},
			expected: []string{"a", "1", "a", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildLabels(tt.labels...)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected length %d, got %d", len(tt.expected), len(got))
			}
			for i := range tt.expected {
				if got[i] != tt.expected[i] {
					t.Errorf("expected %q at index %d, got %q", tt.expected[i], i, got[i])
				}
			}
		})
	}
}

func TestBuildName(t *testing.T) {
	tests := []struct {
		name     string
		metric   string
		labels   []string
		expected string
	}{
		{
			name:     "no labels",
			metric:   "test",
			labels:   []string{},
			expected: "test{}",
		},
		{
			name:     "with labels",
			metric:   "test",
			labels:   []string{"a", "1", "b", "2"},
			expected: `test{a="1",b="2"}`,
		},
		{
			name:     "odd labels",
			metric:   "test",
			labels:   []string{"a", "1", "b"},
			expected: `test{a="1"}`,
		},
		{
			name:     "duplicate labels",
			metric:   "test",
			labels:   []string{"a", "1", "a", "2"},
			expected: `test{a="2"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildName(tt.metric, tt.labels...)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	if opts.Context == nil {
		t.Error("expected non-nil context")
	}
	if opts.Address != DefaultAddress {
		t.Errorf("expected address %q, got %q", DefaultAddress, opts.Address)
	}
	if opts.Path != DefaultPath {
		t.Errorf("expected path %q, got %q", DefaultPath, opts.Path)
	}
}

func TestOptionsContext(t *testing.T) {
	ctx := context.Background()
	opts := NewOptions(Context(ctx))
	if opts.Context != ctx {
		t.Error("expected context to be set")
	}
}

func TestOptionsPath(t *testing.T) {
	opts := NewOptions(Path("/test"))
	if opts.Path != "/test" {
		t.Errorf("expected path /test, got %q", opts.Path)
	}
}

func TestOptionsAddress(t *testing.T) {
	opts := NewOptions(Address(":8080"))
	if opts.Address != ":8080" {
		t.Errorf("expected address :8080, got %q", opts.Address)
	}
}

func TestOptionsName(t *testing.T) {
	opts := NewOptions(Name("test-meter"))
	if opts.Name != "test-meter" {
		t.Errorf("expected name test-meter, got %q", opts.Name)
	}
}

func TestDefaultMeter(t *testing.T) {
	if DefaultMeter == nil {
		t.Error("expected non-nil DefaultMeter")
	}
	if DefaultMeter.String() != "noop" {
		t.Errorf("expected DefaultMeter.String() to be noop, got %q", DefaultMeter.String())
	}
}

func TestMustContext_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustContext with nil meter in context")
		}
	}()
	MustContext(context.Background())
}

func TestMustContext_Valid(t *testing.T) {
	m := NewMeter()
	ctx := NewContext(context.Background(), m)
	got := MustContext(ctx)
	if got != m {
		t.Error("expected meter from context")
	}
}

func TestNoopMeter_SummaryExt(t *testing.T) {
	m := NewMeter()
	s := m.SummaryExt("test-summary-ext", DefaultSummaryWindow, DefaultSummaryQuantiles)
	if s == nil {
		t.Error("expected non-nil summary")
	}
	s.Update(1.0)
}

func TestNoopMeter_HistogramExt(t *testing.T) {
	m := NewMeter()
	h := m.HistogramExt("test-histogram-ext", DefaultSummaryQuantiles)
	if h == nil {
		t.Error("expected non-nil histogram")
	}
	h.Update(1.0)
	h.Reset()
}

func TestNoopMeter_InitWithOpts(t *testing.T) {
	m := NewMeter()
	err := m.Init(Name("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if m.Name() != "test" {
		t.Errorf("expected name test, got %q", m.Name())
	}
}

func TestNoopMeter_CloneWithOpts(t *testing.T) {
	m := NewMeter(Name("original"))
	cloned := m.Clone(Name("cloned"))
	if cloned.Name() != "cloned" {
		t.Errorf("expected cloned name cloned, got %q", cloned.Name())
	}
}

func TestNoopMeter_SetWithOpts(t *testing.T) {
	m := NewMeter(Name("original"))
	s := m.Set(Name("set"))
	if s.Name() != "set" {
		t.Errorf("expected set name set, got %q", s.Name())
	}
}

func TestOptionsQuantiles(t *testing.T) {
	q := []float64{0.1, 0.2}
	opts := NewOptions(Quantiles(q))
	if len(opts.Quantiles) != len(q) {
		t.Fatalf("expected %d quantiles, got %d", len(q), len(opts.Quantiles))
	}
}

func TestOptionsLabels(t *testing.T) {
	opts := NewOptions(Labels("a", "1", "b", "2"))
	if len(opts.Labels) != 4 {
		t.Fatalf("expected 4 labels, got %d", len(opts.Labels))
	}
}

func TestOptionsWriteProcessMetrics(t *testing.T) {
	opts := NewOptions(WriteProcessMetrics(true))
	if !opts.WriteProcessMetrics {
		t.Error("expected WriteProcessMetrics to be true")
	}
}

func TestOptionsWriteFDMetrics(t *testing.T) {
	opts := NewOptions(WriteFDMetrics(true))
	if !opts.WriteFDMetrics {
		t.Error("expected WriteFDMetrics to be true")
	}
}
