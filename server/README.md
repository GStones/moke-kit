# Server

## [Cmux]( https://github.com/soheilhy/cmux)

* Support multiple protocols on the same port
    * HTTP/1.1
    * HTTP/2
    * gRPC
    * WebSocket
    * etc.

## Srpc

* [grpc](https://github.com/grpc/grpc-go)
    * Support gRPC service manager .
    * port=>[$PORT](https://github.com/GStones/moke-kit/blob/main/server/pkg/sfx/settings_module.go#L21)
* [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
    * Support Http service .
    * port=>[$PORT](https://github.com/GStones/moke-kit/blob/main/server/pkg/sfx/settings_module.go#L21)
* [grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)
    * Support authentication, logging, monitoring, recover panic etc.

## [Zinx](https://github.com/aceld/zinx)

* Support TCP/WebSocket server manager
* TCP port => [$ZINX_TCP_PORT](https://github.com/GStones/moke-kit/blob/main/server/pkg/sfx/settings_module.go#L22)
* Websocket port => [$ZINX_WS_PORT](https://github.com/GStones/moke-kit/blob/main/server/pkg/sfx/settings_module.go#L23)
