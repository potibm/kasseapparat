# ADR 002: Readiness Endpoint Scope

**Date:** 2026-09-26
**Status:** Accepted

## Context & Problem

Kasseapparat exposes two unauthenticated probes on the root of the application:

- `GET /health` — liveness. Always answers `200` while the process runs, and performs no dependency checks. This is what the container `HEALTHCHECK` uses.
- `GET /ready` — readiness. Answers `503` when the SQLite database is unreachable and `200` otherwise.

In `oidc` auth mode the identity provider is the gate to the whole application: without it nobody can log in. It is therefore reasonable to ask why `/ready` does not consider it.

This ADR records the scope of the readiness endpoint and, importantly, the reasoning behind the one dependency we deliberately left out.

## Decision

`/ready` checks the **database only**. It does not contact the OIDC identity provider.

### The scope rule

A dependency may gate readiness only if **restarting Kasseapparat would plausibly repair its failure**. If a restart cannot fix it, a failing probe is not actionable — it would only remove healthy capacity and, if wired to a container `HEALTHCHECK` or a Kubernetes `readinessProbe`, turn an external dependency's hiccup into a restart storm across every replica.

The database qualifies: it holds all state, the service genuinely cannot serve traffic without it, and a restart is a reasonable remedy.

### Considered and rejected: probing the OIDC provider

An earlier proposal was to add a cached, timeout-bounded probe of the provider's discovery document to `/ready`. It was dropped. The reasoning, so it is not re-proposed without context:

1. **A restart does not fix it.** If Dex or Keycloak is unavailable, restarting every Kasseapparat replica would not restore logins — it would only replace working processes with identical ones. Gating readiness on the IdP converts a login blip into a full application outage.

2. **Existing sessions are unaffected.** Sessions are verified _offline_. The OIDC middleware (`internal/app/middleware/auth.go`) decodes an encrypted `HttpOnly` session cookie with `gorilla/securecookie` and performs no token introspection and no network call. A user who is already logged in keeps working normally while the provider is down; only _new_ logins break. Reporting the service as not-ready would misdescribe the experience of every currently-authenticated user.

3. **Operational complexity outweighs the benefit.** A correct implementation needs a cached probe with a TTL, a short timeout, a decision on how to surface the result in the response contract, and monitoring guidance. That surface area was judged not worth paying for the signal it yields — the signal is available more cheaply elsewhere.

## Consequences

- **Positive:** Readiness is cheap and deterministic. The endpoint performs exactly one local operation, with no external network dependency that could make it slow or flaky.
- **Positive:** An IdP outage cannot make Kasseapparat instances appear unhealthy, so no restart storm is possible and logged-in staff keep taking payments.
- **Positive:** The rule is general. New dependencies can be evaluated against a single question rather than re-litigating the trade-off each time.
- **Negative:** **An identity provider outage is not observable through `/ready`.** This is the accepted cost. Operators must monitor the provider directly, or alert on login-failure rates and on failed `/api/v3/auth/login` attempts, to detect an IdP outage. This trade-off is the explicit reason for writing this ADR down.
- **Neutral:** Startup is the one point where the provider _is_ contacted. `NewOIDCAuthHandler` performs a live OIDC discovery request, so an unreachable issuer prevents the process from booting. That request is bounded by a 10-second timeout (`oidcDiscoveryTimeout` in `internal/app/handler/http/oidc_auth.go`) rather than `http.DefaultClient`, which has no timeout and would block startup indefinitely. The timeout is generous because failing to start is a harder outcome than a slow boot, and a slow-but-healthy provider — a cold Keycloak start, for example — must not be mistaken for an unreachable one.
