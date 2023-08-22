# Demo

## How to run

* infrastructure:
    * redis:6379
    * nats:4222
    * mongo:27017

* service:
   ```shell
     go run cmd/demo_svc/main.go
   ```
* client:
    ```shell
    go build -o demo_cli cmd/demo_cli/main.go
    ```
* run grpc client:
    ```shell
    $ ./demo_cli grpc
    $ demo hi
    ```
* run zinx client:
    ```shell
    $ ./demo_cli zinx
    $  demo hi
    ```
      