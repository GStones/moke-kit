---
name: create-game
description: Scaffold a new game server from moke-kit templates, start local infra, and run the service. Use when the user asks to create a game, new game project, bootstrap from moke-kit, /create-game, or generate a game server scaffold.
---

# Create a game based on moke-kit

End-to-end for the user: **scaffold → start infra → run service → smoke**. Do **not** use `gonew`, `moke-layout`, or a `demo` rename flow. Do **not** stop after generating files — get the service running unless the user only wants files.

## Gather requirements

Ask only if missing:

1. **Module path** — e.g. `github.com/acme/arena`
2. **Service name** — lowercase (`arena`)
3. **Output directory** — default `../{name}` beside the moke-kit checkout (or `./{name}` if the user prefers)
4. **Platform stubs?** — default no

## 1) Scaffold

From the **moke-kit repository root**:

```bash
chmod +x .cursor/skills/create-game/scripts/scaffold.sh

# Prefer this checkout of moke-kit when scaffolding from source
MOKE_KIT_REPLACE="$PWD" .cursor/skills/create-game/scripts/scaffold.sh \
  --module <module> \
  --name <name> \
  --out <out>
```

Optional: `--with-platform`, `--buf-module buf.build/<org>/<name>`.

The script copies [`assets/template/`](assets/template/), runs `buf generate`, then `go mod tidy`. See [references/template-map.md](references/template-map.md) only if the script cannot run.

## 2) Start infra and run (required unless user declines)

```bash
cd <out>
docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d
go run ./cmd/<name>/service/main.go
```

Run the service in the background / a dedicated terminal so it keeps listening. Defaults: HTTP/gRPC `:8081`, TCP `:8888`.

## 3) Smoke

```bash
curl -s -X POST localhost:8081/v1/hello/hi \
  -H 'Content-Type: application/json' \
  -d '{"uid":"10000","message":"hello","topic":"<name>"}'
```

Tell the user the output path, ports, and the curl/client commands.

## Hard rules

- Never call `gonew` or depend on `github.com/gstones/moke-layout`
- Never leave `demo` / `game0` names unless chosen as `--name`
- Do not hand-edit `api/gen/**`
- Register via `sfx.GrpcServiceResult` / `GatewayServiceResult` / `ZinxServiceResult`
- Do not add platform module deps unless asked

## Done checklist

- [ ] `scaffold.sh` succeeded
- [ ] Infra compose is up
- [ ] Service is running
- [ ] Smoke curl (or equivalent) succeeded
- [ ] User knows where the project is and how to call it

Later: `add-game-rpc`, `compose-moke-modules`.
