//go:build ignore

package flow

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/heimdalr/dag"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/util/id"
)

type microFlow struct {
	opts Options
}

type microWorkflow struct {
	opts   Options
	g      *dag.DAG
	steps  map[string]Step
	id     string
	status Status
	sync.RWMutex
	init bool
}

func (w *microWorkflow) ID() string {
	return w.id
}

func (w *microWorkflow) Status() Status {
	return w.status
}

func (w *microWorkflow) AppendSteps(steps ...Step) error {
	var err error
	w.Lock()
	defer w.Unlock()

	for _, s := range steps {
		w.steps[s.String()] = s
		if _, err = w.g.AddVertex(s); err != nil {
			return err
		}
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			if err = w.g.AddEdge(src.String(), dst.String()); err != nil {
				return err
			}
		}
	}

	w.g.ReduceTransitively()

	return nil
}

func (w *microWorkflow) RemoveSteps(steps ...Step) error {
	// TODO: handle case when some step requires or required by removed step

	w.Lock()
	defer w.Unlock()

	for _, s := range steps {
		delete(w.steps, s.String())
		w.g.DeleteVertex(s.String())
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			w.g.AddEdge(src.String(), dst.String())
		}
	}

	w.g.ReduceTransitively()

	return nil
}

func (w *microWorkflow) Abort(ctx context.Context, id string) error {
	workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", id))
	return workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusAborted.String())})
}

func (w *microWorkflow) Suspend(ctx context.Context, id string) error {
	workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", id))
	return workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusSuspend.String())})
}

func (w *microWorkflow) Resume(ctx context.Context, id string) error {
	workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", id))
	return workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusRunning.String())})
}

func (w *microWorkflow) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (string, error) {
	w.Lock()
	if !w.init {
		w.g.ReduceTransitively()
		w.init = true
	}
	w.Unlock()

	eid, err := id.New()
	if err != nil {
		return "", err
	}

	//	stepStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("steps", eid))
	workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", eid))

	options := NewExecuteOptions(opts...)

	nopts := make([]ExecuteOption, 0, len(opts)+5)

	nopts = append(nopts,
		ExecuteClient(w.opts.Client),
		ExecuteTracer(w.opts.Tracer),
		ExecuteLogger(w.opts.Logger),
		ExecuteMeter(w.opts.Meter),
	)
	nopts = append(nopts, opts...)

	if werr := workflowStore.Write(ctx, "status", &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
		w.opts.Logger.Error(ctx, "store error: %v", werr)
		return eid, werr
	}

	var startID string
	if options.Start == "" {
		mp := w.g.GetRoots()
		if len(mp) != 1 {
			return eid, ErrStepNotExists
		}
		for k := range mp {
			startID = k
		}
	} else {
		for k, v := range w.g.GetVertices() {
			if v == options.Start {
				startID = k
			}
		}
	}

	if startID == "" {
		return eid, ErrStepNotExists
	}

	if options.Async {
		go w.handleWorkflow(startID, nopts...)
		return eid, nil
	}

	return eid, w.handleWorkflow(startID, nopts...)
}

