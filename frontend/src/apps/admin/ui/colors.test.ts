import { describe, expect, it } from "vitest";
import {
  PAYMENT_METHOD_COLORS,
  SERIES_COLORS,
  paymentMethodColor,
} from "./colors";

describe("paymentMethodColor", () => {
  it("gives every known payment method its own colour", () => {
    expect(paymentMethodColor("CASH")).toBe(PAYMENT_METHOD_COLORS.CASH);
    expect(paymentMethodColor("SUMUP", 3)).toBe(PAYMENT_METHOD_COLORS.SUMUP);
    expect(new Set(Object.values(PAYMENT_METHOD_COLORS)).size).toBe(4);
  });

  it("picks a fallback colour by position for unknown methods", () => {
    expect(paymentMethodColor("GIFT", 0)).toBe("#e91e63");
    expect(paymentMethodColor("GIFT", 6)).toBe("#00bcd4");
    expect(paymentMethodColor("GIFT")).toBe("#e91e63");
  });

  it("has ten series colours", () => {
    expect(SERIES_COLORS).toHaveLength(10);
  });
});
