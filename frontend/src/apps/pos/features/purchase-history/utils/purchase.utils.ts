import { Purchase } from "../../../api/schemas";

/**
 * Cashiers may only refund a purchase within 15 minutes of its creation; admins
 * at any time. The backend enforces this (see docs/decisions/001-purchase-refund-
 * time-limit.md); the POS mirrors it so nobody taps a button that can only fail.
 */
export const REFUND_WINDOW_MS = 15 * 60 * 1000;

/** "2× Mate, Beer, T-Shirt" */
export const summarizeItems = (purchase: Purchase): string =>
  (purchase.purchaseItems ?? [])
    .map((item) =>
      item.quantity > 1
        ? `${item.quantity}× ${item.product.name}`
        : item.product.name,
    )
    .join(", ");
