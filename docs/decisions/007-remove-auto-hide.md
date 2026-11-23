# ADR 007: Remove the `auto_hide` out-of-stock behaviour

**Date:** 2026-11-23
**Status:** Accepted – supersedes the `auto_hide` part of [ADR 005](005-out-of-stock-behavior.md)

## Context & Problem

ADR 005 introduced four values for `app.out_of_stock_behavior`. `auto_hide` hid a
limited product from the POS once its stock was gone. With the v4 redesign of the POS
that no longer pays off:

- A sold-out card is useful. It says "Sold out · N interested" and lets the cashier
  register interest, which is what the product interest list and its per-product summary
  are built on. A hidden product collects nothing.
- Sold-out products take little room: a card in the grid, or a struck-through size button
  in the size matrix of a product group (ADR 006).
- `auto_hide` was the most expensive mode. The grid shifted when a product disappeared
  mid-sale, so the POS kept placeholder slots for vanished products until a reload
  (`keepVanishedProducts`), and the mode owned the `hidden` flag, which therefore was not
  a free manual choice for limited products.

## Decision

**`auto_hide` is removed. `ignore`, `fail` and `auto_sold_out` remain.**

- An existing config with `auto_hide` keeps starting: it is rewritten to `auto_sold_out`
  at startup and a warning is logged (`Config.ApplyDeprecations`). The value still passes
  validation for that reason; nothing else accepts it.
- Products that `auto_hide` had hidden are recognised by their flags – managed by the
  system, hidden, not sold out. `auto_sold_out` never produces that combination, because
  it always sets `sold_out` on the products it manages. At the next reconciliation (every
  start) their hidden flag is released and they are marked sold out if their stock is
  still gone. A product an operator hid by hand is either not managed or also sold out,
  so it keeps its flag.
- `hidden` is a purely manual flag again, for every product in every mode.
- The POS no longer keeps placeholder slots; the product list is shown as the backend
  returns it.

## Consequences

- **Positive:** Sold-out products keep collecting interest. Less code in the backend
  (one mode, one flag the system owns) and in the POS (no placeholder bookkeeping).
- **Positive:** Old configs keep working; the log tells the operator to update them.
- **Negative:** An event that wants sold-out merchandise to disappear has to hide it by
  hand in the admin.

## Rejected alternatives

- **Keep `auto_hide`.** The cost above buys a layout that the redesign no longer needs.
- **Reject `auto_hide` at startup.** Safer to spot, but a forgotten config line would
  stop the POS on the day of an event; a warning is enough.
