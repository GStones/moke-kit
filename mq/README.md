# MQ

Message Queue Adapter

## [NATS](https://nats.io/):

A simple, secure and high performance open source messaging system for cloud native
applications, IoT messaging, and microservices architectures.

## Local(Channel):

A simple channel based message queue for local message passing.

## Modules:

* `Module`: mq modules init

## Environment Variables:

| ENV                                | Description                                      | Default               |
|------------------------------------|--------------------------------------------------|-----------------------|
| NATS_URL                           | nats host                                        | nats://localhost:4222 |
| CHANNEL_BUFFER_SIZE                | local channel buffer size                        | 1024                  |
| PERSISTENT                         | local channel persistent                         | false                 |
| BLOCK_PUBLISH_UNTIL_SUBSCRIBER_ACK | local channel block publish until subscriber ack | false                 |





