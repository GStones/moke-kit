# ORM

Database Adapter

## Modules:

* [MongoPureModule](https://github.com/mongodb/mongo-go-driver): MongoDB driver for Go.
* [GormModule](https://gorm.io/): The fantastic ORM library for Golang, aims to be developer friendly.
* DocumentStoreModule: Document store adapter, now support MongoDB.
* RedisModule: redis go client, provide redis(db0) and cache(db1) .
* RedisCacheModule: redis cache adapter implement.

## Environment Variables

| ENV               | Description       | Default                   |
|-------------------|-------------------|---------------------------|
| DATABASE_URL      | Database host     | mongodb://localhost:27017 |
| DATABASE_USER     | Database username | ""                        |
| DATABASE_PASSWORD | Database password | ""                        |
| CACHE_URL         | Cache host        | redis://localhost:6379    |
| CACHE_USER        | Cache username    | ""                        |
| CACHE_PASSWORD    | Cache password    | ""                        |

