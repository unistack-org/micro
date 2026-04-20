package flow

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/heimdalr/dag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStep реализует Step для тестирования
type mockStep struct {
	name        string
	delay       time.Duration
	shouldFail  bool
	compensate  bool
	executed    *sync.Map
	order       *[]string
	orderMu     *sync.Mutex
	compensated *sync.Map
}

func newMockStep(name string, delay time.Duration, shouldFail bool, executed *sync.Map, order *[]string, orderMu *sync.Mutex, compensated *sync.Map) *mockStep {
	return &mockStep{
		name:        name,
		delay:       delay,
		shouldFail:  shouldFail,
		executed:    executed,
		order:       order,
		orderMu:     orderMu,
		compensated: compensated,
	}
}

func (m *mockStep) String() string {
	return m.name
}

func (m *mockStep) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error) {
	if m.executed != nil {
		m.executed.Store(m.name, true)
	}
	if m.order != nil && m.orderMu != nil {
		m.orderMu.Lock()
		*m.order = append(*m.order, m.name)
		m.orderMu.Unlock()
	}

	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if m.shouldFail {
		return nil, fmt.Errorf("step %s failed", m.name)
	}

	return req, nil
}

func (m *mockStep) Compensate(ctx context.Context, req *Message, opts ...ExecuteOption) error {
	if m.compensated != nil {
		m.compensated.Store(m.name, true)
	}
	if m.compensate {
		return fmt.Errorf("compensation for %s failed", m.name)
	}
	return nil
}

// vertexWrapper обертка для использования строки как вершины с dag.Vertexer
type vertexWrapper struct {
	id    string
	value interface{}
}

func (v *vertexWrapper) Vertex() (string, interface{}) {
	return v.id, v.value
}

// visitorFunc адаптер для использования функции как Visitor
type visitorFunc func(dag.Vertexer)

func (vf visitorFunc) Visit(v dag.Vertexer) {
	vf(v)
}

// TestDagLinearPipeline тестирует линейный пайплайн: A -> B -> C
func TestDagLinearPipeline(t *testing.T) {
	d := dag.NewDAG()
	
	// Добавляем вершины с явным указанием ID
	require.NoError(t, d.AddVertexByID("A", "A"))
	require.NoError(t, d.AddVertexByID("B", "B"))
	require.NoError(t, d.AddVertexByID("C", "C"))

	// Добавляем ребра
	require.NoError(t, d.AddEdge("A", "B"))
	require.NoError(t, d.AddEdge("B", "C"))

	// Получаем топологический порядок через OrderedWalk
	var order []string
	d.OrderedWalk(visitorFunc(func(v dag.Vertexer) {
		id, _ := v.Vertex()
		order = append(order, id)
	}))

	// Проверяем порядок
	assert.Equal(t, []string{"A", "B", "C"}, order)
}

// TestDagParallelBranches тестирует параллельное выполнение: A -> (B, C) -> D
func TestDagParallelBranches(t *testing.T) {
	d := dag.NewDAG()

	vertices := []string{"A", "B", "C", "D"}
	for _, v := range vertices {
		require.NoError(t, d.AddVertexByID(v, v))
	}

	// A -> B, A -> C
	require.NoError(t, d.AddEdge("A", "B"))
	require.NoError(t, d.AddEdge("A", "C"))
	// B -> D, C -> D
	require.NoError(t, d.AddEdge("B", "D"))
	require.NoError(t, d.AddEdge("C", "D"))

	var order []string
	d.OrderedWalk(visitorFunc(func(v dag.Vertexer) {
		id, _ := v.Vertex()
		order = append(order, id)
	}))

	// A должен быть первым, D последним
	assert.Equal(t, "A", order[0])
	assert.Equal(t, "D", order[len(order)-1])

	// B и C должны быть между A и D
	foundB, foundC := false, false
	for _, v := range order[1 : len(order)-1] {
		if v == "B" {
			foundB = true
		}
		if v == "C" {
			foundC = true
		}
	}
	assert.True(t, foundB)
	assert.True(t, foundC)
}

