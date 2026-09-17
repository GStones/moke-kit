---
name: compose-moke-modules
description: Compose moke-kit and platform uber/fx modules into fxmain.Main or fxmain.Core (LEGO assembly). Use when wiring infra (NATS, Redis cache, Mongo), choosing Grpc/Http/Tcp modules, or importing platform auth/profile/mail into a game or service entrypoint.
paths:
  - "**/cmd/**/main.go"
  - "**/pkg/modules/**/*.go"
  - "**/pkg/**/module*.go"
---

# Compose moke-kit modules

Assemble services like LEGO. Prefer existing modules over new globals or manual server bootstrap.

## Choose an entry

| Entry | Graph | Use |
| --- | --- | --- |
| `fxmain.Main(opts...)` | `AppModule` = settings + logging + server + orm + mq | Game / API binaries |
| `fxmain.Core(opts...)` | `CoreModule` = settings + logging only | Workers, CLIs, or explicit stacks |

`server.Module` binds gRPC / gateway / zinx itself. Omit it from `Core` when the process should not listen.

## What AppModule already provides

`fxmain.Main(opts...)` always includes `module.AppModule`:

| Area | Module contents |
| --- | --- |
| App settings | `APP_NAME`, `APP_ID`, `DEPLOYMENT`, `VERSION` |
| Server | ports, TLS/mTLS, cmux, OTel settings + binder |
| ORM | Mongo document store, Redis (GORM is opt-in `ofx.GormModule`) |
| Logging | logging module |
| MQ settings | mq setting module (still need a concrete MQ provider) |

## Common extra infra

| Need | Module | Pass into |
| --- | --- | --- |
| NATS JetStream MQ | `mfx.NatsModule` | Main or Core + `mq.Module` |
| In-process MQ | `mfx.LocalModule` | Main or Core + `mq.Module` |
| Redis cache-aside | `ofx.RedisCacheModule` | Main or Core + `orm`/`RedisModule` |
| GORM | `ofx.GormModule` + a `Dialector` | opt-in |
| Agones / IAP | modules under `moke-kit/3rd/...` | either |

MQ topics: `nats://...` and `local://...` only. `kafka://` / `nsq://` are unsupported.

Import paths:

```go
github.com/gstones/moke-kit/fxmain
github.com/gstones/moke-kit/mq/pkg/mfx
github.com/gstones/moke-kit/orm/pkg/ofx
github.com/gstones/moke-kit/server/pkg/module
```

Thin worker example:

```go
fxmain.Core(
    mq.Module,
    mfx.NatsModule,
    myWorker.Module,
)
```

## Game service modules

From the game template `pkg/modules`:

| Module | Transports |
| --- | --- |
| `GrpcModule` | gRPC (+ optional auth) |
| `HttpModule` | gRPC + gateway |
| `TcpModule` | zinx TCP |
| `AllModule` | all of the above |
| `GrpcClientModule` | outbound gRPC client |

Only enable `dfx.AuthModule` / platform auth middleware when you intend authenticated APIs.

## Platform shared modules

Pattern:

```go
import auth "github.com/moke-game/platform/services/auth/pkg/module"

fxmain.Main(
    modules.AllModule,
    auth.AuthAllModule,
    // profile.ProfileModule, mail.MailModule, ...
)
```

Each platform service exports focused modules (`XxxModule`, `XxxClientModule`, middleware variants). Import the smallest set that matches the need.

Single platform binary example: `fxmain.Main(ofx.RedisCacheModule, module.AuthModule)`.

## Adding a new composable module

1. Implement providers in `pkg/<x>fx` or `internal` returning `sfx.*ServiceResult` or plain deps
2. Export `fx.Module("name", ...)` from `pkg/module` or game `pkg/modules`
3. Pass that module into `fxmain.Main` or `fxmain.Core` at the desired binary
4. Register lifecycle cleanup via fx when starting listeners/subscribers

## Anti-patterns

- Calling `grpc.NewServer` / listening ports outside moke-kit server modules
- Forgetting `mfx.NatsModule` or `mfx.LocalModule` while code injects `miface.MessageQueue`
- Adding every platform module “just in case” — keep Main/Core minimal
- Putting game-only logic into platform repos
- Using `fxmain.Main` for a process that only needs MQ or ORM
