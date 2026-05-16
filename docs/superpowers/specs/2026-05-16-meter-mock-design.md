# Design: meter/mock — строгий мок для meter.Meter

**Дата:** 2026-05-16  
**Модуль:** `go.unistack.org/micro/v5`  
**Путь:** `meter/mock/mock.go`, `meter/mock/mock_test.go`

---

## Цель

Реализовать мок-имплементацию интерфейса `meter.Meter` в пакете `go.unistack.org/micro/v5/meter/mock` по аналогии с `broker/mock`. Мок должен позволять:

1. Регистрировать ожидания на методы Meter (`ExpectCounter(...)`, `ExpectInit()` и т.д.)
2. Проверять, что все ожидания были выполнены через `ExpectationsWereMet()`
3. Инспектировать фактически накопленные значения метрик (Counter, Gauge и т.д.)
4. Возвращать ошибку на неожиданные вызовы (для методов, возвращающих `error`)

---

## Интерфейс meter.Meter (справка)

```go
type Meter interface {
    Name() string
    Init(opts ...Option) error
    Clone(opts ...Option) Meter
    Counter(name string, labels ...string) Counter
    FloatCounter(name string, labels ...string) FloatCounter
    Gauge(name string, fn func() float64, labels ...string) Gauge
    Set(opts ...Option) Meter
    Histogram(name string, labels ...string) Histogram
    HistogramExt(name string, quantiles []float64, labels ...string) Histogram
    Summary(name string, labels ...string) Summary
    SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) Summary
    Write(w io.Writer, opts ...Option) error
    Options() Options
    String() string
    Unregister(name string, labels ...string) bool
}
```

---

## Архитектура

### Базовая инфраструктура ожиданий (как в broker/mock)

```go
// expectation — общий интерфейс для всех Expected*
type expectation interface {
    fulfilled() bool
    Lock()
    Unlock()
    String() string
}

// commonExpectation — встраивается во все Expected-типы
type commonExpectation struct {
    sync.Mutex
    triggered bool
    err       error
}

func (e *commonExpectation) fulfilled() bool { return e.triggered }
```

### Expected-типы

| Тип | Метод Meter | Fluent-методы |
|-----|-------------|---------------|
| `ExpectedInit` | `Init()` | `WillReturnError(err)` |
| `ExpectedWrite` | `Write()` | `WillReturnError(err)` |
| `ExpectedUnregister` | `Unregister(name, labels...)` | `WillReturn(bool)` |
| `ExpectedCounter` | `Counter(name, labels...)` | содержит `*mockCounter` |
| `ExpectedFloatCounter` | `FloatCounter(name, labels...)` | содержит `*mockFloatCounter` |
| `ExpectedGauge` | `Gauge(name, fn, labels...)` | содержит `*mockGauge` |
| `ExpectedHistogram` | `Histogram(name, labels...)` | содержит `*mockHistogram` |
| `ExpectedHistogramExt` | `HistogramExt(name, quantiles, labels...)` | содержит `*mockHistogram` |
| `ExpectedSummary` | `Summary(name, labels...)` | содержит `*mockSummary` |
| `ExpectedSummaryExt` | `SummaryExt(name, window, quantiles, labels...)` | содержит `*mockSummary` |

### Мок-объекты метрик (state-tracking)

Каждый мок-объект субинтерфейса реализует соответствующий интерфейс `meter.*` и накапливает реальные значения:

- `mockCounter` — `value uint64`, защищён `sync.Mutex`
- `mockFloatCounter` — `value float64`, защищён `sync.Mutex`
- `mockGauge` — `value float64`, защищён `sync.Mutex`; `fn func() float64` для `Get()` при наличии
- `mockHistogram` — срез `updates []float64`, защищён `sync.Mutex`
- `mockSummary` — срез `updates []float64`, защищён `sync.Mutex`

Тест получает доступ через поле ожидания:

