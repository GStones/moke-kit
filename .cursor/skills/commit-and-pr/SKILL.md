---
name: commit-and-pr
description: Create git commits and pull requests with moke-kit conventions. Use when the user asks to commit, push, open a PR, write a commit message, or prepare changes for review. Prefer explicit invocation via /commit-and-pr.
disable-model-invocation: true
---

# Commit and PR

## Commit messages

Prefer Conventional Commits when the change type is clear:

```text
<type>(optional-scope): <short summary>
```

Common types used in this repo: `fix`, `feat`, `chore`, `build`, `refactor`, `test`, `docs`.

Examples:

- `fix(server): fail closed when auth middleware missing in prod`
- `test(mq): cover nats subscription error paths`
- `chore(deps): bump go.opentelemetry.io/otel to v1.45.0`

Rules:

- Imperative mood, ~72 chars on the subject line.
- Explain **why** in the body when the diff is non-obvious.
- Never commit secrets, tokens, or local env files.
- Do not amend commits unless the user explicitly asks and the amend conditions are met.

## Before committing

1. `git status` / `git diff` — stage only relevant files.
2. Run the `go-test-and-lint` skill checks on the touched scope when Go code changed.
3. Confirm no generated noise or accidental local files are staged.

## Pull requests

- Title: same style as commit subjects; concise and specific.
- Body should include:
  - **Summary**: what changed and why (2–4 bullets).
  - **Test plan**: commands run and results (e.g. `go test -race ./...`, `gofmt`, `staticcheck`).
- Keep PRs focused; split unrelated work.
- Link related issues when known (`Fixes #N` / `Refs #N`).

## Safety

- Do not force-push to `main` / `master`.
- Do not push secrets or rewrite shared history.
- Ask before changing git config or performing destructive git operations.
