# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

#### English | [中文](./README_CN.md)

Go toolkit for assembling game / microservice backends (gRPC, HTTP gateway, TCP, Mongo, Redis, NATS) with [uber/fx](https://github.com/uber-go/fx) — like LEGO.

## Quick start

Needs: Go (see [`go.mod`](./go.mod)), [Docker](https://docs.docker.com/get-docker/), [buf](https://buf.build/docs/installation).

```bash
# 1) From this repo: scaffold a game (uses this moke-kit checkout)
chmod +x .cursor/skills/create-game/scripts/scaffold.sh
MOKE_KIT_REPLACE="$PWD" .cursor/skills/create-game/scripts/scaffold.sh \
  --module github.com/example/arena \
  --name arena \
  --out ../arena

# 2) Start Redis / Mongo / NATS
cd ../arena
docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d

# 3) Run the service (gRPC+HTTP :8081, TCP :8888)
go run ./cmd/arena/service/main.go
```

In another terminal, call the HTTP API:

```bash
curl -s -X POST localhost:8081/v1/hello/hi \
  -H 'Content-Type: application/json' \
  -d '{"uid":"10000","message":"hello","topic":"arena"}'
```

Or use the interactive client:

```bash
cd ../arena
go run ./cmd/arena/client/main.go grpc
# then: arena → hi
```

In Cursor: `/create-game` (or ask “create a game based on moke-kit”).

## What you get

| Piece | Role |
| --- | --- |
| [`fxmain`](./fxmain) | `fxmain.Main(...)` app entry |
| [`server`](./server) | gRPC / grpc-gateway / zinx (TCP·WS·KCP) |
| [`orm`](./orm) | Mongo, Redis/Dragonfly, GORM |
| [`mq`](./mq) | NATS + local MQ |
| [`3rd`](./3rd) | Auth, IAP, Agones, … |

Reference apps: [moke-game/game](https://github.com/moke-game/game), [moke-game/platform](https://github.com/moke-game/platform). Version matrix: [COMPATIBILITY.md](./COMPATIBILITY.md).

## Next steps

- Add RPCs: Cursor skill [`add-game-rpc`](./.cursor/skills/add-game-rpc/SKILL.md)
- Wire more modules: [`compose-moke-modules`](./.cursor/skills/compose-moke-modules/SKILL.md)
- Kit self-test: `go test -race ./...`

## License

[LICENSE](./LICENSE)
