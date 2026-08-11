---
name: compose-moke-modules
description: Compose moke-kit and platform uber/fx modules into fxmain.Main (LEGO assembly). Use when wiring infra (NATS, Redis cache, Mongo), choosing Grpc/Http/Tcp modules, or importing platform auth/profile/mail into a game or service entrypoint.
paths:
  - "**/cmd/**/main.go"
  - "**/pkg/modules/**/*.go"
  - "**/pkg/**/module*.go"
---

# Compose moke-kit modules

Assemble services like LEGO. Prefer existing modules over new globals or manual server bootstrap.

## What AppModule already provides

`fxmain.Main(opts...)` always includes `module.AppModule`:

| Area | Module contents |
| --- | --- |
| App settings | `APP_NAME`, `APP_ID`, `DEPLOYMENT`, `VERSION` |
| Server | ports, TLS/mTLS, cmux, OTel settings |
| ORM | Mongo document store, Redis, GORM drivers |
| Logging | logging module |
| MQ settings | mq setting module (still need a concrete MQ provider) |

Then it `Invoke`s `ServiceBinder`, which binds all group-provided gRPC / gateway / zinx services.

## Common extra infra (pass into Main)

| Need | Module |
| --- | --- |
| NATS JetStream MQ | `mfx.NatsModule` |
| In-process MQ | `mfx.LocalModule` |
| Redis cache-aside | `ofx.RedisCacheModule` |
| Agones / IAP | modules under `moke-kit/3rd/...` |

Import paths:

```go
github.com/gstones/moke-kit/fxmain
github.com/gstones/moke-kit/mq/pkg/mfx
github.com/gstones/moke-kit/orm/pkg/ofx
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
3. Pass that module into `fxmain.Main` at the desired binary
4. Register lifecycle cleanup via fx when starting listeners/subscribers

## Anti-patterns

- Calling `grpc.NewServer` / listening ports outside moke-kit server modules
- Forgetting `mfx.NatsModule` or `mfx.LocalModule` while code injects `miface.MessageQueue`
- Adding every platform module “just in case” — keep Main minimal
- Putting game-only logic into platform repos
