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
- **Negative:** **A transaction callback that used the outer handle instead of the transaction-scoped one would deadlock** rather than merely be slow, because the transaction already holds the only connection. The failure is silent and permanent: the waiting query blocks in `database/sql.(*DB).conn` on a `context.Background()` that nothing cancels, so the process must be killed rather than timing out. **Any transaction callback must use its `tx` parameter.** `WithTransaction` (`internal/app/repository/sqlite/repository.go`) enforces this by handing the callback a repository cloned onto `tx`, but it is not the only way to open a transaction — `seedPurchases` (`internal/app/utils/seed.go`) opens one directly and got this wrong, deadlocking `mise run e2e:setup` and with it the whole Playwright job. GORM's `Session` preserves the transaction's connection pool, so the audit callbacks in `internal/app/store/gorm/audit.go` are also safe. Pinned by `TestSeedDatabaseOnBoundedPoolDoesNotDeadlock`, which runs the seed on the bounded pool under a watchdog, and by `TestBoundedPoolTransactionCallbackUsingOuterHandleErrors`, which pins that a reach-back surfaces as a context error rather than a silent hang.
- **Neutral:** `ConnectToLocalDatabase` (`file::memory:?cache=shared`, used by the e2e suite) is not bounded. Shared-cache mode raises `SQLITE_LOCKED` at table level, which the driver's `unlock_notify` retry path hard-returns on without retrying, so a pool above one is _more_ dangerous on that path, not less. It also benefits from multiple connections for test concurrency. The two paths are genuinely different and one shared helper would have been wrong for both. The cost of that split is that a test on the in-memory path cannot see this decision's failure mode at all — which is why the seed regression above is written against the bounded `ConnectToDatabase` path.

## Follow-up: enable WAL

If write latency under load becomes a problem, the fix is to enable WAL, not to grow the pool on its own:

```go
sqlDB.Exec("PRAGMA journal_mode = WAL;")
```

Under WAL, readers never block the writer, and a pool above one then becomes safe and worthwhile. That is the whole appeal, and it is worth being precise about what WAL does **not** do on its own:

- **WAL alone buys no throughput while the pool is one.** A single connection is serialized by `database/sql` regardless of journal mode. WAL is only worth enabling together with a larger pool, so the two are one decision, not two toggles.
- **A larger pool under WAL still needs a busy handler.** WAL permits one writer at a time, so concurrent writes can still collide and return `SQLITE_BUSY`. Today the DSN is a bare path with no `_pragma`, so there is no handler at all and a losing writer fails immediately. `busy_timeout` must land in the same change; the driver supports it (`modernc.org/sqlite` honours `_busy_timeout` and `_pragma` in the DSN).
- **`journal_mode` is persisted in the database header.** Setting it on an existing file is a one-way door: downgrading the application does not revert it. WAL also depends on shared-memory primitives that network filesystems do not provide reliably, so a `data/` directory on a NAS would need to be ruled out first.

Two things must be dealt with before this is safe to ship, neither of which is a code detail:

- **The documented backup procedure becomes unsafe.** `docs/admin.md` tells operators to back up with `tar cfvz $FILENAME $BASEDIR/data/`, a hot copy taken while the container is still running, and there is no `database backup` command to replace it. In `delete` journal mode the worst case is losing the in-flight transaction. Under WAL, a missing or mismatched `-wal` sidecar silently loses **every transaction since the last checkpoint**, with no error to notice. Every deployed instance has its own hand-rolled `update.sh` built from that snippet, so the procedure and the docs have to change together. **This is the real prerequisite: a `database backup` command on SQLite's online backup API.** It needs no WAL, it fixes the hot-copy weakness that exists today, and it is what makes WAL adoptable later without operator risk.
- **The E2E fixture reset needs sidecar handling.** `resetDatabase` in `frontend/e2e/helpers/db.ts` copies the `.db` file underneath a running backend. Under WAL the meaningful state also lives in `-wal` and `-shm`, so copying the main file alone would silently serve a stale fixture.

There is also a checkpoint policy to settle for the small disks these tills run on, since an uncheckpointed `-wal` grows without bound.

So: this wants its own decision record and its own load test rather than being folded in here. It is the natural next step, and it is deliberately not this one.
