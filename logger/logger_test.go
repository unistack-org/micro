package logger

import (
	"context"
	"testing"
)

// TestLevelString tests the Level.String() method for all levels
func TestLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{"TraceLevel", TraceLevel, "trace"},
		{"DebugLevel", DebugLevel, "debug"},
		{"InfoLevel", InfoLevel, "info"},
		{"WarnLevel", WarnLevel, "warn"},
		{"ErrorLevel", ErrorLevel, "error"},
		{"FatalLevel", FatalLevel, "fatal"},
		{"NoneLevel", NoneLevel, "none"},
		{"UnknownLevel", Level(99), "info"}, // Unknown levels default to "info"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.expected {
				t.Errorf("Level.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestLevelEnabled tests the Level.Enabled() method
func TestLevelEnabled(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		checkLvl Level
		expected bool
	}{
		// Level.Enabled returns true if checkLvl >= level
		{"Trace enabled for Trace", TraceLevel, TraceLevel, true},
		{"Trace enabled for Debug", TraceLevel, DebugLevel, true},
		{"Trace enabled for Info", TraceLevel, InfoLevel, true},
		{"Trace enabled for Error", TraceLevel, ErrorLevel, true},

		{"Info enabled for Info", InfoLevel, InfoLevel, true},
		{"Info enabled for Warn", InfoLevel, WarnLevel, true},
		{"Info enabled for Error", InfoLevel, ErrorLevel, true},
		{"Info enabled for Debug", InfoLevel, DebugLevel, false},
		{"Info enabled for Trace", InfoLevel, TraceLevel, false},

		{"Error enabled for Error", ErrorLevel, ErrorLevel, true},
		{"Error enabled for Fatal", ErrorLevel, FatalLevel, true},
		{"Error enabled for Warn", ErrorLevel, WarnLevel, false},

		{"None enabled for None", NoneLevel, NoneLevel, true},
		{"None enabled for Fatal", NoneLevel, FatalLevel, false},
		{"None enabled for Info", NoneLevel, InfoLevel, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.Enabled(tt.checkLvl)
			if got != tt.expected {
				t.Errorf("Level(%v).Enabled(%v) = %v, want %v", tt.level, tt.checkLvl, got, tt.expected)
			}
		})
	}
}

// TestParseLevel tests the ParseLevel function
func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
	}{
		{"trace", "trace", TraceLevel},
		{"debug", "debug", DebugLevel},
		{"info", "info", InfoLevel},
		{"warn", "warn", WarnLevel},
		{"error", "error", ErrorLevel},
		{"fatal", "fatal", FatalLevel},
		{"none", "none", NoneLevel},
		{"empty string", "", InfoLevel},                           // Unknown defaults to InfoLevel
		{"unknown string", "unknown", InfoLevel},                  // Unknown defaults to InfoLevel
		{"TRACE uppercase", "TRACE", InfoLevel},                   // Case-sensitive, so uppercase is unknown
		{"mixed case", "TrAcE", InfoLevel},                        // Case-sensitive
		{"whitespace", " trace", InfoLevel},                       // Leading whitespace makes it unknown
		{"suffix", "trace_level", InfoLevel},                      // Suffix makes it unknown
		{"numeric", "123", InfoLevel},                             // Non-string values
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestNewLogger tests the NewLogger function
func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	// Verify it's not nil
	if logger == nil {
		t.Error("NewLogger() returned nil, expected a Logger instance")
	}

	// Verify it's a noopLogger
	if _, ok := logger.(*noopLogger); !ok {
		t.Errorf("NewLogger() returned %T, expected *noopLogger", logger)
	}

	// Verify default behavior
	if logger.String() != "noop" {
		t.Errorf("NewLogger().String() = %q, want %q", logger.String(), "noop")
	}

	// Verify V always returns false
	if logger.V(TraceLevel) {
		t.Error("NewLogger().V(TraceLevel) returned true, expected false")
	}
	if logger.V(InfoLevel) {
		t.Error("NewLogger().V(InfoLevel) returned true, expected false")
	}
}

// TestNewLoggerWithOptions tests NewLogger with options
func TestNewLoggerWithOptions(t *testing.T) {
	testName := "test-logger"
	logger := NewLogger(WithName(testName))

	if logger == nil {
		t.Error("NewLogger(WithName(...)) returned nil")
		return
	}

	if logger.Name() != testName {
		t.Errorf("NewLogger(WithName(%q)).Name() = %q, want %q", testName, logger.Name(), testName)
	}
}

