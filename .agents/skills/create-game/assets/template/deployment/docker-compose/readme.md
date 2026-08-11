# Local infrastructure

```bash
docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d
```

Starts Redis `:6379`, MongoDB `:27017`, and NATS JetStream `:4222`.
