# ADR 006: Product sections and product groups

**Date:** 2026-11-20
**Status:** Proposed

## Context & Problem

The POS lays products out as one sorted list (`pos`). The only structure is `wrap_after`,
a boolean that ends the current row after a product. That was enough for a handful of
tickets, but not for a real event:

- There is no way to name a block of products ("Tickets", "T-Shirts", "Merch"), so the
  cashier scans the whole grid for every sale.
- Variants of one article are separate products with nothing tying them together. At
  Evoke the T-shirts are eleven products (`Male S` … `Male 4XL`,
  `Female S` … `Female XXL`). As cards they take three rows of the iPad screen; one
  missing size (3XL) went unnoticed because nothing shows the sizes side by side.

The POS redesign shows products in **sections** with a heading, and variants as a **size
matrix**: one row per group, one button per variant. This ADR decides how that structure
is stored. Stock, sales and interest must keep working per sellable item, including the
out-of-stock behaviour of [ADR 005](005-out-of-stock-behavior.md).

## Decision

**Two small tables, `sections` and `product_groups`, referenced from `products`.
`wrap_after` is removed.**

```
sections        id, name, pos
product_groups  id, section_id, name, pos
products        …, section_id (required), product_group_id (nullable), variant (nullable)
```

A **product stays the sellable unit.** Price, VAT, stock, `sold_out`, `hidden`, purchase
items and interest remain on the product, so a size is sold, counted and reconciled
exactly like any other product today. Sections and groups only describe the layout.

## Rules

**Every product belongs to exactly one section.** `section_id` is required. A fresh setup
(seeder, `config create`) creates one section, and the admin preselects the first section
for a new product.

**A section's name may be empty.** An unnamed section renders as a plain row break without
a heading — the replacement for `wrap_after`.

**A group lives entirely inside one section.** A product's section is derived from its
group: when `product_group_id` is set, the backend sets `section_id` to the group's section
and ignores a different value. Moving a group moves its products with it, and a product can
never sit in a group that belongs to another section.

**Order.** Sections are ordered by `pos`. Inside a section, standalone products and groups
share one `pos` space, so a group can sit between two products. Inside a group, products
are ordered by their own `pos`. Ties are broken by `id`, as for products today.

**Display of a group.** One row per group: the group name on the left, one button per
product, labelled with `variant` (falling back to the product name). If all products of
the group have the same price it is shown once in the row header, otherwise on each
button. A group whose products are all hidden is not shown.

**Deleting.** A section or a group can only be deleted when it is empty. The API answers
`409 Conflict` otherwise. Silently moving or ungrouping products would change the POS
layout as a side effect of a delete.

**Out-of-stock behaviour is unchanged.** It still works per product. Under `auto_hide` a
sold-out size disappears from the backend's list; the POS keeps its slot as a placeholder
(see `keepVanishedProducts`), which in the matrix means a disabled size button.

## API

`sections` and `productGroups` become regular resources with the same list / get /
create / update / delete endpoints and conventions as the other admin resources. As with
products, every user may create and change them; deleting requires the admin role.
Products gain `sectionId`, `productGroupId` and `variant`.

A product created without `sectionId` goes into the default section. On update, an
omitted or `0` `sectionId` means "unchanged"; moving a product into the default section
uses that section's real id.

`productGroupId` on a product update: omitted or `null` means "unchanged", `0` takes
the product out of its group (it stays in the group's section), any other id moves it
into that group. A group created without `sectionId` goes into the default section; on
a group update, `0` keeps the current section, and moving the group moves its products
with it.

The POS loads sections, groups and products and composes the layout itself. All three
lists are small, and they are only reloaded together with the products.

## No data migration

An installation is set up fresh for each event; there are no existing layouts to carry
over. `wrap_after` is dropped from the model, the admin and the POS without converting
old values.

The one exception is that every product needs a section. On startup, products without
one (a database created before sections existed) are put into the default section: the
first by `pos`, or a new unnamed one. This is idempotent and does nothing on a fresh
installation. The seeder and the e2e fixtures are updated in the same pull requests as the
schema, so the daily seeded test instance always shows sections and groups.

## Considered and rejected: fields on the product only

`section_title`, `group` and `variant` as plain fields on the product, with sections and
groups implied by consecutive equal values, would need no new tables. But the admin UI
treats sections and groups as things you rename and drag, and the implicit version has
edge cases for each of those: renaming a group is a bulk update, a product sorted between
two "Male" products splits the group into two rows, an empty section cannot exist, and
group-level properties (display name, later perhaps an image) have no place to live.

## Considered and rejected: deriving groups from product names

"T-Shirt Male L" → group "T-Shirt Male", variant "L" needs no schema change at all, but it
guesses: "Hoodie Zip XL" or "Mug 0,3 l" would be split wrongly, and a typo silently creates
a new group. Explicit fields are unambiguous.

## Considered and rejected: one endpoint for the whole POS catalog

A nested `/pos/catalog` (sections → groups → products) would save the POS two requests.
It would also be a second representation of products next to `/products`, which the
admin, the stock reconciliation and the e2e helpers already use. Three small requests are
cheap; a combined endpoint can be added later without changing the model.

## Consequences

- **Positive:** the POS can show named blocks and size matrices; the Evoke shirts take two
  rows instead of three and missing sizes are obvious.
- **Positive:** sales, stock, exports, statistics and ADR 005 keep working per product,
  unchanged.
- **Positive:** renaming or reordering a section or group is a single-row update.
- **Neutral:** setting up an event means creating a few sections and groups besides the
  products. A "create variants" dialog and a catalog import/export are planned to keep that
  quick.
- **Negative:** two more resources in the API and the admin, and the backend has to keep
  `section_id` consistent with the group (derived on write, see above).
- **Negative:** no compatibility with data from earlier versions; a v3 database is not
  upgraded in place. This is acceptable because installations are set up per event.
