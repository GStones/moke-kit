---
name: fx-module
description: Add or modify uber/fx modules in moke-kit (pkg/module, sfx/mfx/ofx providers, AppModule wiring). Use when creating a new fx module, provider, settings injection, or wiring a service into fxmain.
paths:
  - "**/pkg/**/*.go"
  - "**/fxmain/**/*.go"
  - "**/*fx*/**/*.go"
---

# uber/fx module patterns

moke-kit composes infrastructure with [uber/fx](https://github.com/uber-go/fx). Follow existing package layout instead of inventing new DI styles.

## Layout

| Layer | Location | Role |
| --- | --- | --- |
| App entry | `fxmain/` | `AppModule` aggregates required infra modules |
| Public module | `<area>/pkg/module` | `fx.Module("name", ...)` exported for apps |
| Providers | `<area>/pkg/<x>fx` (e.g. `sfx`, `mfx`, `ofx`) | `fx.Provide(...)`, settings, constructors |
| Internals | `<area>/internal` | implementation not imported by apps |

Examples:

- `server/pkg/module` → `sfx.*Module`
- `mq/pkg/module` → `mfx.MqModule`, `mfx.SettingModule`
- `orm/pkg/module` → `ofx.SettingsModule`, Redis/Mongo/GORM providers
- `fxmain/pkg/module.AppModule` includes server, orm, logging, mq

## Adding a provider

1. Define params/results with `fx.In` / `fx.Out` structs (see `server/pkg/sfx`, `orm/pkg/ofx`).
2. Export a `fx.Provide` variable (`XxxModule`).
3. Compose it into the area's `pkg/module.Module` via `fx.Module("area", ...)`.
4. Only add to `AppModule` if every app should get it by default; otherwise let consumers import the area module explicitly.

## Settings

- Use `envconfig` (or the local settings pattern already in that package).
- Document new env vars in the area `README.md` table (ENV / Description / Default).
- Safe defaults for local; fail closed for prod security-sensitive options.

## Checklist

- [ ] No circular imports between `module` and `internal`.
- [ ] Constructor errors propagate (return `error`); do not `log.Fatal` inside providers.
- [ ] Lifecycle (`OnStart` / `OnStop`) registers cleanup for servers, clients, and subscriptions.
- [ ] Unit or integration smoke covers the happy path when behavior is non-trivial.
- [ ] Public API stays in `pkg/`; keep implementation under `internal/`.

## Anti-patterns

- Global mutable singletons instead of fx-provided deps.
- Importing `internal/` from another area.
- Duplicating settings structs across packages instead of providing them once.
- Starting goroutines in `fx.Provide` without tying them to fx lifecycle.