// TestDagDiamondPattern тестирует паттерн "алмаз": A -> (B, C) -> D
func TestDagDiamondPattern(t *testing.T) {
	d := dag.NewDAG()

	vertices := []string{"Start", "Left", "Right", "End"}
	for _, v := range vertices {
		require.NoError(t, d.AddVertexByID(v, v))
	}

	require.NoError(t, d.AddEdge("Start", "Left"))
	require.NoError(t, d.AddEdge("Start", "Right"))
	require.NoError(t, d.AddEdge("Left", "End"))
	require.NoError(t, d.AddEdge("Right", "End"))

	var order []string
	d.OrderedWalk(visitorFunc(func(v dag.Vertexer) {
		id, _ := v.Vertex()
		order = append(order, id)
	}))

	assert.Equal(t, "Start", order[0])
	assert.Equal(t, "End", order[len(order)-1])
}

// TestDagExecutionWithCompensation тестирует выполнение с компенсациями при ошибке
func TestDagExecutionWithCompensation(t *testing.T) {
	var (
		executed    sync.Map
		compensated sync.Map
		order       []string
		orderMu     sync.Mutex
	)

	d := dag.NewDAG()
	steps := map[string]*mockStep{
		"Step1": newMockStep("Step1", 10*time.Millisecond, false, &executed, &order, &orderMu, &compensated),
		"Step2": newMockStep("Step2", 10*time.Millisecond, false, &executed, &order, &orderMu, &compensated),
		"Step3": newMockStep("Step3", 10*time.Millisecond, true, &executed, &order, &orderMu, &compensated), // упадет
	}

	for name, step := range steps {
		require.NoError(t, d.AddVertexByID(name, step))
	}

	// Линейная зависимость: Step1 -> Step2 -> Step3
	require.NoError(t, d.AddEdge("Step1", "Step2"))
	require.NoError(t, d.AddEdge("Step2", "Step3"))

	// Симулируем выполнение и компенсацию
	ctx := context.Background()
	failed := false
	
	// Выполняем шаги по порядку через OrderedWalk
	var order_exec []string
	d.OrderedWalk(visitorFunc(func(v dag.Vertexer) {
		id, val := v.Vertex()
		order_exec = append(order_exec, id)
		step := val.(*mockStep)
		_, err := step.Execute(ctx, &Message{})
		if err != nil && !failed {
			failed = true
			// Запускаем компенсации в обратном порядке для уже выполненных
			for i := len(order_exec) - 1; i >= 0; i-- {
				compID := order_exec[i]
				if compStep, ok := steps[compID]; ok {
					_ = compStep.Compensate(ctx, &Message{})
				}
			}
		}
	}))

	assert.True(t, failed, "ошибка должна была произойти")
	
	// Проверяем что все выполненные шаги были скомпенсированы
	_, ok1 := executed.Load("Step1")
	assert.True(t, ok1)
	_, ok2 := executed.Load("Step2")
	assert.True(t, ok2)
	_, ok3 := compensated.Load("Step1")
	assert.True(t, ok3)
	_, ok4 := compensated.Load("Step2")
	assert.True(t, ok4)
}

// TestDagComplexGraph тестирует сложный граф с несколькими уровнями
func TestDagComplexGraph(t *testing.T) {
	d := dag.NewDAG()

	// Создаем сложный граф:
	// A -> B -> D
	// A -> C -> D
	// D -> E
	vertices := []string{"A", "B", "C", "D", "E"}
	for _, v := range vertices {
		require.NoError(t, d.AddVertexByID(v, v))
	}

	require.NoError(t, d.AddEdge("A", "B"))
	require.NoError(t, d.AddEdge("A", "C"))
	require.NoError(t, d.AddEdge("B", "D"))
	require.NoError(t, d.AddEdge("C", "D"))
	require.NoError(t, d.AddEdge("D", "E"))

	var order []string
	d.OrderedWalk(visitorFunc(func(v dag.Vertexer) {
		id, _ := v.Vertex()
		order = append(order, id)
	}))

	// Проверка зависимостей
	pos := make(map[string]int)
	for i, v := range order {
		pos[v] = i
	}

	assert.Less(t, pos["A"], pos["B"])
	assert.Less(t, pos["A"], pos["C"])
	assert.Less(t, pos["B"], pos["D"])
	assert.Less(t, pos["C"], pos["D"])
	assert.Less(t, pos["D"], pos["E"])
}

// TestDagCycleDetection тестирует обнаружение циклов
func TestDagCycleDetection(t *testing.T) {
	d := dag.NewDAG()

	require.NoError(t, d.AddVertexByID("A", "A"))
	require.NoError(t, d.AddVertexByID("B", "B"))

	require.NoError(t, d.AddEdge("A", "B"))
	
	// Пытаемся создать цикл
	err := d.AddEdge("B", "A")
	assert.Error(t, err, "должна быть ошибка при создании цикла")
}
