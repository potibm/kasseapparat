# ADR 004: Token-guarded public endpoint

**Date:** 2026-09-27
**Status:** Accepted — contains a stopgap to remove in the next major version

## Context & Problem

`GET /api/v3/purchases/stats` was meant to be protected. It was added in `2dbcc08` inside the authenticated router, and `64cf6f8` ("fix: cors for purchase stats (hot-fix, rather)", PR #164) moved it onto the root router _above_ the CORS middleware to make a cross-origin call work. A CORS fix dropped authentication as a side effect, so for two years the endpoint has served per-product sales quantities to anyone who asks, with `Access-Control-Allow-Origin: *` on the response.

The endpoint is genuinely needed without a login: a statistics display may sit anywhere, such as a shop window screen on a machine nobody logs into.

The decision is therefore not "auth or no auth" but "what authorises a read of this one endpoint, given that a login is not available to the caller".

## Decision

**A shared bearer token, checked per request, on that route only.**

```
Authorization: Bearer <auth.public_endpoint_token>
```

`middleware.PublicEndpointAuth` validates the header and nothing else. It grants no identity, no session and no access to any other endpoint. Comparison is constant-time, and an empty configured token rejects every request rather than accepting an empty one.

**CORS is handled by a dedicated middleware instance, not a global one.** The route is registered as its own router group _before_ the app-wide `r.Use(CreateCorsMiddleware(...))`, which works because gin snapshots a group's handlers when the group is created. It also gets an explicit `OPTIONS` route: a preflight does not match the `GET` route, and without it gin's fallback chain would run the restrictive app-wide CORS middleware and reject the origin.

`/health` and `/ready` stay unauthenticated, since the container `HEALTHCHECK` probes the former and there is nowhere to put a credential. See [ADR 002](002-readiness-endpoint-scope.md).

## Considered and rejected: moving the route into the protected group

The obvious fix. It was rejected because the display has no way to log in, so the endpoint would be unusable — the feature is the reason the route was opened up in the first place.

## Considered and rejected: a separate signing or HMAC scheme

A per-request signature would stop a leaked token from being replayed and would survive token rotation. It is more machinery than this endpoint warrants, and the caller is a display, not a payment flow.

## Consequences

- **Positive:** the endpoint is no longer world-readable, and the change is small and local: one config value, one middleware, one route.
- **Positive:** existing deployments that never set the token keep working after the upgrade, which is what a two-year-old accidental hole requires.
- **Negative:** **the token is a shared secret with no expiry, no rotation and no per-caller identity.** Anyone holding it can read sales quantities. A browser that stores it exposes it to that origin.
- **Negative:** **the stopgap.** An unset token produces a random value at each start, logged as a warning, so the endpoint stays reachable on a fresh install. The cost is that the value changes on every restart, breaking any display or script until an operator sets one, and the secret is written to the logs — which, with the shipped OpenObserve setup, means a third party may hold it. **This is a stopgap, not a design.** It should be replaced in the next major version, where a breaking change is acceptable, by one of:
  - refusing to start when the token is unset, so the value must be chosen deliberately; or
  - a per-caller signing scheme, if the display needs rotation.
- **Negative:** CORS is wide open for this route, so the token is the only control. If the display's origin is known, pin it in `app.cors_allow_origins` and give this route the same treatment as the rest of the API.
