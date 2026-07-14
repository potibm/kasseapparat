# ADR 001: Purchase Refund Time Limit

**Date:** 2026-07-09
**Status:** Accepted

## Context & Problem

To handle erroneous bookings efficiently without keeping transactions open indefinitely, we need a controlled mechanism for refunding purchases. Leaving refunds open without a time limit increases the risk of abuse or complex state changes long after a transaction has occurred.

## Decision

We enforce a time limit on refunds with the following rules:

- **Standard Users:** Purchases can only be refunded within a **15-minute window** after creation to correct immediate mistakes.
- **Administrators:** Admins are exempt from this time limit and are permitted to refund _any_ purchase at any given time.

## Consequences

- **Positive:** Reduces the risk of late, unauthorized refunds while providing users enough time to revert accidental bookings.
- **Positive:** Administrators retain full operational flexibility to resolve edge cases manually.
- **Negative:** The backend logic must strictly check both the age of the purchase (`time.Since(purchase.CreatedAt)`) and the authorization role of the executing user before processing a refund.
