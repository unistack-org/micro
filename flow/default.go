package flow

import (
	"context"
	"fmt"
	"sync"

	"github.com/silas/dag"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
	moptions "go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/util/id"
)

type microFlow struct {
	opts Options
}

type microWorkflow struct {
	opts   Options
	g      *dag.AcyclicGraph
	steps  map[string]Step
	id     string
	status Status
	sync.RWMutex
	init bool
}

func (w *microWorkflow) ID() string {
	return w.id
}

func (w *microWorkflow) Steps() ([][]Step, error) {
	return w.getSteps("", false)
}

func (w *microWorkflow) Status() Status {
	return w.status
}

func (w *microWorkflow) AppendSteps(steps ...Step) error {
	w.Lock()

	for _, s := range steps {
		w.steps[s.String()] = s
		w.g.Add(s)
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			w.g.Connect(dag.BasicEdge(src, dst))
		}
	}

	if err := w.g.Validate(); err != nil {
		w.Unlock()
		return err
	}

	w.g.TransitiveReduction()

	w.Unlock()

	return nil
}

func (w *microWorkflow) RemoveSteps(steps ...Step) error {
	// TODO: handle case when some step requires or required by removed step

	w.Lock()

	for _, s := range steps {
		delete(w.steps, s.String())
		w.g.Remove(s)
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			w.g.Connect(dag.BasicEdge(src, dst))
		}
	}

	if err := w.g.Validate(); err != nil {
		w.Unlock()
		return err
	}

	w.g.TransitiveReduction()

	w.Unlock()

	return nil
}

func (w *microWorkflow) getSteps(start string, reverse bool) ([][]Step, error) {
	var steps [][]Step
	var root dag.Vertex
	var err error

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

	if start != "" {
		var ok bool
		w.RLock()
		root, ok = w.steps[start]
		w.RUnlock()
		if !ok {
			return nil, ErrStepNotExists
		}
	} else {
		root, err = w.g.Root()
		if err != nil {
			return nil, err
		}
	}

	if reverse {
		err = w.g.SortedReverseDepthFirstWalk([]dag.Vertex{root}, fn)
	} else {
		err = w.g.SortedDepthFirstWalk([]dag.Vertex{root}, fn)
	}
	if err != nil {
		return nil, err
	}

	return steps, nil
}

func (w *microWorkflow) Abort(ctx context.Context, id string) error {
	workflowStore := store.NewNamespaceStore(w.opts.Store, "workflows"+w.opts.Store.Options().Separator+id)
	return workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusAborted.String())})
}

func (w *microWorkflow) Suspend(ctx context.Context, id string) error {
	workflowStore := store.NewNamespaceStore(w.opts.Store, "workflows"+w.opts.Store.Options().Separator+id)
	return workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusSuspend.String())})
}

func (w *microWorkflow) Resume(ctx context.Context, id string) error {
	workflowStore := store.NewNamespaceStore(w.opts.Store, "workflows"+w.opts.Store.Options().Separator+id)
	return workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusRunning.String())})
}

