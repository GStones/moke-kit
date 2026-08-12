# create-game template map

Templates live in `assets/template/`. Placeholders:

| Placeholder | Example |
| --- | --- |
| `__MODULE__` | `github.com/acme/arena` |
| `__NAME__` | `arena` |
| `__NAME_TITLE__` | `Arena` (first letter uppercased) |
| `__BUF_MODULE__` | `buf.build/acme/arena` |

**Replace `__NAME_TITLE__` and `__BUF_MODULE__` before `__NAME__` / `__MODULE__`.**

## Output tree

```text
<out>/
  go.mod
  buf.yaml
  buf.gen.yaml
  README.md
  .env.example                # AUTH/TLS/CORS/NATS/Mongo/Redis
  api/<name>/<name>.proto
  api/gen/...                 # from buf generate
  cmd/<name>/service/main.go       # AllModule + dfx.AuthModule stub (TCP opt-in)
  cmd/<name>/service-thin/main.go  # game-only thin topology (swap stub for platform middleware)
  cmd/<name>/client/main.go
  internal/services/<name>/...  # no WithoutAuth; UID from auth context only
  internal/clients/<name>/...
  pkg/dfx/...                 # AuthModule stub (swap for platform in main; not both)
  pkg/modules/<name>_module.go  # Grpc/Http/All modules have no auth provider
  deployment/docker-compose/...
  build/package/docker/Dockerfile
  .github/workflows/go.yml      # lint/vet/govulncheck/buf breaking
  .golangci.yml
  tests/common/common.js
  tests/<name>/<name>.js        # stub bearer or USE_PLATFORM_AUTH=1 Authenticate
  tests/<name>/*-k6.proto
  internal/services/<name>/domain/watch_test.go
```

## Preferred creation command

```bash
.agents/skills/create-game/scripts/scaffold.sh \
  --module github.com/acme/arena \
  --name arena \
  --out ./arena
```

Do not use `gonew` / `moke-layout`.