```go
exp := m.ExpectCounter("requests", "method", "GET")
// ... тестируемый код ...
if exp.Counter().Get() != 5 {
    t.Errorf("expected 5, got %d", exp.Counter().Get())
}
```

### MockMeter

```go
type MockMeter struct {
    opts       meter.Options
    mu         sync.Mutex
    expected   []expectation
    unexpected []string  // список неожиданных вызовов
}

func NewMockMeter(opts ...meter.Option) *MockMeter
```

**Поведение при вызовах:**

- Первое подходящее незатронутое ожидание срабатывает и помечается `triggered = true`
- Если подходящего ожидания нет:
  - `Init()`, `Write()` → `fmt.Errorf("unexpected call to ...")`
  - `Counter()`, `FloatCounter()`, `Gauge()`, `Histogram()`, `Summary()` → noop-объект + запись в `unexpected`
  - `Unregister()` → `false` + запись в `unexpected`
- `Clone()` и `Set()` → возвращают тот же `*MockMeter` (без необходимости регистрировать ожидание)

### ExpectationsWereMet()

```go
func (m *MockMeter) ExpectationsWereMet() error
```

Возвращает ошибку если:
1. Хотя бы одно ожидание не было выполнено (`triggered == false`)
2. Были неожиданные вызовы (непустой `unexpected`)

---

## Файлы

```
meter/
  mock/
    mock.go       — MockMeter, Expected*, mock*Counter/Gauge/Histogram/Summary
    mock_test.go  — тесты: Init, Counter.Inc, Histogram.Update, ExpectationsWereMet
```

---

## Примеры использования

```go
// Базовый сценарий
m := mock.NewMockMeter()
m.ExpectInit()
m.ExpectCounter("requests", "method", "GET")

_ = m.Init()
c := m.Counter("requests", "method", "GET")
c.Inc()
c.Inc()

exp := m.expected[1].(*mock.ExpectedCounter)
if exp.Counter().Get() != 2 {
    t.Error("expected 2 increments")
}
if err := m.ExpectationsWereMet(); err != nil {
    t.Error(err)
}

// Ожидание ошибки
m2 := mock.NewMockMeter()
m2.ExpectInit().WillReturnError(errors.New("init failed"))
if err := m2.Init(); err == nil {
    t.Error("expected error")
}

// Unregister
m3 := mock.NewMockMeter()
m3.ExpectUnregister("old_metric").WillReturn(true)
if !m3.Unregister("old_metric") {
    t.Error("expected true")
}
```

---

## Тесты

`mock_test.go` покрывает:

1. `TestMockMeter_Init` — ожидание, ошибка, неожиданный вызов
2. `TestMockMeter_Counter` — `Inc`, `Dec`, `Add`, `Set`, `Get` через `ExpectedCounter.Counter()`
3. `TestMockMeter_FloatCounter` — `Add`, `Sub`, `Set`, `Get`
4. `TestMockMeter_Gauge` — `Inc`, `Dec`, `Set`, `Get`
5. `TestMockMeter_Histogram` — `Update`, `UpdateDuration`, `Reset`
6. `TestMockMeter_HistogramExt` — то же что Histogram
7. `TestMockMeter_Summary` — `Update`, `UpdateDuration`
8. `TestMockMeter_SummaryExt` — то же что Summary
9. `TestMockMeter_Write` — ожидание, ошибка
10. `TestMockMeter_Unregister` — `WillReturn(true/false)`
11. `TestMockMeter_ExpectationsWereMet` — незакрытые ожидания, неожиданные вызовы
12. `TestMockMeter_CloneSet` — Clone и Set без ожиданий

---

## Соответствие стандартам проекта

- Compile-time check: `var _ meter.Meter = (*MockMeter)(nil)`
- Документация всех публичных типов и методов (обязательна по CLAUDE.md)
- Пакет `mock`, `package mock`
- Импорт: `go.unistack.org/micro/v5/meter`
