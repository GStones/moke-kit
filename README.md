# moke-kit
A dependency framework structure kit based on uber/fx, which provides dependency injection of various components, as well as initialization of various components, and lifecycle management of various components.

## warning
The current framework is still in the development stage and cannot be used in a production environment
## server
Provides services that support `http` `tcp` `grpc` injection mechanism
 * http: Based on grpc-gateway proxy implementation
 * grpc: Based on `grpc` implementation
 * tcp: Based on `zinx` implementation

## gorm
nosql orm framework
* mongodb: Provides basic adapter for mongodb
* mock: Provides mock implementation (TODO)
## mq
messageQueue basic adapter
* nats: Provides basic adapter for nats
* kafka: Provides basic adapter for kafka (TODO)
* rabbitmq: Provides basic adapter for rabbitmq (TODO)
* mock: Provides mock implementation (TODO)
## fxmain
Manage the startup and shutdown of all modules
* Provide basic creation method of service
* Provide basic lifecycle management of services
* Provide basic dependency injection of services
## tracing(TODO)
Provide basic tracing support
## logging
Provide `uber/zap` basic logging support

## demo 
  An example for moke-kits




