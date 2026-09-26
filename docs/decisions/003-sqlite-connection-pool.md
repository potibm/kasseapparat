# ADR 003: SQLite Connection Pool Size

**Date:** 2026-09-26
**Status:** Accepted

## Context & Problem

Kasseapparat opens its SQLite database through GORM, which hands back a `*gorm.DB` backed by a `database/sql` pool. Nothing configured that pool, so it ran on `database/sql` defaults, and the relevant default is permissive:

```go
maxOpen int // <= 0 means unlimited
```

`MaxOpenConns` therefore had to be explicitly set to have any effect at all. The pool size matters here for two reasons that pull in opposite directions:

- **Unbounded connections are a contention risk.** Every connection is a separate handle on one file. In the default `delete` journal mode a writer takes an exclusive lock and blocks readers, so extra connections do not add throughput — they convert "wait your turn" into competing for the same lock. Nothing configures a busy handler: the DSN is a bare path (`utils.ConnectToDatabase` passes no `_pragma`), so SQLite's default of _no handler_ applies and a writer that loses the race is handed `SQLITE_BUSY` immediately instead of waiting its turn. That surfaces to a customer mid-purchase.
- **A single connection has a failure mode too.** With `MaxOpenConns(1)`, a slow query occupies the only connection and everything behind it waits. That is the case the readiness probe cares about: a probe queued behind a long write cannot report.

## Decision

**`ConnectToDatabase` sets `MaxOpenConns(1)`.** The in-memory test path is deliberately left alone — see Consequences.

`database/sql` clamps `MaxIdleConns` to `MaxOpenConns`, so the single connection is also the single idle connection; setting `MaxIdleConns` explicitly would be a no-op.

## Considered and rejected: a pool above one

A larger pool is the right answer _in general_, and it is the wrong answer today for one specific reason: **multi-connection concurrency is only safe under WAL.** In `delete` journal mode, readers still block the writer, so the pool cannot buy read/write overlap — it can only move the queueing from `database/sql` into SQLite's lock manager, where the failure mode is a customer-visible error instead of a request that simply waits.

The upside of a larger pool is also already covered. The scenario motivating it — the readiness probe stalling behind a slow query — is handled by the 2-second bound on `GET /ready` (see [ADR 002](002-readiness-endpoint-scope.md)), which reports the condition rather than hanging on it. So raising the pool size would add a new failure mode without removing an existing one.

`SetMaxIdleConns(MaxOpenConns)` was also dropped: `database/sql` already applies that cap internally, so the second call cannot do anything.

## Consequences

- **Positive:** `SQLITE_BUSY` and `SQLITE_LOCKED` cannot arise from intra-process contention, because there is no second connection to contend with. This is a real improvement over the previous unlimited pool.
- **Positive:** One file handle and one page cache instead of one per concurrent request, which matters on the small hardware a till typically runs on.
- **Positive:** The bound is asserted by a test, so a future change cannot silently reopen the pool.
- **Negative:** **All queries serialize.** Under load a slow report or export delays purchases. This is the deliberate trade: queued latency instead of lock errors, which is the right choice for a point of sale.
- **Negative:** **A transaction callback that used the outer repository instead of the transaction-scoped one would deadlock** rather than merely be slow, because the transaction already holds the only connection. The current call sites in `internal/app/service/purchase/purchase.go` all use `txRepo` inside `WithTransaction`, and `TestConnectToDatabaseTransactionDoesNotDeadlock` pins that. GORM's `Session` preserves the transaction's connection pool, so the audit callbacks in `internal/app/store/gorm/audit.go` are also safe. **A new transaction callback must not reach for the outer repository.**
- **Neutral:** `ConnectToLocalDatabase` (`file::memory:?cache=shared`, used by the e2e suite) is not bounded. Shared-cache mode raises `SQLITE_LOCKED` at table level, which the driver's `unlock_notify` retry path hard-returns on without retrying, so a pool above one is _more_ dangerous on that path, not less. It also benefits from multiple connections for test concurrency. The two paths are genuinely different and one shared helper would have been wrong for both.

## Follow-up: enable WAL

If write latency under load becomes a problem, the fix is to enable WAL, not to grow the pool:

```go
sqlDB.Exec("PRAGMA journal_mode = WAL;")
```

Under WAL, readers never block the writer, and a pool above one then becomes safe and worthwhile. This is a behaviour change rather than a tuning knob — it alters the on-disk format, is not reversible on an existing file without care, and interacts with network filesystems — so it wants its own load test and its own decision record rather than being folded in here. It is the natural next step, and it is deliberately not this one.
