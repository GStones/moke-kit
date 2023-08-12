# moke-kit 
 一个基于uber/fx搭建的依赖框架结构套件，提供各种组建的依赖注入，以及各种组件的初始化，以及各种组件的生命周期管理。
## server
提供可支持`http` `tcp` `grpc` 的服务注入机制
 * http: 基于 grpc-gateway 代理实现
 * grpc: 基于`grpc`实现 
 * tcp： 基于`zinx`实现
 
## gorm
 * mongodb: 提供mongodb的基本adapter
 * mock: 提供mock实现(TODO)
## mq
 * nats: 提供nats的基本adapter 
 * kafka: 提供kafka的基本adapter(TODO)
 * rabbitmq: 提供rabbitmq的基本adapter(TODO)
 * mock: 提供mock实现(TODO)
## fxmain
 * 提供服务的基本创建方法
 * 提供服务的基本生命周期管理
 * 提供服务的基本依赖注入
