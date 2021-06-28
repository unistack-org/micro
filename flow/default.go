package flow

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/silas/dag"
)

type microFlow struct {
	opts Options
}

type microWorkflow struct {
	id   string
	g    *dag.AcyclicGraph
	init bool
	sync.RWMutex
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

	var steps [][]string
	fn := func(n dag.Vertex, idx int) error {
		if idx == 0 {
			steps = make([][]string, 1)
			steps[0] = make([]string, 0, 1)
		} else if idx >= len(steps) {
			tsteps := make([][]string, idx+1)
			copy(tsteps, steps)
			steps = tsteps
			steps[idx] = make([]string, 0, 1)
		}
		steps[idx] = append(steps[idx], fmt.Sprintf("%s", n))
		return nil
	}

	w.RLock()
	err = w.g.SortedDepthFirstWalk([]dag.Vertex{start}, fn)
	w.RUnlock()

	if err != nil {
		return "", err
	}

	return uid.String(), nil
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
	w := &microWorkflow{id: id, g: &dag.AcyclicGraph{}}

	for _, s := range steps {
		w.g.Add(s.Options().ID)
	}
	for _, s := range steps {
		for _, req := range s.Requires() {
			w.g.Connect(dag.BasicEdge(s, req))
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
	opts     StepOptions
	service  string
	method   string
	requires []Step
}

func (s *microCallStep) ID() string {
	return s.opts.ID
}

func (s *microCallStep) Options() StepOptions {
	return s.opts
}

func (s *microCallStep) Endpoint() string {
	return s.method
}

func (s *microCallStep) Requires() []Step {
	return s.requires
}

func (s *microCallStep) Require(steps ...Step) error {
	s.requires = append(s.requires, steps...)
	return nil
}

func (s *microCallStep) Execute(ctx context.Context, req interface{}, opts ...ExecuteOption) error {
	return nil
}

type microPublishStep struct {
	opts     StepOptions
	topic    string
	requires []Step
}

func (s *microPublishStep) ID() string {
	return s.opts.ID
}

func (s *microPublishStep) Options() StepOptions {
	return s.opts
}

func (s *microPublishStep) Endpoint() string {
	return s.topic
}

func (s *microPublishStep) Requires() []Step {
	return s.requires
}

func (s *microPublishStep) Require(steps ...Step) error {
	s.requires = append(s.requires, steps...)
	return nil
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
