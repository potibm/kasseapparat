import { describe, expect, it } from "vitest";
import { summarizePurchaseItems } from "./purchaseItems.utils";

describe("summarizePurchaseItems", () => {
  it("lists the products with their quantities", () => {
    expect(
      summarizePurchaseItems([
        { quantity: 2, product: { name: "Mate" } },
        { quantity: 1, product: { name: "Mug" } },
      ]),
    ).toBe("2× Mate, Mug");
  });

  it("copes with missing data", () => {
    expect(summarizePurchaseItems(undefined)).toBe("");
    expect(summarizePurchaseItems([{ quantity: 1, product: null }])).toBe(
      "Unknown product",
    );
  });
});
