/**
 * Colours shared by the admin: the dashboard charts and the purchase list use the
 * same colour for a payment method, so they can be matched at a glance.
 */
export const PAYMENT_METHOD_COLORS: Record<string, string> = {
  CASH: "#4caf50",
  CC: "#2196f3",
  SUMUP: "#ff9800",
  VOUCHER: "#9c27b0",
};

const FALLBACK_COLORS = ["#e91e63", "#00bcd4", "#ff5722", "#607d8b", "#795548"];

/** The colour of a payment method; unknown methods get one by position. */
export const paymentMethodColor = (code: string, index = 0): string =>
  PAYMENT_METHOD_COLORS[code] ??
  FALLBACK_COLORS[index % FALLBACK_COLORS.length];

/** Colours for series of products, in order. */
export const SERIES_COLORS = [
  "#1f77b4",
  "#ff7f0e",
  "#2ca02c",
  "#d62728",
  "#9467bd",
  "#8c564b",
  "#e377c2",
  "#7f7f7f",
  "#bcbd22",
  "#17becf",
];
