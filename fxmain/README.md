# MAIN

Bind,Manager and Run all services with [uber fx](https://github.com/uber-go/fx)

## Modules

* `AppModule`: app module init with  `server`, `orm`, `mq`, `logging` modules and inject custom modules

## Environment Variables

| ENV        | Description                                                                                                          | Default |
|------------|----------------------------------------------------------------------------------------------------------------------|---------|
| APP_NAME   | Application name                                                                                                     | app     |
| APP_ID     | Application id                                                                                                       | app     |
| DEPLOYMENT | local,dev,prod <br/> prod: <br/>  logger change to production mod <br/> grpc will add recovery panic with middleware | local   |
| VERSION    | Application version                                                                                                  | 0.0.1   |