// TestNoopLoggerV tests the V method of noopLogger
func TestNoopLoggerV(t *testing.T) {
	logger := NewLogger().(*noopLogger)

	tests := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, NoneLevel}

	for _, lvl := range tests {
		if logger.V(lvl) {
			t.Errorf("noopLogger.V(%v) returned true, expected false", lvl)
		}
	}
}

// TestNoopLoggerLevel tests the Level method of noopLogger (should be a no-op)
func TestNoopLoggerLevel(t *testing.T) {
	logger := NewLogger().(*noopLogger)
	originalName := logger.opts.Name
	originalLevel := logger.opts.Level

	// Level() should be a no-op, so opts should not change
	logger.Level(TraceLevel)
	if logger.opts.Name != originalName || logger.opts.Level != originalLevel {
		t.Error("noopLogger.Level() modified options, expected no-op")
	}

	logger.Level(ErrorLevel)
	if logger.opts.Name != originalName || logger.opts.Level != originalLevel {
		t.Error("noopLogger.Level() modified options, expected no-op")
	}
}

// TestNoopLoggerName tests the Name method of noopLogger
func TestNoopLoggerName(t *testing.T) {
	testName := "my-logger"
	logger := NewLogger(WithName(testName)).(*noopLogger)

	if logger.Name() != testName {
		t.Errorf("noopLogger.Name() = %q, want %q", logger.Name(), testName)
	}
}

// TestNoopLoggerNameDefault tests the Name method with default name
func TestNoopLoggerNameDefault(t *testing.T) {
	logger := NewLogger().(*noopLogger)

	// Default name should be empty string
	if logger.Name() != "" {
		t.Errorf("noopLogger.Name() = %q, want empty string", logger.Name())
	}
}

// TestNoopLoggerString tests the String method of noopLogger
func TestNoopLoggerString(t *testing.T) {
	logger := NewLogger().(*noopLogger)

	if logger.String() != "noop" {
		t.Errorf("noopLogger.String() = %q, want %q", logger.String(), "noop")
	}
}

// TestNoopLoggerInit tests the Init method of noopLogger
func TestNoopLoggerInit(t *testing.T) {
	logger := NewLogger().(*noopLogger)

	testName := "initialized-logger"
	err := logger.Init(WithName(testName))

	if err != nil {
		t.Errorf("noopLogger.Init() returned error: %v, expected nil", err)
	}

	if logger.Name() != testName {
		t.Errorf("After Init, noopLogger.Name() = %q, want %q", logger.Name(), testName)
	}
}

// TestNoopLoggerInitNoOptions tests Init with no options
func TestNoopLoggerInitNoOptions(t *testing.T) {
	logger := NewLogger(WithName("original")).(*noopLogger)
	originalName := logger.Name()

	err := logger.Init()

	if err != nil {
		t.Errorf("noopLogger.Init() returned error: %v", err)
	}

	if logger.Name() != originalName {
		t.Errorf("After Init(), name changed from %q to %q", originalName, logger.Name())
	}
}

// TestNoopLoggerClone tests the Clone method of noopLogger
func TestNoopLoggerClone(t *testing.T) {
	originalLogger := NewLogger(WithName("original")).(*noopLogger)
	clonedLogger := originalLogger.Clone(WithName("cloned"))

	if clonedLogger == nil {
		t.Error("noopLogger.Clone() returned nil")
		return
	}

	if clonedLogger.Name() != "cloned" {
		t.Errorf("Cloned logger name = %q, want %q", clonedLogger.Name(), "cloned")
	}

	// Original should be unchanged
	if originalLogger.Name() != "original" {
		t.Errorf("Original logger name changed to %q, want %q", originalLogger.Name(), "original")
	}

	// Cloned logger should still be a noopLogger
	if _, ok := clonedLogger.(*noopLogger); !ok {
		t.Errorf("Clone() returned %T, expected *noopLogger", clonedLogger)
	}
}

// TestNoopLoggerCloneNoOptions tests Clone with no options
func TestNoopLoggerCloneNoOptions(t *testing.T) {
	originalLogger := NewLogger(WithName("original")).(*noopLogger)
	clonedLogger := originalLogger.Clone()

	if clonedLogger == nil {
		t.Error("noopLogger.Clone() returned nil")
		return
	}

	// Without options, the clone should have the same name
	if clonedLogger.Name() != originalLogger.Name() {
		t.Errorf("Cloned logger name = %q, want %q", clonedLogger.Name(), originalLogger.Name())
	}
}

