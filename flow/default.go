package flow

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/silas/dag"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/codec"
)

type microFlow struct {
	opts Options
}

type microWorkflow struct {
	id   string
	g    *dag.AcyclicGraph
	init bool
	sync.RWMutex
	opts  Options
	steps map[string]Step
}

func (w *microWorkflow) ID() string {
	return w.id
}

func (w *microWorkflow) Steps() [][]Step {
	return nil
}

func (w *microWorkflow) AppendSteps(ctx context.Context, steps ...Step) error {
	return nil
}

func (w *microWorkflow) RemoveSteps(ctx context.Context, steps ...Step) error {
	return nil
}

func (w *microWorkflow) Execute(ctx context.Context, req interface{}, opts ...ExecuteOption) (string, error) {
	w.Lock()
	if !w.init {
		if err := w.g.Validate(); err != nil {
			w.Unlock()
			return "", err
		}
		w.g.TransitiveReduction()
		w.init = true
	}
	w.Unlock()

	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	options := NewExecuteOptions(opts...)
	var steps [][]Step
	fn := func(n dag.Vertex, idx int) error {
		if idx == 0 {
			steps = make([][]Step, 1)
			steps[0] = make([]Step, 0, 1)
		} else if idx >= len(steps) {
			tsteps := make([][]Step, idx+1)
			copy(tsteps, steps)
			steps = tsteps
			steps[idx] = make([]Step, 0, 1)
		}
		steps[idx] = append(steps[idx], n.(Step))
		return nil
	}

	var root dag.Vertex
	if options.Start != "" {
		var ok bool
		w.RLock()
		root, ok = w.steps[options.Start]
		w.RUnlock()
		if !ok {
			return "", ErrStepNotExists
		}
	} else {
		root, err = w.g.Root()
		if err != nil {
			return "", err
		}
	}
	if options.Reverse {
		err = w.g.SortedReverseDepthFirstWalk([]dag.Vertex{root}, fn)
	} else {
		err = w.g.SortedDepthFirstWalk([]dag.Vertex{root}, fn)
	}
	if err != nil {
		return "", err
	}

	var wg sync.WaitGroup
	cherr := make(chan error, 1)
	defer close(cherr)

	nctx, cancel := context.WithCancel(ctx)
	defer cancel()
	nopts := make([]ExecuteOption, 0, len(opts)+5)
	nopts = append(nopts, ExecuteClient(w.opts.Client), ExecuteTracer(w.opts.Tracer), ExecuteLogger(w.opts.Logger), ExecuteMeter(w.opts.Meter), ExecuteStore(w.opts.Store))

	go func() {
		for idx := range steps {
			wg.Add(len(steps[idx]))
			for nidx := range steps[idx] {
				go func(step Step) {
					defer wg.Done()
					if err = step.Execute(nctx, req, nopts...); err != nil {
						cherr <- err
						cancel()
					}
				}(steps[idx][nidx])
			}
			wg.Wait()
		}
		cherr <- nil
	}()

	err = <-cherr

	return uid.String(), err
}

func NewFlow(opts ...Option) Flow {
	options := NewOptions(opts...)
	return &microFlow{opts: options}
}

func (f *microFlow) Options() Options {
	return f.opts
}

func (f *microFlow) Init(opts ...Option) error {
	for _, o := range opts {
		o(&f.opts)
	}
	if err := f.opts.Client.Init(); err != nil {
		return err
	}
	if err := f.opts.Tracer.Init(); err != nil {
		return err
	}
	if err := f.opts.Logger.Init(); err != nil {
		return err
	}
	if err := f.opts.Meter.Init(); err != nil {
		return err
	}
	if err := f.opts.Store.Init(); err != nil {
		return err
	}
	return nil
}

func (f *microFlow) WorkflowList(ctx context.Context) ([]Workflow, error) {
	return nil, nil
}

func (f *microFlow) WorkflowCreate(ctx context.Context, id string, steps ...Step) (Workflow, error) {
	w := &microWorkflow{opts: f.opts, id: id, g: &dag.AcyclicGraph{}, steps: make(map[string]Step, len(steps))}

	for _, s := range steps {
		w.steps[s.String()] = s
		w.g.Add(s)
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return nil, ErrStepNotExists
			}
			w.g.Connect(dag.BasicEdge(src, dst))
		}
	}

	if err := w.g.Validate(); err != nil {
		return nil, err
	}
	w.g.TransitiveReduction()

	w.init = true

	return w, nil
}

func (f *microFlow) WorkflowRemove(ctx context.Context, id string) error {
	return nil
}

func (f *microFlow) WorkflowSave(ctx context.Context, w Workflow) error {
	return nil
}

func (f *microFlow) WorkflowLoad(ctx context.Context, id string) (Workflow, error) {
	return nil, nil
}

type microCallStep struct {
	opts    StepOptions
	service string
	method  string
}

func (s *microCallStep) ID() string {
	return s.String()
}

func (s *microCallStep) Options() StepOptions {
	return s.opts
}

func (s *microCallStep) Endpoint() string {
	return s.method
}

func (s *microCallStep) Requires() []string {
	return s.opts.Requires
}

func (s *microCallStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.opts.Requires = append(s.opts.Requires, step.String())
	}
	return nil
}

func (s *microCallStep) String() string {
	if s.opts.ID != "" {
		return s.opts.ID
	}
	return fmt.Sprintf("%s.%s", s.service, s.method)
}

func (s *microCallStep) Name() string {
	return s.String()
}

func (s *microCallStep) Hashcode() interface{} {
	return s.String()
}

func (s *microCallStep) Execute(ctx context.Context, req interface{}, opts ...ExecuteOption) error {
	options := NewExecuteOptions(opts...)
	if options.Client == nil {
		return fmt.Errorf("client not set")
	}
	rsp := &codec.Frame{}
	copts := []client.CallOption{client.WithRetries(0)}
	if options.Timeout > 0 {
		copts = append(copts, client.WithRequestTimeout(options.Timeout), client.WithDialTimeout(options.Timeout))
	}
	err := options.Client.Call(ctx, options.Client.NewRequest(s.service, s.method, req), rsp)
	return err
}

type microPublishStep struct {
	opts  StepOptions
	topic string
}

func (s *microPublishStep) ID() string {
	return s.String()
}

func (s *microPublishStep) Options() StepOptions {
	return s.opts
}

func (s *microPublishStep) Endpoint() string {
	return s.topic
}

func (s *microPublishStep) Requires() []string {
	return s.opts.Requires
}

func (s *microPublishStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.opts.Requires = append(s.opts.Requires, step.String())
	}
	return nil
}

func (s *microPublishStep) String() string {
	if s.opts.ID != "" {
		return s.opts.ID
	}
	return fmt.Sprintf("%s", s.topic)
}

func (s *microPublishStep) Name() string {
	return s.String()
}

func (s *microPublishStep) Hashcode() interface{} {
	return s.String()
}

func (s *microPublishStep) Execute(ctx context.Context, req interface{}, opts ...ExecuteOption) error {
	return nil
}

func NewCallStep(service string, method string, opts ...StepOption) Step {
	options := NewStepOptions(opts...)
	return &microCallStep{service: service, method: method, opts: options}
}

func NewPublishStep(topic string, opts ...StepOption) Step {
	options := NewStepOptions(opts...)
	return &microPublishStep{topic: topic, opts: options}
}
