# Bugbot review rules (moke-kit)

## Auth / binder
- Flag public gRPC/HTTP services that embed `utility.WithoutAuth` or skip AuthMiddleware.
- Flag handlers that trust client-supplied uid (e.g. request field) when `UIDContextKey` should be authoritative.
- Flag duplicate `sfx.AuthMiddleware` providers in the same fx graph (stub AuthModule + platform auth).
- Flag production binder paths that fail open on missing auth/CORS when fail-closed is intended.

## create-game templates
- Scaffold `AllModule` must not embed an auth provider; auth is paired in `main`.
- k6 / CLI smoke against authenticated APIs must send `authorization` metadata.
- TCP/zinx paths are unauthenticated — flag docs or defaults that expose them publicly.

## MQ
- Flag discarded `miface.Subscription` when the code unsubscribes before parent context cancel (must retain + `Unsubscribe`).
