import Decimal from "decimal.js";

/** Bank notes the customer is likely to hand over. */
export const BANKNOTES = [5, 10, 20, 50, 100, 200];

/**
 * Amounts a customer is likely to pay with: the next round sums above the total,
 * one per bank note, without duplicates (e.g. 37 → 40, 50, 100, 200).
 */
export const suggestReceivedAmounts = (total: Decimal, max = 4): Decimal[] => {
  const suggestions: Decimal[] = [];

  for (const note of BANKNOTES) {
    const amount = total.div(note).ceil().mul(note);

    if (
      amount.greaterThan(total) &&
      !suggestions.some((suggestion) => suggestion.equals(amount))
    ) {
      suggestions.push(amount);
    }

    if (suggestions.length >= max) break;
  }

  return suggestions;
};

/**
 * Reads an amount typed by the cashier. Accepts a comma or a dot as decimal
 * separator; returns null for anything that is not a non-negative number.
 */
export const parseAmount = (input: string): Decimal | null => {
  const normalized = input.trim().replace(",", ".");
  if (normalized === "" || !/^\d+(\.\d+)?$/.test(normalized)) return null;

  return new Decimal(normalized);
};
