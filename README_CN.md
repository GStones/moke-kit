# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

#### [English](./README.md) | 中文

用 [uber/fx](https://github.com/uber-go/fx) 组装游戏 / 微服务后端的 Go 工具包（gRPC、HTTP gateway、TCP、Mongo、Redis、NATS）——像搭 LEGO。

## 快速开始

任意支持 [Agent Skills](https://agentskills.io) 的 AI 编程助手均可（Cursor、Claude Code、Codex、Windsurf 等）：

1. 在 AI 工具中打开本仓库。
2. 让它创建游戏，例如：
   - *基于 moke-kit 创建一个 game*
   - *使用 create-game skill*
   - `/create-game`（若工具支持 Skill 斜杠命令）
3. 提供 Go module 路径和服务名（例如 `github.com/example/arena`、`arena`）。

[`create-game`](./.agents/skills/create-game/SKILL.md) Skill 会脚手架、拉起本地 Redis/Mongo/NATS，并启动服务。

本机需具备：Go（见 [`go.mod`](./go.mod)）、[Docker](https://docs.docker.com/get-docker/)、[buf](https://buf.build/docs/installation)。

## 模块一览

| 目录 | 作用 |
| --- | --- |
| [`fxmain`](./fxmain) | `fxmain.Main(...)` 应用入口 |
| [`server`](./server) | gRPC / grpc-gateway / zinx（TCP·WS·KCP） |
| [`orm`](./orm) | Mongo、Redis/Dragonfly、GORM |
| [`mq`](./mq) | NATS + 本地 MQ |
| [`3rd`](./3rd) | Auth、IAP、Agones 等 |

Skill 位于 [`.agents/skills/`](./.agents/skills)（[Agent Skills](https://agentskills.io) 标准布局）。参考项目：[moke-game/game](https://github.com/moke-game/game)、[moke-game/platform](https://github.com/moke-game/platform)。版本矩阵：[COMPATIBILITY.md](./COMPATIBILITY.md)。

## 接下来

- 加 RPC — 让 AI 执行 / [`add-game-rpc`](./.agents/skills/add-game-rpc/SKILL.md)
- 组装更多模块 — 让 AI 执行 / [`compose-moke-modules`](./.agents/skills/compose-moke-modules/SKILL.md)
- 开发本仓库：`go test -race ./...`

## 许可证

[LICENSE](./LICENSE)
