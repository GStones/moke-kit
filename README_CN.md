# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

#### [English](./README.md) | 中文

用 [uber/fx](https://github.com/uber-go/fx) 组装游戏 / 微服务后端的 Go 工具包（gRPC、HTTP gateway、TCP、Mongo、Redis、NATS）——像搭 LEGO。

## 快速开始

1. 用 [Cursor](https://cursor.com) 打开本仓库。
2. 在 Agent 对话里执行 **`/create-game`**（或直接说：*基于 moke-kit 创建一个 game*）。
3. 告诉 Agent 你的 Go module 路径和服务名（例如 `github.com/example/arena`、`arena`）。

[`create-game`](./.cursor/skills/create-game/SKILL.md) Skill 会自动脚手架、拉起本地 Redis/Mongo/NATS，并启动服务。

本机需具备：Go（见 [`go.mod`](./go.mod)）、[Docker](https://docs.docker.com/get-docker/)、[buf](https://buf.build/docs/installation)。

## 模块一览

| 目录 | 作用 |
| --- | --- |
| [`fxmain`](./fxmain) | `fxmain.Main(...)` 应用入口 |
| [`server`](./server) | gRPC / grpc-gateway / zinx（TCP·WS·KCP） |
| [`orm`](./orm) | Mongo、Redis/Dragonfly、GORM |
| [`mq`](./mq) | NATS + 本地 MQ |
| [`3rd`](./3rd) | Auth、IAP、Agones 等 |

参考项目：[moke-game/game](https://github.com/moke-game/game)、[moke-game/platform](https://github.com/moke-game/platform)。版本矩阵：[COMPATIBILITY.md](./COMPATIBILITY.md)。

## 接下来

- 加 RPC：`/add-game-rpc` · [`add-game-rpc`](./.cursor/skills/add-game-rpc/SKILL.md)
- 组装更多模块：`/compose-moke-modules` · [`compose-moke-modules`](./.cursor/skills/compose-moke-modules/SKILL.md)
- 开发本仓库：`go test -race ./...`

## 许可证

[LICENSE](./LICENSE)
