# MAIN

Bind, manage, and run services with [uber fx](https://github.com/uber-go/fx).

## Entries

| Function | Graph | Use |
| --- | --- | --- |
| `fxmain.Main(opts...)` | `AppModule` = settings + logging + server + orm + mq | Typical game / API binary |
| `fxmain.Core(opts...)` | `CoreModule` = settings + logging | Workers, CLIs, or an explicit LEGO stack |

`server.Module` binds gRPC / gateway / zinx itself. A `Core` process that omits it does not listen on `PORT`.

## Modules

* `CoreModule`: app settings + logging
* `AppModule`: `CoreModule` plus `server`, `orm`, and `mq`

## Environment Variables

| ENV        | Description                                                                                                            | Default |
|------------|------------------------------------------------------------------------------------------------------------------------|---------|
| APP_NAME   | Application name                                                                                                       | app     |
| APP_ID     | Application id                                                                                                         | app     |
| DEPLOYMENT | local,dev,prod <br/> you can customize it as your need <br/>local_{name} = local, dev_{name} = dev, prod_{name} = prod | local   |
| VERSION    | Application version                                                                                                    | 0.0.1   |





