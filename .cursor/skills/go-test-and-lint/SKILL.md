---
name: go-test-and-lint
description: Run and fix Go tests and linters for moke-kit. Use when adding or changing Go code, fixing CI failures, or when the user asks to test, lint, vet, gofmt, staticcheck, or verify a package.
paths:
  - "**/*.go"
  - "**/go.mod"
  - "**/go.sum"
  - ".github/workflows/**"
---

# Go test and lint

Align local checks with CI (`.github/workflows/go.yml`, `.github/workflows/reviewdog.yml`).

## Workflow

1. Identify the change scope (package or `./...`).
2. Format: `gofmt -s -w` on touched `.go` files.
3. Vet: `go vet ./...` (or the narrowed package path).
4. Test with race + coverage (same spirit as CI):

   ```bash
   go test -v -race -coverpkg=./... -coverprofile=coverage.txt ./...
   ```

   For a single package, prefer:

   ```bash
   go test -v -race ./path/to/package
   ```

5. If available, run `staticcheck ./...` (primary lint gate in reviewdog). Fix failures before finishing.
6. Re-run the failing command after each fix until clean.

## Rules

- Prefer package-scoped tests while iterating; run full `./...` before commit/PR when changes span modules.
- Do not ignore race failures; fix data races or document why a test cannot use `-race` (rare).
- Keep tests deterministic; no sleeps for synchronization when channels/`sync` primitives suffice.
- Match existing import grouping and `gofmt` style; do not fight the formatter.

## CI mapping

| Check | Command / tool |
| --- | --- |
| Unit tests | `go test -v -race ...` |
| Vet | `go vet ./...` |
| Format | `gofmt -s` (reviewdog) |
| Static analysis | `staticcheck` (reviewdog, fail-on-error) |
| Style (informational) | `golint` via reviewdog |

## Done when

- Touched packages pass `go test` (with `-race` when practical).
- `go vet` and `gofmt -s` are clean on changed files.
- No new `staticcheck` findings introduced by the change.
