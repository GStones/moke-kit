---
name: add-game-rpc
description: Add a protobuf RPC and implement it on a moke-kit game service (gRPC, grpc-gateway, optional zinx). Use when adding an API, proto method, handler, or extending game0/demo service contracts.
paths:
  - "**/api/**/*.proto"
  - "**/internal/services/**/*.go"
  - "**/pkg/modules/**/*.go"
---

# Add a game RPC (moke-kit)

Use this when the game repo already exists and the user wants a new API method. For a brand-new game project, use `create-game` first.

## Steps

1. **Proto** — edit `api/{name}/{name}.proto`:
   - Add RPC + messages
   - Add `google.api.http` annotations if the method should be on the HTTP gateway
   - Keep `package {name}.pb` and `option go_package` consistent

2. **Generate**

   ```bash
   buf generate
   ```

   Do not hand-edit `api/gen/**`.

3. **Domain** — put business logic under `internal/services/{name}/domain` (or equivalent). Keep transport handlers thin.

4. **Service handlers** — on the service struct in `internal/services/{name}`:
   - Implement the generated gRPC method
   - Ensure `RegisterWithGrpcServer` registers the server (`pb.RegisterXxxServiceServer`)
   - For gateway: `RegisterWithGatewayServer` → `RegisterXxxServiceHandlerFromEndpoint`
   - For zinx (optional): map a msg ID in `Handle` / `RegisterWithServer` with `proto.Marshal` / `Unmarshal`

5. **Providers** — only change `fx.Provide` constructors if new deps are required (DB collection, MQ, redis, settings). Still return `sfx.GrpcServiceResult` / `GatewayServiceResult` / `ZinxServiceResult`.

6. **Client / tests**
   - Extend interactive client under `internal/clients/{name}` if one exists
   - Update k6 under `tests/{name}` with the new RPC FQDN
   - Smoke: `go run ./cmd/{name}/service/main.go` + client or HTTP call

## Auth note

- Embedding `utility.WithoutAuth` skips auth for that service
- In `prod` / `prod_*` / `prod-*`, exposing gRPC/gateway without auth middleware fails closed — do not “fix” that by disabling checks; wire real auth (e.g. platform `AuthMiddlewareModule`) instead

## Done when

- `buf generate` clean
- New RPC works on the intended transports
- No edits under `api/gen` except via buf
