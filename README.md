# moke-kit

[![Go Report Card](https://goreportcard.com/badge/github.com/gstones/moke-kit)](https://goreportcard.com/report/github.com/gstones/moke-kit)
[![Go Reference](https://pkg.go.dev/badge/github.com/GStones/moke-kit.svg)](https://pkg.go.dev/github.com/GStones/moke-kit)
[![Release](https://img.shields.io/github/v/release/gstones/moke-kit.svg?style=flat-square)](https://github.com/GStones/moke-kit)

## What is moke-kit?

moke-kit is a toolkit for building a [Domain-Driven Hexagon](https://github.com/Sairyss/domain-driven-hexagon)
microservices/monolithic in Go. You can develop as a monolithic service and deploy it as a microservice.
Just like building with LEGO, you can assemble the service as you like.

## Diagram

![moke-kit](./assets/moke-kit-diagram.drawio.png)

## Features

* Inversion of control with [uber/fx](https://github.com/uber-go/fx),assemble your service as you like.
* Builtin TLS,mTLS to build [Zero Trust security](https://www.wikiwand.com/en/Zero_trust_security_model).
* Builtin middlewares (rate limit, open telemetry, auth override,logging, panic recovery, etc.).
* Builtin [Cache-Aside pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside) for orm and
  nosql.
* Builtin [Compare-and-swap](https://www.wikiwand.com/en/Compare-and-swap) to ensure db update consistency.
* Command client to interact with the server for independent testing.
* One command to generate proto, grpc, gateway, swagger and client code with [buf](https://buf.build/).

## Builtin Kits

* [Servers](https://github.com/GStones/moke-kit/tree/main/server):
    * [gRPC](https://grpc.io/)
    * HTTP[[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)]
    * TCP [[zinx](https://github.com/aceld/zinx)]
    * Websocket [[zinx](https://github.com/aceld/zinx)]
* [MQ](https://github.com/GStones/moke-kit/tree/main/mq):
    * [nats](https://nats.io/)
* [Orm](https://github.com/GStones/moke-kit/tree/main/orm):
    * [gorm](https://gorm.io/)
    * [mongodb](https://github.com/mongodb/mongo-go-driver)
* [Cache](https://github.com/GStones/moke-kit/tree/main/orm/nosql/cache):
    * redis
    * [dragonfly](https://github.com/dragonflydb/dragonfly)
* [Third Party](https://github.com/GStones/moke-kit/tree/main/3rd):
    * [IAP](https://github.com/awa/go-iap): Verifies the purchase receipt via AppStore, GooglePlayStore or Amazon
      AppStore.
    * [Agones](https://agones.dev/site/):  Host, Run and Scale dedicated game servers on Kubernetes.

## Getting started

* install gonew:

 ``` bash 
    go install golang.org/x/tools/cmd/gonew@latest
 ```

* create a new project:

 ``` bash 
    gonew github.com/gstones/moke-layout your.domain/myprog
 ```