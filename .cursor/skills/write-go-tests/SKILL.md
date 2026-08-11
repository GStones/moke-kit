---
name: write-go-tests
description: Write or extend Go tests for moke-kit using testify and test/utils.TestHelper. Use when adding unit tests, fixing flaky tests, increasing coverage, or the user asks to test a package.
paths:
  - "**/*_test.go"
  - "**/test/**/*.go"
  - "**/*.go"
---

# Write Go tests

## Where tests live

- Prefer `_test.go` next to the code under test for unit tests.
- Broader/integration-style suites live under `test/<area>/` (e.g. `test/mq`, `test/orm`, `test/server`).
- Shared helpers: `test/utils` (`TestHelper`).

## Preferred patterns

```go
helper := utils.NewTestHelper(t)
defer helper.Cleanup()

t.Run("CaseName", func(t *testing.T) {
    // arrange / act
    helper.AssertNoError(err)
    helper.AssertEqual(want, got)
})
```

- Use `github.com/stretchr/testify/assert` and `require` (directly or via `TestHelper`).
- Use `zaptest` / `helper.Logger()` for logs; do not print with `fmt` in tests.
- Table-driven tests for pure logic and multiple input variants.
- Prefer `require` for preconditions; `assert` for multiple checks in one case.

## Quality bar

- Cover success path + at least one failure/error path for new logic.
- Exercise context cancel/timeout when the code accepts `context.Context`.
- Avoid real external network/disk when a fake, stub, or existing mock (`orm/nosql/mock`, interfaces in `miface` / `diface`) will do.
- No time-based flakes: do not `time.Sleep` to "wait for" async work; use channels, `testify` eventually helpers only if already used nearby, or deterministic fakes.
- Name cases after behavior: `ClosesSubscriptionOnUnsubscribe`, not `Test1`.

## Commands

```bash
# package under test
go test -v -race ./path/to/package

# area suite
go test -v -race ./test/mq/...

# full module (pre-PR)
go test -v -race ./...
```

## Done when

- New behavior has tests that fail if the implementation is reverted.
- `-race` clean for the touched packages.
- Helper `Cleanup` (or `t.Cleanup`) releases resources.
