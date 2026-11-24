/** At or below this many units left, the stock counts as low (as in the POS). */
export const LOW_STOCK_THRESHOLD = 3;

export type StockTone = "unlimited" | "normal" | "low" | "out";

export interface StockSummary {
  /** Units left, or null when the stock is not limited. */
  left: number | null;
  sold: number;
  total: number | null;
  tone: StockTone;
}

/** Stock figures of a product for the admin list. totalStock 0 means unlimited. */
export const summarizeStock = (product: {
  totalStock?: number;
  unitsSold?: number;
}): StockSummary => {
  const sold = product.unitsSold ?? 0;
  const total = product.totalStock ?? 0;

  if (total <= 0) {
    return { left: null, sold, total: null, tone: "unlimited" };
  }

  const left = Math.max(0, total - sold);
  let tone: StockTone = "normal";
  if (left === 0) {
    tone = "out";
  } else if (left <= LOW_STOCK_THRESHOLD) {
    tone = "low";
  }

  return { left, sold, total, tone };
};