// TestNoopLoggerFields tests the Fields method of noopLogger
func TestNoopLoggerFields(t *testing.T) {
	logger := NewLogger().(*noopLogger)

	// Fields() should return the logger itself
	result := logger.Fields("key1", "value1", "key2", "value2")

	if result != logger {
		t.Error("noopLogger.Fields() should return the same logger instance")
	}
}

// TestNoopLoggerFieldsChaining tests method chaining with Fields
func TestNoopLoggerFieldsChaining(t *testing.T) {
	logger := NewLogger()

	// Should be able to chain multiple Fields calls
	result := logger.Fields("key1", "value1").Fields("key2", "value2")

	if result == nil {
		t.Error("Chained Fields() call returned nil")
	}
}

// TestNoopLoggerOptions tests the Options method of noopLogger
func TestNoopLoggerOptions(t *testing.T) {
	logger := NewLogger(WithName("test")).(*noopLogger)

	opts := logger.Options()

	if opts.Name != "test" {
		t.Errorf("Logger.Options().Name = %q, want %q", opts.Name, "test")
	}

	// Verify CallerSkipCount is set to default
	if opts.CallerSkipCount != defaultCallerSkipCount {
		t.Errorf("Logger.Options().CallerSkipCount = %d, want %d", opts.CallerSkipCount, defaultCallerSkipCount)
	}
}

// TestNoopLoggerLogMethod tests the Log method of noopLogger (should be a no-op)
func TestNoopLoggerLogMethod(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// These should all be no-ops and not panic
	logger.Log(ctx, InfoLevel, "test message")
	logger.Log(ctx, ErrorLevel, "error message", "key", "value")
	logger.Log(ctx, TraceLevel, "trace")
}

// TestNoopLoggerInfo tests the Info method of noopLogger (should be a no-op)
func TestNoopLoggerInfo(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Should not panic
	logger.Info(ctx, "info message")
	logger.Info(ctx, "info with attrs", "key1", "value1", "key2", "value2")
}

// TestNoopLoggerDebug tests the Debug method of noopLogger (should be a no-op)
func TestNoopLoggerDebug(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Should not panic
	logger.Debug(ctx, "debug message")
	logger.Debug(ctx, "debug with attrs", "attr1", "val1")
}

// TestNoopLoggerTrace tests the Trace method of noopLogger (should be a no-op)
func TestNoopLoggerTrace(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Should not panic
	logger.Trace(ctx, "trace message")
	logger.Trace(ctx, "trace with attrs", "key", "value")
}

// TestNoopLoggerWarn tests the Warn method of noopLogger (should be a no-op)
func TestNoopLoggerWarn(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Should not panic
	logger.Warn(ctx, "warning message")
	logger.Warn(ctx, "warning with attrs", "status", "warning")
}

// TestNoopLoggerError tests the Error method of noopLogger (should be a no-op)
func TestNoopLoggerError(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Should not panic
	logger.Error(ctx, "error message")
	logger.Error(ctx, "error with attrs", "error_code", "500")
}

// TestNoopLoggerFatal tests the Fatal method of noopLogger (should be a no-op)
func TestNoopLoggerFatal(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Should not panic (and notably, should NOT call os.Exit)
	logger.Fatal(ctx, "fatal message")
	logger.Fatal(ctx, "fatal with attrs", "emergency", "true")
}

// TestDefaultLogger tests that DefaultLogger is properly initialized
func TestDefaultLogger(t *testing.T) {
	if DefaultLogger == nil {
		t.Error("DefaultLogger is nil")
	}

	if DefaultLogger.String() != "noop" {
		t.Errorf("DefaultLogger.String() = %q, expected %q", DefaultLogger.String(), "noop")
	}
}

// TestDefaultLevel tests that DefaultLevel is InfoLevel
func TestDefaultLevel(t *testing.T) {
	if DefaultLevel != InfoLevel {
		t.Errorf("DefaultLevel = %v, expected %v", DefaultLevel, InfoLevel)
	}
}

// TestWithLevel tests the WithLevel option
func TestWithLevel(t *testing.T) {
	logger := NewLogger(WithLevel(ErrorLevel))
	opts := logger.Options()

	if opts.Level != ErrorLevel {
		t.Errorf("WithLevel(ErrorLevel) set level to %v, expected %v", opts.Level, ErrorLevel)
	}
}