func (w *microWorkflow) Execute(ctx context.Context, req *Message, opts ...options.Option) (string, error) {
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

	eid, err := id.New()
	if err != nil {
		return "", err
	}

	stepStore := store.NewNamespaceStore(w.opts.Store, "steps"+w.opts.Store.Options().Separator+eid)
	workflowStore := store.NewNamespaceStore(w.opts.Store, "workflows"+w.opts.Store.Options().Separator+eid)

	options := NewExecuteOptions(opts...)

	steps, err := w.getSteps(options.Start, options.Reverse)
	if err != nil {
		if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusPending.String())}); werr != nil {
			w.opts.Logger.Error(w.opts.Context, "store write error", "error", werr.Error())
		}
		return "", err
	}

	var wg sync.WaitGroup
	cherr := make(chan error, 1)
	chstatus := make(chan Status, 1)

	nctx, cancel := context.WithCancel(ctx)
	defer cancel()

	nopts := make([]moptions.Option, 0, len(opts)+5)

	nopts = append(nopts,
		moptions.Client(w.opts.Client),
		moptions.Tracer(w.opts.Tracer),
		moptions.Logger(w.opts.Logger),
		moptions.Meter(w.opts.Meter),
	)
	nopts = append(nopts, opts...)
	done := make(chan struct{})

	if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
		w.opts.Logger.Error(w.opts.Context, "store write error", "error", werr.Error())
		return eid, werr
	}
	for idx := range steps {
		for nidx := range steps[idx] {
			cstep := steps[idx][nidx]
			if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusPending.String())}); werr != nil {
				return eid, werr
			}
		}
	}

	go func() {
		for idx := range steps {
			for nidx := range steps[idx] {
				wStatus := &codec.Frame{}
				if werr := workflowStore.Read(w.opts.Context, "status", wStatus); werr != nil {
					cherr <- werr
					return
				}
				if status := StringStatus[string(wStatus.Data)]; status != StatusRunning {
					chstatus <- status
					return
				}
				if w.opts.Logger.V(logger.TraceLevel) {
					w.opts.Logger.Trace(nctx, fmt.Sprintf("step will be executed %v", steps[idx][nidx]))
				}
				cstep := steps[idx][nidx]
				// nolint: nestif
				if len(cstep.Requires()) == 0 {
					wg.Add(1)
					go func(step Step) {
						defer wg.Done()
						if werr := stepStore.Write(ctx, step.ID()+w.opts.Store.Options().Separator+"req", req); werr != nil {
							cherr <- werr
							return
						}
						if werr := stepStore.Write(ctx, step.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
							cherr <- werr
							return
						}
						rsp, serr := step.Execute(nctx, req, nopts...)
						if serr != nil {
							step.SetStatus(StatusFailure)
							if werr := stepStore.Write(ctx, step.ID()+w.opts.Store.Options().Separator+"rsp", serr); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
								w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
							}
							if werr := stepStore.Write(ctx, step.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
								w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
							}
							cherr <- serr
							return
						}
						if werr := stepStore.Write(ctx, step.ID()+w.opts.Store.Options().Separator+"rsp", rsp); werr != nil {
							w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
							cherr <- werr
							return
						}
						if werr := stepStore.Write(ctx, step.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
							w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
							cherr <- werr
							return
						}
					}(cstep)
					wg.Wait()
				} else {
					if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"req", req); werr != nil {
						cherr <- werr
						return
					}
					if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
						cherr <- werr
						return
					}
					rsp, serr := cstep.Execute(nctx, req, nopts...)
					if serr != nil {
						cstep.SetStatus(StatusFailure)
						if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"rsp", serr); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
							w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
						}
						if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
							w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
						}
						cherr <- serr
						return
					}
					if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"rsp", rsp); werr != nil {
						w.opts.Logger.Error(ctx, "store write error", "error", werr.Error())
						cherr <- werr
						return
					}
					if werr := stepStore.Write(ctx, cstep.ID()+w.opts.Store.Options().Separator+"status", &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
						cherr <- werr
						return
					}
				}
			}
		}
		close(done)
	}()

	if options.Async {
		return eid, nil
	}

	w.opts.Logger.Trace(ctx, "wait for finish or error")
	select {
	case <-nctx.Done():
		err = nctx.Err()
	case cerr := <-cherr:
		err = cerr
	case <-done:
		close(cherr)
	case <-chstatus:
		close(chstatus)
		return eid, nil
	}

	switch {
	case nctx.Err() != nil:
		if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusAborted.String())}); werr != nil {
			w.opts.Logger.Error(w.opts.Context, "store write error", "error", werr.Error())
		}
	case err == nil:
		if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
			w.opts.Logger.Error(w.opts.Context, "store write error", "error", werr.Error())
		}
	case err != nil:
		if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil {
			w.opts.Logger.Error(w.opts.Context, "store write error", "error", werr.Error())
		}
	}

	return eid, err
}

// NewFlow create new flow
func NewFlow(opts ...options.Option) Flow {
	options := NewOptions(opts...)
	return &microFlow{opts: options}
}

func (f *microFlow) Options() Options {
	return f.opts
}

func (f *microFlow) Init(opts ...options.Option) error {
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
	rsp     *Message
	req     *Message
	service string
	method  string
	opts    StepOptions
	status  Status
}

func (s *microCallStep) Request() *Message {
	return s.req
}

func (s *microCallStep) Response() *Message {
	return s.rsp
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

func (s *microCallStep) GetStatus() Status {
	return s.status
}

func (s *microCallStep) SetStatus(status Status) {
	s.status = status
}

func (s *microCallStep) Execute(ctx context.Context, req *Message, opts ...options.Option) (*Message, error) {
	options := NewExecuteOptions(opts...)
	if options.Client == nil {
		return nil, ErrMissingClient
	}
	rsp := &codec.Frame{}
	copts := []moptions.Option{client.Retries(0)}
	if options.Timeout > 0 {
		copts = append(copts,
			client.RequestTimeout(options.Timeout),
			client.DialTimeout(options.Timeout))
	}
	nctx := metadata.NewOutgoingContext(ctx, req.Header)
	err := options.Client.Call(nctx, options.Client.NewRequest(s.service, s.method, &codec.Frame{Data: req.Body}), rsp, copts...)
	if err != nil {
		return nil, err
	}
	md, _ := metadata.FromOutgoingContext(nctx)
	return &Message{Header: md, Body: rsp.Data}, err
}

type microPublishStep struct {
	req    *Message
	rsp    *Message
	topic  string
	opts   StepOptions
	status Status
}

func (s *microPublishStep) Request() *Message {
	return s.req
}

func (s *microPublishStep) Response() *Message {
	return s.rsp
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
	return s.topic
}

func (s *microPublishStep) Name() string {
	return s.String()
}

func (s *microPublishStep) Hashcode() interface{} {
	return s.String()
}

func (s *microPublishStep) GetStatus() Status {
	return s.status
}

func (s *microPublishStep) SetStatus(status Status) {
	s.status = status
}

func (s *microPublishStep) Execute(ctx context.Context, req *Message, opts ...options.Option) (*Message, error) {
	return nil, nil
}

// NewCallStep create new step with client.Call
func NewCallStep(service string, name string, method string, opts ...options.Option) Step {
	options := NewStepOptions(opts...)
	return &microCallStep{service: service, method: name + "." + method, opts: options}
}

// NewPublishStep create new step with client.Publish
func NewPublishStep(topic string, opts ...options.Option) Step {
	options := NewStepOptions(opts...)
	return &microPublishStep{topic: topic, opts: options}
}
