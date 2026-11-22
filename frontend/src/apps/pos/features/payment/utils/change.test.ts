import { describe, expect, it } from "vitest";
import Decimal from "decimal.js";
import { parseAmount, suggestReceivedAmounts } from "./change";

const amounts = (total: number) =>
  suggestReceivedAmounts(new Decimal(total)).map((amount) => amount.toNumber());

describe("suggestReceivedAmounts", () => {
  it("suggests the next round sums without duplicates", () => {
    expect(amounts(37)).toEqual([40, 50, 100, 200]);
  });

  it("only suggests amounts above the total", () => {
    expect(amounts(40)).toEqual([45, 50, 60, 100]);
  });

  it("works with cents", () => {
    expect(amounts(3.5)).toEqual([5, 10, 20, 50]);
  });

  it("can be limited", () => {
    expect(
      suggestReceivedAmounts(new Decimal(12), 2).map((a) => a.toNumber()),
    ).toEqual([15, 20]);
  });
});

describe("parseAmount", () => {
  it("accepts dots and commas", () => {
    expect(parseAmount("12.50")?.toNumber()).toBe(12.5);
    expect(parseAmount(" 12,5 ")?.toNumber()).toBe(12.5);
    expect(parseAmount("50")?.toNumber()).toBe(50);
  });

  it("rejects everything else", () => {
    expect(parseAmount("")).toBeNull();
    expect(parseAmount("abc")).toBeNull();
    expect(parseAmount("-5")).toBeNull();
    expect(parseAmount("1,2,3")).toBeNull();
  });
});
