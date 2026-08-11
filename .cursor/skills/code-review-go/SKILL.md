---
name: code-review-go
description: Review Go changes in moke-kit for correctness, concurrency, security, fx wiring, and CI risk. Use when the user asks for a code review, PR review, or to check a diff before merge.
paths:
  - "**/*.go"
  - "**/go.mod"
  - "**/go.sum"
---

# Go code review (moke-kit)

Review the actual diff and surrounding call sites. Prefer findings over praise. Order by severity.

## Severity order

1. **Correctness** — logic bugs, nil deref, wrong error handling, broken API contracts
2. **Security** — authz bypass, TLS/mTLS misuse, path traversal, secret leakage, unsafe defaults in prod
3. **Concurrency** — races, missing locks, ctx cancellation ignored, goroutine leaks
4. **Reliability** — resource leaks, missing timeouts, silent error drops
5. **Maintainability** — fx module boundaries, naming, test gaps (only if actionable)

## moke-kit-specific checks

- **Deployment fail-closed**: prod (`prod`, `prod_*`, `prod-*`) must not weaken auth/TLS defaults. Local/dev convenience must not leak into prod paths.
- **fx wiring**: new providers use `fx.In` / `fx.Out` consistently; modules live under `pkg/module` and `pkg/*fx` (e.g. `sfx`, `mfx`, `ofx`); avoid circular module imports.
- **Servers**: gRPC / gateway / zinx changes must respect cmux, TLS/mTLS settings, and middleware order.
- **ORM / MQ**: document store, Redis, Mongo, NATS changes should preserve interface contracts in `diface` / `miface` and error types in `nerrors` / `qerrors`.
- **Config**: envconfig defaults documented; empty/zero values must not silently disable security controls.

## Review output format

For each finding:

```markdown
### [severity] Title
- **Where**: `path/file.go` (function or symbol)
- **Why**: concrete failure mode
- **Fix**: smallest actionable suggestion
```

End with:

- **Blockers**: must fix before merge
- **Nits**: optional

## Do not

- Nitpick style that `gofmt` / `staticcheck` already enforce.
- Request broad refactors unrelated to the diff.
- Claim a race or leak without pointing at the code path.
