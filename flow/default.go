package flow

import (
	"context"
	"encoding/json"
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
		if err := w.g.DeleteVertex(s.String()); err != nil {
			return fmt.Errorf("failed to delete vertex %s: %w", s.String(), err)
		}
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			if err := w.g.AddEdge(src.String(), dst.String()); err != nil {
				return fmt.Errorf("failed to add edge %s -> %s: %w", src.String(), dst.String(), err)
			}
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
		go func() {
			if err := w.handleWorkflow(startID, nopts...); err != nil {
				w.opts.Logger.Error(context.Background(), "async workflow execution failed", "error", err, "workflow_id", eid)
			}
		}()
		return eid, nil
	}

	return eid, w.handleWorkflow(startID, nopts...)
}

func (w *microWorkflow) handleWorkflow(startID string, opts ...ExecuteOption) error {
	options := NewExecuteOptions(opts...)
	
	// Создаем store для workflow
	eid := w.id
	workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", eid))
	stepStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("steps", eid))

	// Трек успешных шагов для компенсации
	executedSteps := make([]string, 0)
	var execMu sync.Mutex
	
	// Завершаем workflow с ошибкой и выполняем компенсацию
	failWorkflow := func(err error) error {
		w.opts.Logger.Error(options.Context, "workflow failed: %v", err)
		
		// Обновляем статус workflow
		if werr := workflowStore.Write(options.Context, "status", &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil {
			w.opts.Logger.Error(options.Context, "store error: %v", werr)
		}
		
		// Выполняем компенсацию в обратном порядке
		execMu.Lock()
		stepsToCompensate := make([]string, len(executedSteps))
		copy(stepsToCompensate, executedSteps)
		execMu.Unlock()
		
		for i := len(stepsToCompensate) - 1; i >= 0; i-- {
			stepID := stepsToCompensate[i]
			step, ok := w.steps[stepID]
			if !ok {
				continue
			}
			
			w.opts.Logger.Info(options.Context, "compensating step: %s", stepID)
			
			// Читаем запрос из store
			reqFrame := &codec.Frame{}
			if rerr := stepStore.Read(options.Context, filepath.Join(stepID, "req"), reqFrame); rerr != nil {
				w.opts.Logger.Error(options.Context, "failed to read request for compensation: %v", rerr)
				continue
			}
			
			req := &Message{Body: reqFrame.Data}
			
			// Выполняем компенсацию
			if cerr := step.Compensate(options.Context, req, opts...); cerr != nil {
				w.opts.Logger.Error(options.Context, "compensation failed for step %s: %v", stepID, cerr)
				// Продолжаем компенсацию остальных шагов даже если один не удался
			} else {
				// Обновляем статус шага после компенсации
				if werr := stepStore.Write(options.Context, filepath.Join(stepID, "status"), &codec.Frame{Data: []byte(StatusPending.String())}); werr != nil {
					w.opts.Logger.Error(options.Context, "store error: %v", werr)
				}
			}
		}
		
		return err
	}

	// Хранилище результатов шагов
	stepResults := make(map[string]*Message)
	var resultsMu sync.RWMutex

	// Получаем все вершины графа
	vertices := w.g.GetVertices()
	if len(vertices) == 0 {
		return failWorkflow(fmt.Errorf("no steps to execute"))
	}

	// Трек выполненных шагов для определения готовности
	completedSteps := make(map[string]bool)
	var completedMu sync.RWMutex

	// Канал для сбора ошибок
	errChan := make(chan error, len(vertices))
	
	// WaitGroup для ожидания завершения всех шагов
	var wg sync.WaitGroup

	// Функция проверки готовности шага (все зависимости выполнены)
	isReady := func(stepID string) bool {
		step, ok := w.steps[stepID]
		if !ok {
			return false
		}
		requires := step.Requires()
		if len(requires) == 0 {
			return true
		}
		completedMu.RLock()
		defer completedMu.RUnlock()
		for _, reqID := range requires {
			if !completedSteps[reqID] {
				return false
			}
		}
		return true
	}

	// Запускаем шаги в параллель с учетом зависимостей
	for stepID := range vertices {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			
			// Ждем пока шаг станет готовым (все зависимости выполнены)
			for {
				// Проверяем статус workflow
				statusFrame := &codec.Frame{}
				if rerr := workflowStore.Read(options.Context, "status", statusFrame); rerr != nil {
					errChan <- rerr
					return
				}
				
				currentStatus := StringStatus[string(statusFrame.Data)]
				if currentStatus != StatusRunning {
					errChan <- fmt.Errorf("workflow %s", currentStatus)
					return
				}
				
				if isReady(id) {
					break
				}
				
				// Небольшая пауза перед следующей проверкой
				select {
				case <-options.Context.Done():
					errChan <- options.Context.Err()
					return
				default:
					// Продолжаем ожидание
				}
			}
			
			step, ok := w.steps[id]
			if !ok {
				errChan <- ErrStepNotExists
				return
			}

			// Собираем результаты зависимых шагов
			requires := step.Requires()
			inputMsg := &Message{Body: []byte{}, Header: metadata.Metadata{}}
			
			if len(requires) > 0 {
				resultsMu.RLock()
				// Объединяем результаты всех зависимостей
				for _, reqID := range requires {
					if res, exists := stepResults[reqID]; exists {
						// Простая эвристика: используем последнюю зависимость или объединяем
						if len(res.Body) > 0 {
							inputMsg.Body = res.Body
							inputMsg.Header = res.Header
						}
					}
				}
				resultsMu.RUnlock()
			}

			// Сохраняем запрос в store
			if werr := stepStore.Write(options.Context, filepath.Join(id, "req"), &codec.Frame{Data: inputMsg.Body}); werr != nil {
				errChan <- werr
				return
			}

			// Устанавливаем статус Running
			step.SetStatus(StatusRunning)
			if werr := stepStore.Write(options.Context, filepath.Join(id, "status"), &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
				errChan <- werr
				return
			}

			w.opts.Logger.Info(options.Context, "executing step: %s", id)

			// Выполняем шаг
			rsp, execErr := step.Execute(options.Context, inputMsg, opts...)
			
			if execErr != nil {
				step.SetStatus(StatusFailure)
				// Сохраняем ошибку в store
				if werr := stepStore.Write(options.Context, filepath.Join(id, "rsp"), &codec.Frame{Data: []byte(execErr.Error())}); werr != nil {
					w.opts.Logger.Error(options.Context, "store error: %v", werr)
				}
				if werr := stepStore.Write(options.Context, filepath.Join(id, "status"), &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil {
					w.opts.Logger.Error(options.Context, "store error: %v", werr)
				}
				
				errChan <- execErr
				return
			}

			// Успешное выполнение
			step.SetStatus(StatusSuccess)
			
			// Сохраняем результат в store
			if rsp != nil {
				if werr := stepStore.Write(options.Context, filepath.Join(id, "rsp"), &codec.Frame{Data: rsp.Body}); werr != nil {
					errChan <- werr
					return
				}
				// Сохраняем результат для последующих шагов
				resultsMu.Lock()
				stepResults[id] = rsp
				resultsMu.Unlock()
			}
			
			if werr := stepStore.Write(options.Context, filepath.Join(id, "status"), &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
				errChan <- werr
				return
			}

			// Добавляем в список выполненных шагов для возможной компенсации
			execMu.Lock()
			executedSteps = append(executedSteps, id)
			execMu.Unlock()
			
			// Помечаем шаг как завершенный
			completedMu.Lock()
			completedSteps[id] = true
			completedMu.Unlock()

			w.opts.Logger.Info(options.Context, "step completed: %s", id)
		}(stepID)
	}

	// Ждем завершения всех горутин
	wg.Wait()
	close(errChan)

	// Проверяем наличие ошибок
	for err := range errChan {
		if err != nil {
			return failWorkflow(err)
		}
	}

	// Все шаги выполнены успешно
	if werr := workflowStore.Write(options.Context, "status", &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
		w.opts.Logger.Error(options.Context, "store error: %v", werr)
	}

	w.opts.Logger.Info(options.Context, "workflow completed successfully")
	return nil
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

func (f *microFlow) WorkflowRemove(ctx context.Context, id string) error {
	return nil
}

func (f *microFlow) WorkflowCreate(ctx context.Context, id string, steps ...Step) (Workflow, error) {
	w := &microWorkflow{opts: f.opts, id: id, g: &dag.DAG{}, steps: make(map[string]Step, len(steps))}

	for _, s := range steps {
		w.steps[s.String()] = s
		if _, err := w.g.AddVertex(s); err != nil {
			return nil, fmt.Errorf("failed to add vertex %s: %w", s.String(), err)
		}
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return nil, ErrStepNotExists
			}
			if err := w.g.AddEdge(src.String(), dst.String()); err != nil {
				return nil, fmt.Errorf("failed to add edge %s -> %s: %w", src.String(), dst.String(), err)
			}
		}
	}

	w.g.ReduceTransitively()

	w.init = true

	return w, nil
}

