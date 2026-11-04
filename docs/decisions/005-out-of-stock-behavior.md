# ADR 005: Out-of-stock behaviour

**Date:** 2026-09-27 (revised 2026-10-15)
**Status:** Accepted

## Context & Problem

Products carry a `total_stock`, and the quantity sold is derived by summing the purchase
items of confirmed purchases — there is no counter column to keep in step. Before this
change nothing consulted that number during a sale, so a limited product could be sold
past its stock. That is the right behaviour for drinks and open-ended tickets, and the
wrong one for a fixed run of T-shirts or mugs.

The decision is not whether to restrict stock, but what should happen when it runs out,
because that differs per event: some merchandise should simply stop selling, some should
disappear from the POS, and some should keep selling with a nudge.

## Decision

**A single global setting, `app.out_of_stock_behavior`, with four values: `ignore`,
`fail`, `auto_sold_out` and `auto_hide`.**

`ignore` is the default, and an unset or unrecognised value also behaves as `ignore`. That
is deliberate: the value is validated at startup so a typo cannot pass unnoticed, but the
sales logic must not fail _closed_ on a missing setting, because that would silently start
rejecting sales on an existing deployment.

`total_stock: 0` means unlimited and is excluded from every mode.

## Where the check runs

The check runs **inside the purchase transaction**, on the transaction-scoped repository,
not in the price validation that runs before it. Two reasons:

- A check outside the transaction races the insert that follows it.
- [ADR 003](003-sqlite-connection-pool.md) bounds the pool to a single connection, so a
  read through the outer repository from inside a transaction would not merely be stale,
  it would wait for a connection that is never released.

The quantity requested is **aggregated per product across the whole cart** before checking.
Checking each cart line separately lets two lines of the same product both pass against the
same availability, because neither sees the other's quantity.

The depletion comparison is `>=`, not `==`, so a product that is already oversold still
converges instead of never being flagged.

Only **confirmed** purchases count. A pending SumUp payment holds no units: it becomes sold
when it is confirmed, and a refund, a failure or a deletion releases it again. An abandoned
pending purchase would otherwise hold the last units forever, because nothing expires it.

## The product flags are derived, not maintained incrementally

`auto_sold_out` and `auto_hide` work through the `sold_out` and `hidden` columns the POS
already reads: a hidden product is filtered out of the POS list, a sold-out one renders the
existing "Sold Out" card with its interest dialog. Nothing consumes a flag the API did not
already expose, so a second client needs nothing new to display it.

`PurchaseService.ReconcileStockFlags` derives both columns from the stock in one pass over
all products, and runs after a purchase is stored, after every purchase status change, after
a product is edited, after a purchase is deleted, and once at startup:

```sql
-- one query, not one per product
SELECT product_id, SUM(quantity) FROM purchase_items JOIN purchases ... WHERE status = 'confirmed' GROUP BY product_id
```

Derived state beats incremental state here. The pass is idempotent, and it repairs what an
update per purchase could never reach: a `total_stock` an operator raised, a deleted
purchase, or a flag left behind by a previous mode. It is also the whole reason there is no
second source of truth that can drift.

**A flag an operator set by hand is not the behaviour's.** `stock_flags_managed` on the
product records whether the behaviour set the flag; a product carrying a flag the behaviour
does not own is skipped entirely, so a refund cannot put a product an operator hid back on
sale. Existing rows count as unmanaged, so manual flags survive the first pass.

A product that already carries the state the behaviour wants is not written at all, which
keeps a no-op pass — and every pass under `ignore` and `fail`, which never manage a flag —
free of writes.

## Considered and rejected: a per-product override

A global setting is what the story asked for and is the smallest thing that covers the four
cases. A per-product override would be a natural follow-up, and the aggregation point is
already per product, so it would not be a rewrite.

## Considered and rejected: a stock counter column

Maintaining a `units_sold` column would make the check a single-row read. It would also
introduce a second source of truth that can drift from the purchase rows — through a failed
transaction, a deleted purchase, or a manual database edit. Deriving the number keeps it
correct by construction, and the cost is one aggregate query per product in the cart.

## Considered and rejected: reserving units for a pending payment

Counting pending purchases as sold would stop a second client from taking the last units
while a SumUp payment is in flight. Nothing expires a pending purchase, so an abandoned one
would hold the units for the rest of the event — with one POS the risk it prevents is smaller
than the risk it introduces. The payment is the operator's to complete, so the units become
sold when it is confirmed.

## Consequences

- **Positive:** overselling is prevented in the three enforcing modes, and the POS cannot
  offer a quantity the server will reject, because both sides use the same definition.
- **Positive:** the default is the previous behaviour, so no existing deployment changes what
  it can sell.
- **Positive:** switching the mode no longer strands products flagged by the old one, because
  startup re-evaluates them.
- **Positive:** a product that runs out and later comes back — through a refund, a failure, or
  a raised stock — returns to sale by itself.
- **Neutral:** a manual flag wins over the behaviour. A product an operator sold out or hid
  by hand stays that way until the operator clears it, even under `ignore`.
- **Negative:** the pass runs on the request that caused it. With the stock query being a
  single aggregate, that is two queries per purchase; moving it off the request would only
  move the latency, because the pool is a single connection ([ADR 003](003-sqlite-connection-pool.md)).
- **Neutral:** the setting lives in the config file, not the admin interface. Changing it
  needs a config edit and a restart — which is also when existing flags are re-evaluated.
