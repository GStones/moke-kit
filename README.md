# moke-kit
[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

## What is moke-kit?
moke-kit is a toolkit for building a [Domain-Driven Hexagon](https://github.com/Sairyss/domain-driven-hexagon) microservices in Go.

## Diagram

![moke-kit](./assets/moke-kit-diagram.drawio.png)

## Features
* [x] [Dependency injection](https://www.wikiwand.com/en/Dependency_injection) service/module
* [x] [Interactive client](https://github.com/GStones/moke-kit/blob/main/demo/cmd/demo_cli/main.go)
* [Server](https://github.com/GStones/moke-kit/tree/main/server):
    * [x] [gRPC](https://grpc.io/)
    * [x] HTTP[[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)]
    * [x] TCP [[zinx](https://github.com/aceld/zinx)]
    * [x] Websocket [[zinx](https://github.com/aceld/zinx)]
* [MQ](https://github.com/GStones/moke-kit/tree/main/mq):
    * [x] [nats](https://nats.io/)
    * [ ] local
    * [ ] mock
    * [ ] kafka
    * [ ] rabbitmq
    * [ ] nsq
* [Orm](https://github.com/GStones/moke-kit/tree/main/orm):
    * [x] [gorm](https://gorm.io/)
    * [x] [mongodb](https://github.com/mongodb/mongo-go-driver)
      * Atomicity and Transactions: [Set/Get](https://github.com/GStones/moke-kit/blob/main/orm/nosql/mongo/internal/driver.go#L25)
      * What is [CAS](https://www.wikiwand.com/en/Compare-and-swap)?
* [cache](https://github.com/GStones/moke-kit/blob/main/orm/nosql/diface/icache.go):
    * [x] redis
    * [x] [Cache Aside Pattern](https://blog.cdemi.io/design-patterns-cache-aside-pattern/)

## Getting started
* install gonew:
 ``` bash 
    go install golang.org/x/tools/cmd/gonew@latest
 ```
* create a new project:
 ``` bash 
    gonew github.com/gstones/moke-layout your.domain/myprog
 ```





