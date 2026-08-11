# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

#### [English](./README.md) | 中文

用 [uber/fx](https://github.com/uber-go/fx) 组装游戏 / 微服务后端的 Go 工具包（gRPC、HTTP gateway、TCP、Mongo、Redis、NATS）——像搭 LEGO。

## 快速开始

需要：Go（见 [`go.mod`](./go.mod)）、[Docker](https://docs.docker.com/get-docker/)、[buf](https://buf.build/docs/installation)。

```bash
# 1) 在本仓库根目录生成一个游戏（使用当前 moke-kit 源码）
chmod +x .cursor/skills/create-game/scripts/scaffold.sh
MOKE_KIT_REPLACE="$PWD" .cursor/skills/create-game/scripts/scaffold.sh \
  --module github.com/example/arena \
  --name arena \
  --out ../arena

# 2) 启动 Redis / Mongo / NATS
cd ../arena
docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d

# 3) 跑服务（gRPC+HTTP :8081，TCP :8888）
go run ./cmd/arena/service/main.go
```

另开终端打一枪 HTTP：

```bash
curl -s -X POST localhost:8081/v1/hello/hi \
  -H 'Content-Type: application/json' \
  -d '{"uid":"10000","message":"hello","topic":"arena"}'
```

或用交互式客户端：

```bash
cd ../arena
go run ./cmd/arena/client/main.go grpc
# 然后: game → hi
```

在 Cursor 里：`/create-game`，或直接说「基于 moke-kit 创建一个 game」。

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

- 加 RPC：Cursor Skill [`add-game-rpc`](./.cursor/skills/add-game-rpc/SKILL.md)
- 组装更多模块：[`compose-moke-modules`](./.cursor/skills/compose-moke-modules/SKILL.md)
- 开发本仓库：`go test -race ./...`

## 许可证

[LICENSE](./LICENSE)
