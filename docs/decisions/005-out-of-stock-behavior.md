# ADR 005: Out-of-stock behaviour

**Date:** 2026-09-27
**Status:** Accepted

## Context & Problem

Products carry a `total_stock`, and the quantity sold is derived by summing the purchase items of confirmed purchases — there is no counter column to keep in step. Before this change nothing consulted that number during a sale, so a limited product could be sold past its stock. That is the right behaviour for drinks and open-ended tickets, and the wrong one for a fixed run of T-shirts or mugs.

The decision is not whether to restrict stock, but what should happen when it runs out, because that differs per event: some merchandise should simply stop selling, some should disappear from the POS, and some should keep selling with a nudge.

## Decision

**A single global setting, `app.out_of_stock_behavior`, with four values: `ignore`, `fail`, `auto_sold_out` and `auto_hide`.**

`ignore` is the default, and an unset or unrecognised value also behaves as `ignore`. That is deliberate: the value is validated at startup so a typo cannot pass unnoticed, but the sales logic must not fail _closed_ on a missing setting, because that would silently start rejecting sales on an existing deployment.

`total_stock: 0` means unlimited and is excluded from every mode.

## Where the check runs, and why

The check runs **inside the purchase transaction**, on the transaction-scoped repository, not in the price validation that runs before it. Two reasons:

- A check outside the transaction races the insert that follows it. [ADR 003](003-sqlite-connection-pool.md) bounds the pool to a single connection, so a read through the outer repository from inside a transaction would not merely be stale, it would wait for a connection that is never released.
- Availability is a query, not a counter, so there is nothing to make atomic beyond reading and writing inside one transaction.

The quantity requested is **aggregated per product across the whole cart** before checking. Checking each cart line separately lets two lines of the same product both pass against the same availability, because neither sees the other's quantity.

The depletion comparison is `>=`, not `==`, so a product that is already oversold still converges instead of never being flagged.

## Considered and rejected: a per-product override

A global setting is what the story asked for and is the smallest thing that covers the four cases. A per-product override would be a natural follow-up, and the aggregation point is already per product, so it would not be a rewrite.

## Considered and rejected: a stock counter column

Maintaining a `units_sold` column would make the check a single-row read. It would also introduce a second source of truth that can drift from the purchase rows — through a failed transaction, a deleted purchase, or a manual database edit. Deriving the number keeps it correct by construction, and the cost is one aggregate query per product in the cart.

## Consequences

- **Positive:** overselling is prevented in the three enforcing modes, and the POS cannot offer a quantity the server will reject, because both sides use the same definition.
- **Positive:** the default is the previous behaviour, so no existing deployment changes what it can sell.
- **Positive:** a pending SumUp purchase reserves its units, and releasing them on failure or cancellation keeps the number honest without a reservation table.
- **Negative:** **a refund re-enables a product that an operator hid by hand.** The product row does not record whether a flag was set by the behaviour or by a person, so both are cleared when the stock comes back. Distinguishing them would need columns tracking auto-managed flags; that was judged not worth the schema for now. A product that is still sold out keeps its flag, so a refund of an unrelated product cannot put it back on sale.
- **Negative:** changing the setting does not re-evaluate products that are already flagged. Existing flags are the operator's to clear.
- **Negative:** the check is per product, not across products, so a cart can still fail as a whole after each line passed — the individual line is what is validated.
- **Neutral:** the setting lives in the config file, not the admin interface. Changing it needs a config edit and a restart.