func (f *microFlow) WorkflowSave(ctx context.Context, w Workflow) error {
	mw, ok := w.(*microWorkflow)
	if !ok {
		return fmt.Errorf("invalid workflow type")
	}
	
	workflowStore := store.NewNamespaceStore(f.opts.Store, filepath.Join("workflows", w.ID()))
	
	// Сохраняем статус workflow
	statusFrame := &codec.Frame{Data: []byte(w.Status().String())}
	if err := workflowStore.Write(ctx, "status", statusFrame); err != nil {
		return err
	}
	
	// Сохраняем информацию о шагах и их зависимостях
	stepsData := make(map[string][]string)
	for stepID, step := range mw.steps {
		stepsData[stepID] = step.Requires()
	}
	
	stepData, err := json.Marshal(stepsData)
	if err != nil {
		return err
	}
	
	stepFrame := &codec.Frame{Data: stepData}
	if err := workflowStore.Write(ctx, "steps", stepFrame); err != nil {
		return err
	}
	
	f.opts.Logger.Info(ctx, "workflow %s saved", w.ID())
	return nil
}

func (f *microFlow) WorkflowLoad(ctx context.Context, id string) (Workflow, error) {
	workflowStore := store.NewNamespaceStore(f.opts.Store, filepath.Join("workflows", id))
	
	// Читаем статус
	statusFrame := &codec.Frame{}
	if err := workflowStore.Read(ctx, "status", statusFrame); err != nil {
		return nil, err
	}
	
	status := StringStatus[string(statusFrame.Data)]
	
	// Создаем новый workflow
	w := &microWorkflow{
		opts:   f.opts,
		id:     id,
		g:      &dag.DAG{},
		steps:  make(map[string]Step),
		status: status,
	}
	
	// Читаем информацию о шагах
	stepFrame := &codec.Frame{}
	if err := workflowStore.Read(ctx, "steps", stepFrame); err != nil {
		f.opts.Logger.Warn(ctx, "failed to read steps for workflow %s: %v", id, err)
		return nil, err
	}
	
	// Десериализуем шаги
	var stepsData map[string][]string
	if err := json.Unmarshal(stepFrame.Data, &stepsData); err != nil {
		f.opts.Logger.Warn(ctx, "failed to unmarshal steps for workflow %s: %v", id, err)
		return nil, err
	}
	
	// Восстанавливаем граф из сохраненных данных
	// Примечание: для полноценного восстановления нужны зарегистрированные шаги
	// В текущей реализации шаги должны быть созданы отдельно и добавлены через AppendSteps
	// Здесь мы только восстанавливаем структуру графа
	
	w.init = true
	
	f.opts.Logger.Info(ctx, "workflow %s loaded with status %s", id, status)
	return w, nil
}

func (f *microFlow) WorkflowList(ctx context.Context) ([]Workflow, error) {
	workflowStore := store.NewNamespaceStore(f.opts.Store, "workflows")
	
	// Получаем список всех workflow по ключам
	keys, err := workflowStore.List(ctx)
	if err != nil {
		return nil, err
	}
	
	workflows := make([]Workflow, 0, len(keys))
	for _, key := range keys {
		// Извлекаем ID workflow из ключа (последняя часть пути)
		id := filepath.Base(key)
		w, err := f.WorkflowLoad(ctx, id)
		if err != nil {
			f.opts.Logger.Error(ctx, "failed to load workflow %s: %v", key, err)
			continue
		}
		workflows = append(workflows, w)
	}
	
	return workflows, nil
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

// Compensate performs rollback for this step (default implementation returns nil)
func (s *microCallStep) Compensate(ctx context.Context, req *Message, opts ...ExecuteOption) error {
	// Default implementation does nothing - override in custom steps if compensation is needed
	return nil
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

// Compensate performs rollback for this step (default implementation returns nil)
func (s *microPublishStep) Compensate(ctx context.Context, req *Message, opts ...ExecuteOption) error {
	// Default implementation does nothing - override in custom steps if compensation is needed
	return nil
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
