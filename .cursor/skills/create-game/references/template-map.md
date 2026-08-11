# Template file map

Use when renaming `demo` (moke-layout) or `game0` (moke-game/game) to `{name}`.

## Paths to copy / rename

| From (game0) | To |
| --- | --- |
| `cmd/game0/service/main.go` | `cmd/{name}/service/main.go` |
| `cmd/game0/client/main.go` | `cmd/{name}/client/main.go` |
| `api/game0/game0.proto` | `api/{name}/{name}.proto` |
| `api/gen/game0/` | regenerate via `buf generate` (do not hand-copy) |
| `internal/services/game0/**` | `internal/services/{name}/**` |
| `internal/clients/game0/**` | `internal/clients/{name}/**` |
| `pkg/modules/game0_module.go` | `pkg/modules/{name}_module.go` |
| `pkg/dfx/game0_setting.go` | `pkg/dfx/{name}_setting.go` (or keep generic settings) |
| `pkg/dfx/game0_client.go` | `pkg/dfx/{name}_client.go` |
| `pkg/dfx/auth.go` | keep / adapt if using auth middleware |
| `tests/game0/**` | `tests/{name}/**` |
| `buf.yaml` | update `name:` for new BSR module |
| `go.mod` | new module path; depend on `github.com/gstones/moke-kit` |

moke-layout uses `demo` instead of `game0`; same relative layout.

## Proto shape (game0)

```protobuf
syntax = "proto3";
package {name}.pb;
import "google/api/annotations.proto";
option go_package = "{name}/api;pb";

service {Name}Service {
  rpc Hi(HiRequest) returns (HiResponse) {
    option (google.api.http) = { post: "/v1/hello/hi" body: "*" };
  }
}
```

Generated Go import pattern:

```go
pb "github.com/<org>/<game>/api/gen/{name}/api"
```

## pkg/modules LEGO pieces

Export (names can stay descriptive):

- `GrpcModule` → settings + `GrpcService`
- `HttpModule` → settings + grpc + gateway providers
- `TcpModule` → settings + `TcpService` / zinx
- `AllModule` → all transports
- `GrpcClientModule` → client provider for CLI / other services

## fx Provide results

| Transport | Result type | Field |
| --- | --- | --- |
| gRPC | `sfx.GrpcServiceResult` | `GrpcService` |
| Gateway | `sfx.GatewayServiceResult` | `GatewayService` |
| Zinx TCP/WS/KCP | `sfx.ZinxServiceResult` | `ZinxService` |

## Infra env defaults (moke-kit)

| ENV | Default |
| --- | --- |
| `PORT` | `8081` |
| `ZINX_TCP_PORT` | `8888` |
| `DATABASE_URL` | `mongodb://localhost:27017` |
| `CACHE_URL` | `redis://localhost:6379` |
| `NATS_URL` | `nats://localhost:4222` |
| `DEPLOYMENT` | `local` |
| `GAME_URL` (game settings) | `localhost:8081` |
| `DB_NAME` (game settings) | `game` |

## Platform modules (optional imports)

Import path pattern:

```go
github.com/moke-game/platform/services/<svc>/pkg/module
```

Common assemblers (see game0 `main.go`): `AuthAllModule`, `ProfileModule`, `MailModule`, `AnalyticsModule`, `KnapsackModule`, `PartyModule`, `BuddyModule`, `LeaderboardModule`, `ChatModule`, `MatchmakingModule`.
