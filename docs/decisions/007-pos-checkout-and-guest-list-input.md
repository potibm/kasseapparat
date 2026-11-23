# ADR 007: Change calculator, guest search and ticket scanning on the POS

**Date:** 2026-11-23
**Status:** Accepted

## Context & Problem

The v4 redesign made the POS a touch-first screen for an iPad in landscape. Three
decisions in that work change how sales and check-ins behave, not only how they look:

1. Cash sales need the change. Card sales do not.
2. The guest search is typed on an on-screen A–Z keyboard without umlauts, but names on the
   lists contain them ("Müller", "Straße").
3. Tickets are read by barcode scanners that the iPad sees as keyboards. We do not know
   whether a given scanner ends a code with Enter.

## Decision

### Change calculator per payment method

A payment method can show a dialog that asks for the amount received and shows the change
before the purchase is booked. It is configured per method with `change_calculator` in
`payment_methods`; unset means on for `CASH` and off for everything else. "Exact" is
preselected, so an exact payment costs one extra tap. A cart that costs nothing is booked
without the dialog. The amount received is not stored.

### Guest search ignores case and diacritics

The POS search (`GET /products/:id/guests?q=`) folds names and query to lower case without
diacritics (`ä → a`, `ß → ss`, `ø → o`, …) and matches the name as a substring; codes still
match as a case-insensitive prefix. The matching runs in Go after loading the product's
unattended guests: those lists are short, and SQLite's `LIKE` only ignores the case of
ASCII letters. The admin guest search is unchanged.

### Scanned tickets are checked in automatically

The scanner types into the search field of the guest list window, which has the focus for
lists with codes. When the results contain **exactly one** ticket whose code equals the
search term, it is added to the cart at once, the field is cleared and the window stays
open for the next scan. This depends only on the results, not on an Enter. Not automatic:
tickets that allow additional guests (the cashier picks +1/+2), tickets already in the
cart, and tickets without stock left. The window shows whether the field has the focus
("Ready to scan"), because that is when a scan arrives.

## Consequences

- **Positive:** Cash needs no mental arithmetic; card payments stay one tap. Cashiers
  find "Müller" by typing MULLER. Scanning works the same with every keyboard scanner.
- **Negative:** Typing a complete code by hand also checks the ticket in – the same intent,
  but no confirmation step. A focus lost to another button has to be restored by a tap on
  the scanner status line.

## Rejected alternatives

- **A global change-calculator switch.** Card and voucher methods would show a dialog
  that only costs time.
- **Fold in SQL with a custom SQLite function.** Needs a direct dependency on the
  driver's function API for a list that fits in memory anyway.
- **Rely on Enter to detect a scan.** Not every scanner sends it, and an Enter alone does
  not say which ticket was meant.
