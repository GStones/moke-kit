# moke-kit
A dependency framework structure kit based on uber/fx, which provides dependency injection of various components, as well as initialization of various components, and lifecycle management of various components.

## Warning
The current framework is still in the development stage and cannot be used in a production environment

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
    * [ ] nats
* [Orm](https://github.com/GStones/moke-kit/tree/main/orm):
    * [x] [gorm](https://gorm.io/)
    * [x] [mongodb](https://github.com/mongodb/mongo-go-driver)
      * Atomicity and Transactions: [Set/Get](https://github.com/GStones/moke-kit/blob/main/orm/nosql/mongo/internal/driver.go#L25)
      * What is [CAS](https://www.wikiwand.com/en/Compare-and-swap)?
* [cache](https://github.com/GStones/moke-kit/blob/main/orm/nosql/diface/icache.go):
    * [x] [Cache Aside Pattern](https://blog.cdemi.io/design-patterns-cache-aside-pattern/)

## Getting started
 * [Introduction](https://github.com/GStones/moke-kit/wiki/Introduction)
 * [Demo](./demo)



