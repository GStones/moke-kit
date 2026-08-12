# Compatibility

Use this matrix to track which `moke-kit` versions are validated with the
reference platform and game repositories.

| moke-kit version | platform | game | Notes |
| --- | --- | --- | --- |
| `v1.0.5-0.20260812031519-a995d2ada33f` ([#229](https://github.com/GStones/moke-kit/pull/229)) + [#230](https://github.com/GStones/moke-kit/pull/230) (+ smoke/DocumentBase/NATS on tip) | [platform#28](https://github.com/moke-game/platform/pull/28) / [#29](https://github.com/moke-game/platform/pull/29) | [game#24](https://github.com/moke-game/game/pull/24) / [#25](https://github.com/moke-game/game/pull/25) | create-game thin/smoke CI; AUTH_URL defaults commented; DocumentBase + NATS lifecycle tests |
| `v1.0.5-0.20260812022140-acb9f313d7fd` ([#228](https://github.com/GStones/moke-kit/pull/228)) | [platform#27](https://github.com/moke-game/platform/pull/27) | [game#22](https://github.com/moke-game/game/pull/22) | CAS/StopServing/MQ race tests + CI lint/vuln; platform jwt/v5 + CAS/chat lifecycle |
| `v1.0.5-0.20260811094419-bcdfe55515cd` (#224) + create-game [#226](https://github.com/GStones/moke-kit/pull/226) | [platform#25](https://github.com/moke-game/platform/pull/25) | [game#20](https://github.com/moke-game/game/pull/20) | Binder fail-closed + scaffold auth defaults |
| `v1.0.5-0.20260811094419-bcdfe55515cd` (#224) | [platform#24](https://github.com/moke-game/platform/pull/24) | [game#19](https://github.com/moke-game/game/pull/19) | Subscribe/CAS/CORS/TLS/Auth + binder fail-closed |
| `v1.0.5-0.20260811085843-60e104db6db7` (#221) | platform#24 base | game#19 base | Request-path auth fail-closed only |
| `v1.0.4` | older | older | Pre-hardening |

`create-game` skill templates follow game secure defaults (`AllModule` + stub auth in
`main`, no `WithoutAuth` on public services, gRPC+HTTP default, TCP opt-in, cancelable Watch).

`create-game` also ships `service-thin`, CI lint templates, k6 stub/platform auth modes, and
keeps auth out of transport modules (pair stub or platform middleware in `main` only).

## Tracking issues

- kit plan: https://github.com/GStones/moke-kit/issues/223
- platform plan: https://github.com/moke-game/platform/issues/23
- game plan: https://github.com/moke-game/game/issues/18
