# increase-coverage Design

## Overview

Increase unit test coverage to 85% for 7 core packages in the micro framework: server, client, tracer, broker, store, meter, logger. Additionally, these tests will validate framework behavior by covering interfaces, error handling, concurrent safety, and options.

## Current State

| Package   | Current Coverage | Target |
|-----------|------------------|--------|
| server    | 7.2%            | 85%    |
| client    | 12.4%           | 85%    |
| tracer    | 11.9%           | 85%    |
| store     | 13.9%           | 85%    |
| broker    | 25.0%           | 85%    |
| meter     | 55.8%           | 85%    |
| logger    | 84.0%           | 85%    |

## Approach

**Coverage-guided unit testing**: Use `go test -coverprofile=coverage.out` and `go tool cover -html=coverage.out` to identify uncovered lines, then write targeted unit tests.

**Order of work**: From lowest to highest coverage:
1. server (7.2%)
2. client (12.4%)
3. tracer (11.9%)
4. store (13.9%)
5. broker (25.0%)
6. meter (55.8%)
7. logger (84.0%)

## Testing Strategy

- **Tools**: `go test -coverprofile`, table-driven tests (standard for Go), existing mocks (store/mock, flow/mock).
- **Process per package**:
  1. Generate coverage profile.
  2. Analyze uncovered functions/branches.
  3. Write unit tests (mock external dependencies if needed).
  4. Re-check coverage.
  5. Repeat until ≥85%.
- **Linting**: Run `golangci-lint run` to ensure code quality.

## Framework Validation

Tests will verify:
1. **Interface correctness** — all implementations (memory, noop, mock) satisfy contracts.
2. **Error handling** — cover error scenarios (timeout, nil pointers, invalid input).
3. **Concurrent safety** — race condition tests for server, broker, store.
4. **Options/Defaults** — ensure default values are set correctly.

## Package-specific Notes

- **logger** (84%→85%): Minimal work, likely edge cases in slog adapter.
- **meter** (55.8%→85%): Add tests for metrics, histogram, counter.
- **tracer** (11.9%→85%): Cover memory implementation and context helpers.
- **server, client, store, broker**: Major effort, follow coverage-guided approach.

## Success Criteria

- Each targeted package reaches ≥85% coverage.
- All tests pass.
- No new linting errors introduced.
- Design approved and committed.
