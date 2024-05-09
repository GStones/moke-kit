# Server

## Modules:

* `AuthService`: GRpc Authentication service, if you want to use it, you need to implement the `AuthService` interface
  and inject it. Per request will be checked by the `AuthService` service .
* `ConnectionMuxModule`: GRpc,Http will listen on the same port, and the connection will be handled by
  the `ConnectionMuxModule` .
* `OTelModule`: grpc open telemetry module:  https://github.com/open-telemetry/opentelemetry-go.
* `ServiceBinder`: Bind all injected services: grpc,grpc-gateway, tcp, otel modules.

## Environment Variables:

### Basic:

| ENV                 | Description              | Default |
|---------------------|--------------------------|---------|
| PORT                | http/grpc listen port    | 8081    |
| ZINX_TCP_PORT       | tcp/udp listen port      | 8888    |
| ZINX_WS_PORT        | ws listen port           | ""      |
| MAX_PACKET_SIZE     | zinx max packet size     | 4096    |
| WORKER_POOL_SIZE    | zinx worker pool size    | 64      |
| MAX_WORKER_TASK_LEN | zinx max worker task len | 1024    |
| MAX_MSG_CHAN_LEN    | zinx max msg chan len    | 1024    |
| TIMEOUT             | tcp heartbeat timeout(s) | 10      |
| RATE_LIMIT          | rate limit per second    | 1000    |
| OTEL_ENABLE         | enable open telemetry    | false   |

### TLS:

| ENV            | Description               | Default                           |
|----------------|---------------------------|-----------------------------------|
| MTLS_ENABLE    | enable mTLS mod           | false                             |
| TLS_ENABLE     | enable server tls         | false                             |
| TCP_TLS_ENABLE | enable TCP tls            | false                             |
| CLIENT_CA_CERT | client ca cert path(mTls) | "./configs/tls-client/ca.crt"     |
| CLIENT_CERT    | client cert path(mTls)    | "./configs/tls-client/client.crt" |
| CLIENT_KEY     | client key path(mTls)     | "./configs/tls-client/client.key" |
| SERVER_CA_CERT | server ca cert path       | "./configs/tls-server/ca.crt"     |
| SERVER_CERT    | server cert path          | "./configs/tls-server/server.crt" |
| SERVER_KEY     | server key path           | "./configs/tls-server/server.key" |
| SERVER_NAME    | sever name                | ""                                |   


