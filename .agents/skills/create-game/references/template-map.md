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
  api/<name>/<name>.proto
  api/gen/...                 # from buf generate
  cmd/<name>/service/main.go
  cmd/<name>/client/main.go
  internal/services/<name>/...
  internal/clients/<name>/...
  pkg/dfx/...
  pkg/modules/<name>_module.go
  deployment/docker-compose/...
  build/package/docker/Dockerfile
  tests/<name>/<name>.js
```

## Preferred creation command

```bash
.agents/skills/create-game/scripts/scaffold.sh \
  --module github.com/acme/arena \
  --name arena \
  --out ./arena
```

Do not use `gonew` / `moke-layout`.
