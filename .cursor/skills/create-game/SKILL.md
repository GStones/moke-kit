---
name: create-game
description: Scaffold a new game server directly from moke-kit create-game templates (no gonew, no moke-layout). Use when the user asks to create a game, new game project, bootstrap from moke-kit, or generate a game server scaffold.
---

# Create a game based on moke-kit

Scaffold the project **from this skill’s templates**. Do **not** use `gonew`, clone `moke-layout`, or create a `demo` service then rename it.

## Mental model

```text
moke-kit create-game skill  → writes a ready game repo (templates + scaffold.sh)
moke-kit modules            → runtime infra via fxmain.Main
platform (optional)         → shared services imported later as fx modules
```

## Gather requirements

Ask only if missing:

1. **Module path** — e.g. `github.com/acme/arena`
2. **Service name** — lowercase (`arena`); becomes paths, proto package, binary name
3. **Output directory** — default `./{name}` (outside or beside this repo as the user prefers)
4. **Platform stubs?** — default no; pass `--with-platform` only if requested

## Scaffold (required path)

From the **moke-kit repository root**:

```bash
chmod +x .cursor/skills/create-game/scripts/scaffold.sh

.cursor/skills/create-game/scripts/scaffold.sh \
  --module github.com/<org>/<game> \
  --name <name> \
  --out ./<name>
```

Optional:

```bash
# keep commented platform import stubs in cmd/<name>/service/main.go
--with-platform

# override Buf Schema Registry module name
--buf-module buf.build/<org>/<name>
```

The script:

1. Copies [`assets/template/`](assets/template/) with `__NAME__` / `__MODULE__` / `__NAME_TITLE__` / `__BUF_MODULE__` substituted
2. Runs `go get github.com/gstones/moke-kit@latest` + `go mod tidy`
3. Runs `buf dep update` + `buf generate` when `buf` is installed

If the script cannot run, manually copy `assets/template/`, rename `__NAME__` path segments, substitute placeholders (**`__NAME_TITLE__` before `__NAME__`**), then run `go mod tidy` and `buf generate`. See [references/template-map.md](references/template-map.md).

## After scaffold

```bash
cd <out>
docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d
go run ./cmd/<name>/service/main.go
```

Smoke:

```bash
go build -o <name> ./cmd/<name>/client/main.go
./<name> grpc
# HTTP: POST localhost:8081/v1/hello/hi
```

## What gets generated

| Area | Contents |
| --- | --- |
| `cmd/<name>/service` | `fxmain.Main` + NATS/local MQ + Redis cache + `AllModule` |
| `cmd/<name>/client` | Interactive grpc/tcp CLI |
| `api/<name>` | Proto with `Hi` + streaming `Watch` |
| `internal/services/<name>` | Service registration (`sfx.*ServiceResult`), domain, nosql DAO |
| `pkg/modules` | `GrpcModule` / `HttpModule` / `TcpModule` / `AllModule` / client |
| `pkg/dfx` | Settings, auth middleware stub, grpc client provider |
| `deployment/docker-compose` | Redis, Mongo, NATS |
| `tests/<name>` | k6 smoke for `Hi` |

Default transports: gRPC + gateway + TCP. Service embeds `utility.WithoutAuth` so local smoke works; wire real auth before prod.

## Hard rules

- **Never** call `gonew` or depend on `github.com/gstones/moke-layout` for this workflow
- **Never** leave a `demo` / `game0` name unless the user chose that as `--name`
- Do not hand-edit `api/gen/**`; only via `buf generate`
- Register handlers via `sfx.GrpcServiceResult` / `GatewayServiceResult` / `ZinxServiceResult`
- Do not add platform module deps unless the user asked

## Done checklist

- [ ] `scaffold.sh` succeeded (or equivalent manual render from `assets/template`)
- [ ] No `moke-layout` / `gonew` / leftover `demo` rename steps
- [ ] `buf generate` produced `api/gen`
- [ ] `go build ./cmd/<name>/service` succeeds (after infra if runtime-tested)
- [ ] README in the new repo matches `<name>` and module path

For adding RPCs later, use `add-game-rpc`. For wiring more fx modules, use `compose-moke-modules`.