func (w *microWorkflow) handleWorkflow(startID string, opts ...ExecuteOption) error {
	w.RLock()
	defer w.RUnlock()

	//	stepStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("steps", eid))
	// workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", eid))

	// Get IDs of all descendant vertices.
	flowIDs, errDes := w.g.GetDescendants(startID)
	if errDes != nil {
		return errDes
	}

	// inputChannels provides for input channels for each of the descendant vertices (+ the start-vertex).
	inputChannels := make(map[string]chan FlowResult, len(flowIDs)+1)

	// Iterate vertex IDs and create an input channel for each of them and a single
	// output channel for leaves. Note, this "pre-flight" is needed to ensure we
	// really have an input channel regardless of how we traverse the tree and spawn
	// workers.
	leafCount := 0

	for id := range flowIDs {

		// Get all parents of this vertex.
		parents, errPar := w.g.GetParents(id)
		if errPar != nil {
			return errPar
		}

		// Create a buffered input channel that has capacity for all parent results.
		inputChannels[id] = make(chan FlowResult, len(parents))

		if ok, err := w.g.IsLeaf(id); ok && err == nil {
			leafCount += 1
		}
	}

	// outputChannel caries the results of leaf vertices.
	outputChannel := make(chan FlowResult, leafCount)

	// To also process the start vertex and to have its results being passed to its
	// children, add it to the vertex IDs. Also add an input channel for the start
	// vertex and feed the inputs to this channel.
	flowIDs[startID] = struct{}{}
	inputChannels[startID] = make(chan FlowResult, len(inputs))
	for _, i := range inputs {
		inputChannels[startID] <- i
	}

	wg := sync.WaitGroup{}

	// Iterate all vertex IDs (now incl. start vertex) and handle each worker (incl.
	// inputs and outputs) in a separate goroutine.
	for id := range flowIDs {

		// Get all children of this vertex that later need to be notified. Note, we
		// collect all children before the goroutine to be able to release the read
		// lock as early as possible.
		children, errChildren := w.g.GetChildren(id)
		if errChildren != nil {
			return errChildren
		}

		// Remember to wait for this goroutine.
		wg.Add(1)

		go func(id string) {
			// Get this vertex's input channel.
			// Note, only concurrent read here, which is fine.
			c := inputChannels[id]

			// Await all parent inputs and stuff them into a slice.
			parentCount := cap(c)
			parentResults := make([]FlowResult, parentCount)
			for i := 0; i < parentCount; i++ {
				parentResults[i] = <-c
			}

			// Execute the worker.
			errWorker := callback(w.g, id, parentResults)
			if errWorker != nil {
				return errWorker
			}

			// Send this worker's FlowResult onto all children's input channels or, if it is
			// a leaf (i.e. no children), send the result onto the output channel.
			if len(children) > 0 {
				for child := range children {
					inputChannels[child] <- flowResult
				}
			} else {
				outputChannel <- flowResult
			}

			// "Sign off".
			wg.Done()
		}(id)
	}

	// Wait for all go routines to finish.
	wg.Wait()

	// Await all leaf vertex results and stuff them into a slice.
	resultCount := cap(outputChannel)
	results := make([]FlowResult, resultCount)
	for i := 0; i < resultCount; i++ {
		results[i] = <-outputChannel
	}

	/*
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
						w.opts.Logger.Tracef(nctx, "will be executed %v", steps[idx][nidx])
					}
					cstep := steps[idx][nidx]
					// nolint: nestif
					if len(cstep.Requires()) == 0 {
						wg.Add(1)
						go func(step Step) {
							defer wg.Done()
							if werr := stepStore.Write(ctx, filepath.Join(step.ID(), "req"), req); werr != nil {
								cherr <- werr
								return
							}
							if werr := stepStore.Write(ctx, filepath.Join(step.ID(), "status"), &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
								cherr <- werr
								return
							}
							rsp, serr := step.Execute(nctx, req, nopts...)
							if serr != nil {
								step.SetStatus(StatusFailure)
								if werr := stepStore.Write(ctx, filepath.Join(step.ID(), "rsp"), serr); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
									w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
								}
								if werr := stepStore.Write(ctx, filepath.Join(step.ID(), "status"), &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
									w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
								}
								cherr <- serr
								return
							}
							if werr := stepStore.Write(ctx, filepath.Join(step.ID(), "rsp"), rsp); werr != nil {
								w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
								cherr <- werr
								return
							}
							if werr := stepStore.Write(ctx, filepath.Join(step.ID(), "status"), &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
								w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
								cherr <- werr
								return
							}
						}(cstep)
						wg.Wait()
					} else {
						if werr := stepStore.Write(ctx, filepath.Join(cstep.ID(), "req"), req); werr != nil {
							cherr <- werr
							return
						}
						if werr := stepStore.Write(ctx, filepath.Join(cstep.ID(), "status"), &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
							cherr <- werr
							return
						}
						rsp, serr := cstep.Execute(nctx, req, nopts...)
						if serr != nil {
							cstep.SetStatus(StatusFailure)
							if werr := stepStore.Write(ctx, filepath.Join(cstep.ID(), "rsp"), serr); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
								w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
							}
							if werr := stepStore.Write(ctx, filepath.Join(cstep.ID(), "status"), &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil && w.opts.Logger.V(logger.ErrorLevel) {
								w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
							}
							cherr <- serr
							return
						}
						if werr := stepStore.Write(ctx, filepath.Join(cstep.ID(), "rsp"), rsp); werr != nil {
							w.opts.Logger.Errorf(ctx, "store write error: %v", werr)
							cherr <- werr
							return
						}
						if werr := stepStore.Write(ctx, filepath.Join(cstep.ID(), "status"), &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
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

		logger.Tracef(ctx, "wait for finish or error")
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
				w.opts.Logger.Errorf(w.opts.Context, "store error: %v", werr)
			}
		case err == nil:
			if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
				w.opts.Logger.Errorf(w.opts.Context, "store error: %v", werr)
			}
		case err != nil:
			if werr := workflowStore.Write(w.opts.Context, "status", &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil {
				w.opts.Logger.Errorf(w.opts.Context, "store error: %v", werr)
			}
		}
	*/
	return err
}

// NewFlow create new flow
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
	w := &microWorkflow{opts: f.opts, id: id, g: &dag.DAG{}, steps: make(map[string]Step, len(steps))}

	for _, s := range steps {
		w.steps[s.String()] = s
		w.g.AddVertex(s)
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return nil, ErrStepNotExists
			}
			w.g.AddEdge(src.String(), dst.String())
		}
	}

	w.g.ReduceTransitively()

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

func (s *microCallStep) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error) {
	options := NewExecuteOptions(opts...)
	if options.Client == nil {
		return nil, ErrMissingClient
	}
	rsp := &codec.Frame{}
	copts := []client.CallOption{client.WithRetries(0)}
	if options.Timeout > 0 {
		copts = append(copts,
			client.WithRequestTimeout(options.Timeout),
			client.WithDialTimeout(options.Timeout))
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

func (s *microPublishStep) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error) {
	return nil, nil
}

// NewCallStep create new step with client.Call
func NewCallStep(service string, name string, method string, opts ...StepOption) Step {
	options := NewStepOptions(opts...)
	return &microCallStep{service: service, method: name + "." + method, opts: options}
}

// NewPublishStep create new step with client.Publish
func NewPublishStep(topic string, opts ...StepOption) Step {
	options := NewStepOptions(opts...)
	return &microPublishStep{topic: topic, opts: options}
}
