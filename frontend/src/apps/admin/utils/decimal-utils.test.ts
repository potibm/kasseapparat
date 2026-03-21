import { describe, it, expect } from "vitest";
import { parseDecimal, formatDecimal, decimalValidator } from "./decimal-utils";

describe("Decimal Utils", () => {
  it("should parse various inputs to a clean dot-separated string", () => {
    expect(parseDecimal("12,34")).toBe("12.34");
    expect(parseDecimal("1.234,56")).toBe("1.23456");
    expect(parseDecimal("abc 12.3")).toBe("12.3");
  });

  it("should format dot-strings to comma-strings for the user", () => {
    expect(formatDecimal(10.5)).toBe("10,5");
    expect(formatDecimal("10.5")).toBe("10,5");
    expect(formatDecimal(undefined)).toBe("");
  });

  it("should validate correctly", () => {
    expect(decimalValidator("-1")).toBe("Negative number");
    expect(decimalValidator("abc")).toBe("Invalid number");
    expect(decimalValidator({})).toBe("Invalid input type");
    expect(decimalValidator("12.34")).toBeUndefined();
  });
});
