# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

#### English | [中文](./README_CN.md)

Go toolkit for assembling game / microservice backends (gRPC, HTTP gateway, TCP, Mongo, Redis, NATS) with [uber/fx](https://github.com/uber-go/fx) — like LEGO.

## Quick start

Use any AI coding agent that supports [Agent Skills](https://agentskills.io) (Cursor, Claude Code, Codex, Windsurf, …):

1. Open this repository in your AI tool.
2. Ask it to create a game, for example:
   - *create a game based on moke-kit*
   - *use the create-game skill*
   - `/create-game` (if your tool supports skill slash commands)
3. Give a Go module path and service name (e.g. `github.com/example/arena`, `arena`).

The [`create-game`](./.agents/skills/create-game/SKILL.md) skill scaffolds the game, starts local Redis/Mongo/NATS, and runs the service.

Needs on the machine: Go (see [`go.mod`](./go.mod)), [Docker](https://docs.docker.com/get-docker/), [buf](https://buf.build/docs/installation).

## What you get

| Piece | Role |
| --- | --- |
| [`fxmain`](./fxmain) | `fxmain.Main(...)` app entry |
| [`server`](./server) | gRPC / grpc-gateway / zinx (TCP·WS·KCP) |
| [`orm`](./orm) | Mongo, Redis/Dragonfly, GORM |
| [`mq`](./mq) | NATS + local MQ |
| [`3rd`](./3rd) | Auth, IAP, Agones, … |

Skills live under [`.agents/skills/`](./.agents/skills) ([Agent Skills](https://agentskills.io) layout). Reference apps: [moke-game/game](https://github.com/moke-game/game), [moke-game/platform](https://github.com/moke-game/platform). Version matrix: [COMPATIBILITY.md](./COMPATIBILITY.md).

## Next steps

- Add RPCs — ask your AI / [`add-game-rpc`](./.agents/skills/add-game-rpc/SKILL.md)
- Wire more modules — ask your AI / [`compose-moke-modules`](./.agents/skills/compose-moke-modules/SKILL.md)
- Kit self-test: `go test -race ./...`

## License

[LICENSE](./LICENSE)