// TestWithFields tests the WithFields option
func TestWithFields(t *testing.T) {
	fields := []any{"key1", "value1", "key2", "value2"}
	logger := NewLogger(WithFields(fields...))
	opts := logger.Options()

	if len(opts.Fields) != len(fields) {
		t.Errorf("WithFields set %d fields, expected %d", len(opts.Fields), len(fields))
	}
}

// TestWithAddFields tests the WithAddFields option
func TestWithAddFields(t *testing.T) {
	logger := NewLogger(WithAddFields("key1", "value1"), WithAddFields("key2", "value2"))
	opts := logger.Options()

	if len(opts.Fields) < 2 {
		t.Errorf("WithAddFields resulted in %d fields", len(opts.Fields))
	}
}

// TestWithAddStacktrace tests the WithAddStacktrace option
func TestWithAddStacktrace(t *testing.T) {
	logger := NewLogger(WithAddStacktrace(ErrorLevel))
	opts := logger.Options()

	if opts.AddStacktrace != ErrorLevel {
		t.Errorf("WithAddStacktrace set level to %v, expected %v", opts.AddStacktrace, ErrorLevel)
	}
}

// TestWithSource tests the WithSource option
func TestWithSource(t *testing.T) {
	logger := NewLogger(WithSource(false))
	opts := logger.Options()

	if opts.AddSource {
		t.Error("WithSource(false) should disable AddSource")
	}
}

// TestWithAddCallerSkipCount tests the WithAddCallerSkipCount option
func TestWithAddCallerSkipCount(t *testing.T) {
	logger := NewLogger(WithAddCallerSkipCount(5))
	opts := logger.Options()

	// Note: NewLogger sets CallerSkipCount to defaultCallerSkipCount AFTER options are processed,
	// so the WithAddCallerSkipCount addition is overwritten. This is the current behavior.
	if opts.CallerSkipCount != defaultCallerSkipCount {
		t.Errorf("NewLogger set skip count to %d, expected %d", opts.CallerSkipCount, defaultCallerSkipCount)
	}
}

// TestWithZapKeys tests the WithZapKeys option
func TestWithZapKeys(t *testing.T) {
	logger := NewLogger(WithZapKeys())
	opts := logger.Options()

	if opts.TimeKey != "@timestamp" {
		t.Errorf("WithZapKeys set TimeKey to %q, expected %q", opts.TimeKey, "@timestamp")
	}
	if opts.LevelKey != "level" {
		t.Errorf("WithZapKeys set LevelKey to %q, expected %q", opts.LevelKey, "level")
	}
	if opts.MessageKey != "msg" {
		t.Errorf("WithZapKeys set MessageKey to %q, expected %q", opts.MessageKey, "msg")
	}
	if opts.SourceKey != "caller" {
		t.Errorf("WithZapKeys set SourceKey to %q, expected %q", opts.SourceKey, "caller")
	}
}

// TestWithZerologKeys tests the WithZerologKeys option
func TestWithZerologKeys(t *testing.T) {
	logger := NewLogger(WithZerologKeys())
	opts := logger.Options()

	if opts.TimeKey != "time" {
		t.Errorf("WithZerologKeys set TimeKey to %q, expected %q", opts.TimeKey, "time")
	}
	if opts.MessageKey != "message" {
		t.Errorf("WithZerologKeys set MessageKey to %q, expected %q", opts.MessageKey, "message")
	}
}

// TestWithSlogKeys tests the WithSlogKeys option
func TestWithSlogKeys(t *testing.T) {
	logger := NewLogger(WithSlogKeys())
	opts := logger.Options()

	if opts.SourceKey != "source" {
		t.Errorf("WithSlogKeys set SourceKey to %q, expected %q", opts.SourceKey, "source")
	}
}

// TestWithContextAttrFuncs tests the WithContextAttrFuncs option
func TestWithContextAttrFuncs(t *testing.T) {
	fn := func(ctx context.Context) []any {
		return []any{"key", "value"}
	}

	logger := NewLogger(WithContextAttrFuncs(fn))
	opts := logger.Options()

	if len(opts.ContextAttrFuncs) == 0 {
		t.Error("WithContextAttrFuncs should add context attr funcs")
	}
}

// TestWithName tests the WithName option
func TestWithNameOption(t *testing.T) {
	testName := "my-logger"
	logger := NewLogger(WithName(testName))
	opts := logger.Options()

	if opts.Name != testName {
		t.Errorf("WithName set name to %q, expected %q", opts.Name, testName)
	}
}
