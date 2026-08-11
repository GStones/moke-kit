# Compatibility

Use this matrix to track which `moke-kit` versions are validated with the
reference platform and game repositories.

| moke-kit version | platform | game | Notes |
| --- | --- | --- | --- |
| `v1.0.5-0.20260811094419-bcdfe55515cd` (#224) | [platform#24](https://github.com/moke-game/platform/pull/24)+ | [game#19](https://github.com/moke-game/game/pull/19)+ | Subscribe/CAS/CORS/TLS/Auth + binder fail-closed |
| `v1.0.5-0.20260811085843-60e104db6db7` (#221) | platform#24 base | game#19 base | Request-path auth fail-closed only |
| `v1.0.4` | older | older | Pre-hardening |

`create-game` skill templates follow game secure defaults (no `WithoutAuth` on public
services, AuthModule stub, gRPC+HTTP default, TCP opt-in, cancelable Watch).

## Tracking issues

- kit plan: https://github.com/GStones/moke-kit/issues/223
- platform plan: https://github.com/moke-game/platform/issues/23
- game plan: https://github.com/moke-game/game/issues/18
