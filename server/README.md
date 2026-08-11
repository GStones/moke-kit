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

| ENV            | Description                                                                 | Default                           |
|----------------|-----------------------------------------------------------------------------|-----------------------------------|
| MTLS_ENABLE    | enable mTLS (server TLS + require client certs). Implies TLS for cmux.       | false                             |
| TLS_ENABLE     | enable server TLS for grpc/http (cmux). Does **not** require client certs.  | false                             |
| TCP_TLS_ENABLE | enable TCP tls                                                              | false                             |
| CLIENT_CA_CERT | client ca cert path (required when `MTLS_ENABLE=true`)                      | "./configs/tls-client/ca.crt"     |
| CLIENT_CERT    | client cert path (used by mTLS clients / gateway dial)                      | "./configs/tls-client/tls.crt"    |
| CLIENT_KEY     | client key path (used by mTLS clients / gateway dial)                       | "./configs/tls-client/tls.key"    |
| SERVER_CA_CERT | server ca cert path (used by clients verifying the server)                  | "./configs/tls-server/ca.crt"     |
| SERVER_CERT    | server cert path                                                            | "./configs/tls-server/tls.crt"    |
| SERVER_KEY     | server key path                                                             | "./configs/tls-server/tls.key"    |
| SERVER_NAME    | server name                                                                 | ""                                |

### CORS / Auth:

| ENV                | Description                                                                 | Default |
|--------------------|-----------------------------------------------------------------------------|---------|
| CORS_ALLOW_ORIGINS | comma-separated allowed Origins; empty disables CORS; `*` allows any origin | ""      |

Notes:
* Auth middleware is optional in local/dev. In **prod** deployments (`prod`, `prod_*`, `prod-*`), missing auth middleware fails closed (Unauthenticated).
* Gateway dials gRPC with TLS/mTLS credentials when `TLS_ENABLE` or `MTLS_ENABLE` is set.


