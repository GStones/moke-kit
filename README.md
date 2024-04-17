# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

## What is moke-kit?

moke-kit is a toolkit for building a [Domain-Driven Hexagon](https://github.com/Sairyss/domain-driven-hexagon)
microservices/monolithic in Go. You can develop as a monolithic service and deploy it as microservice. 
Also like play with lego, you can assemble the service as you like.

## Diagram

![moke-kit](./assets/moke-kit-diagram.drawio.png)

## Features
* [Server](https://github.com/GStones/moke-kit/tree/main/server):
    * [x] [gRPC](https://grpc.io/)
    * [x] HTTP[[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)]
    * [x] TCP [[zinx](https://github.com/aceld/zinx)]
    * [x] Websocket [[zinx](https://github.com/aceld/zinx)]
* [MQ](https://github.com/GStones/moke-kit/tree/main/mq):
    * [x] [nats](https://nats.io/)
* [Orm](https://github.com/GStones/moke-kit/tree/main/orm):
    * [x] [gorm](https://gorm.io/)
    * [x] [mongodb](https://github.com/mongodb/mongo-go-driver)
* [Cache](https://github.com/GStones/moke-kit/tree/main/orm/nosql/cache):
    * [x] redis
    * [x] [dragonfly](https://github.com/dragonflydb/dragonfly)
* [Third-party](https://github.com/GStones/moke-kit/tree/main/3rd):
    * [x] [IAP](https://github.com/awa/go-iap)
    * [x] [Agones](https://agones.dev/site/)

## Getting started

* install gonew:

 ``` bash 
    go install golang.org/x/tools/cmd/gonew@latest
 ```

* create a new project:

 ``` bash 
    gonew github.com/gstones/moke-layout your.domain/myprog
 ```