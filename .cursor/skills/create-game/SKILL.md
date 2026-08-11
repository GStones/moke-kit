---
name: create-game
description: Scaffold a new game server based on moke-kit (gonew moke-layout or fork game0). Use when the user asks to create a game, new game project, bootstrap from moke-kit, clone moke-layout, or start a game server template.
---

# Create a game based on moke-kit

moke-kit is LEGO bricks (`fxmain`, server, orm, mq). A game is a composed app. There is no in-repo generator — scaffold from the official layout, then rename and wire.

## Mental model

```text
moke-kit          → infra modules (always via fxmain.Main)
moke-layout/game  → game server template (one named service)
platform (opt)    → shared services imported as fx modules
```

Reference implementations:

- Layout template: https://github.com/gstones/moke-layout (service name `demo`)
- Full game example: https://github.com/moke-game/game (service name `game0`)
- Platform modules: https://github.com/moke-game/platform

When this workspace has `/agent/repos/game` or `/agent/repos/platform`, prefer reading those trees over guessing.

## Gather requirements first

Ask only if missing:

1. **Module path** — e.g. `github.com/acme/arena`
2. **Service name** — short, lowercase, no spaces (e.g. `arena`; template defaults are `demo` / `game0`)
3. **Scope** — minimal game only, or compose platform modules (auth, profile, mail, …)
4. **Transports** — gRPC / HTTP gateway / TCP(zinx); default = all three like `AllModule`

## Path A — Greenfield (recommended)

Documented in moke-kit README:

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/gstones/moke-layout github.com/<org>/<game>
cd <game>
```

Then rename `demo` → `{name}` across the tree. Read [references/template-map.md](references/template-map.md) for the file map.

### Required renames

- Paths: `cmd/demo` → `cmd/{name}`, `api/demo`, `internal/services/demo`, `internal/clients/demo`, `tests/demo`, `pkg/modules/*`
- Proto: `package demo.pb` → `{name}.pb`, service/messages, `go_package`, HTTP routes
- Go imports / fx module vars / env defaults (`GAME_URL`, `DB_NAME`, key namespaces)
- `buf.yaml` `name:` if publishing to a new BSR module
- k6 tests: fix RPC FQDN to the new service (game0’s k6 may still say `game.pb.DemoService` — correct it)

### Wire entrypoint

`cmd/{name}/service/main.go` should look like:

```go
fxmain.Main(
    // infra beyond AppModule (include only what you use)
    mfx.NatsModule,
    mfx.LocalModule,
    ofx.RedisCacheModule,

    // game modules
    modules.AllModule, // or GrpcModule / HttpModule / TcpModule

    // optional platform shared services
    // auth.AuthAllModule, profile.ProfileModule, ...
)
```

`fxmain.Main` already injects `AppModule` (server + orm + logging + mq settings) and launches `ServiceBinder`.

### Codegen and run

```bash
buf generate

docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d
go run ./cmd/{name}/service/main.go
```

Defaults: HTTP/gRPC `:8081`, TCP `:8888`, Mongo/Redis/NATS on localhost (see moke-kit orm/mq/server READMEs).

Smoke:

```bash
go build -o {name} ./cmd/{name}/client/main.go
./{name} grpc   # or tcp; HTTP via Postman localhost:8081
```

Docker image (template Dockerfile):

```bash
docker buildx build -t <registry>/<name>:latest \
  --build-arg APP_NAME={name} \
  --build-arg GIT_PWD=<git_token_if_private_deps> \
  -f ./build/package/docker/Dockerfile . --push
```

## Path B — Copy from existing game0 in workspace

If `/agent/repos/game` is available and the user wants the platform-composed template:

1. Copy the `game0` trees listed in [references/template-map.md](references/template-map.md)
2. Rename packages/proto/modules/env as in Path A
3. Trim platform imports in `main.go` unless needed
4. `buf generate`, start infra, run service

Do **not** put game-specific logic into platform; only add platform services when they are shared middleware.

## Hard rules (do not invent)

- Register services by implementing `siface.IGrpcService` / `IGatewayService` / `IZinxService` and returning `sfx.GrpcServiceResult` / `GatewayServiceResult` / `ZinxServiceResult` from `fx.Provide`
- Export LEGO pieces from `pkg/modules` (`GrpcModule`, `HttpModule`, `TcpModule`, `AllModule`, client module)
- Keep settings in `pkg/dfx` via `utility.Load` / envconfig
- Persist with `nosql.DocumentBase` + key helpers when following the template DAO pattern
- Embed `utility.WithoutAuth` only when intentionally skipping auth (prod fails closed without auth middleware)

## Done checklist

- [ ] Module path and `{name}` consistent in go.mod, imports, proto, paths
- [ ] `buf generate` succeeds; no stale `demo`/`game0` references in new code
- [ ] `fxmain.Main` wires needed moke-kit + optional platform modules
- [ ] Infra compose up; `go run ./cmd/{name}/service/main.go` starts
- [ ] Client or HTTP smoke works for at least one RPC
- [ ] README run/build snippets use `{name}`

If the user only wants a new RPC inside an existing game, use the `add-game-rpc` skill instead.
