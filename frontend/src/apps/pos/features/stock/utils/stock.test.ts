import { describe, it, expect } from "vitest";
import { availableStock, restrictsStock, canAddMore } from "./stock";
import { Product } from "../../../api/schemas";

const product = (overrides: Partial<Product> = {}): Product =>
  ({
    id: 1,
    name: "T-Shirt",
    netPrice: "10",
    grossPrice: "12.5",
    vatRate: "25",
    vatAmount: "2.5",
    wrapAfter: false,
    hidden: false,
    soldOut: false,
    apiExport: false,
    pos: 1,
    totalStock: 5,
    unitsSold: 0,
    soldOutRequestCount: 0,
    guestlists: null,
    ...overrides,
  }) as Product;

describe("availableStock", () => {
  it("is unlimited for a product without a stock limit", () => {
    expect(availableStock(product({ totalStock: 0, unitsSold: 99 }))).toBe(
      Number.POSITIVE_INFINITY,
    );
  });

  it("subtracts what was sold and what is in the cart", () => {
    expect(availableStock(product({ totalStock: 10, unitsSold: 4 }), 2)).toBe(
      4,
    );
  });

  it("never goes below zero", () => {
    expect(availableStock(product({ totalStock: 2, unitsSold: 5 }))).toBe(0);
  });
});

describe("restrictsStock", () => {
  it("does not restrict under ignore", () => {
    expect(restrictsStock("ignore")).toBe(false);
  });

  it.each(["fail", "auto_sold_out", "auto_hide"] as const)(
    "restricts under %s",
    (behavior) => {
      expect(restrictsStock(behavior)).toBe(true);
    },
  );
});

describe("canAddMore", () => {
  it("is always true under ignore, even past the stock", () => {
    expect(
      canAddMore(product({ totalStock: 2, unitsSold: 9 }), 99, "ignore"),
    ).toBe(true);
  });

  it("is true for an unlimited product in every enforcing mode", () => {
    for (const behavior of ["fail", "auto_sold_out", "auto_hide"] as const) {
      expect(
        canAddMore(product({ totalStock: 0, unitsSold: 500 }), 20, behavior),
      ).toBe(true);
    }
  });

  it.each(["fail", "auto_sold_out", "auto_hide"] as const)(
    "is true while stock remains under %s",
    (behavior) => {
      expect(
        canAddMore(product({ totalStock: 3, unitsSold: 0 }), 1, behavior),
      ).toBe(true);
    },
  );

  it.each(["fail", "auto_sold_out", "auto_hide"] as const)(
    "is false once the cart holds the remaining stock under %s",
    (behavior) => {
      // One unit left, one already in the cart.
      expect(
        canAddMore(product({ totalStock: 1, unitsSold: 0 }), 1, behavior),
      ).toBe(false);
    },
  );

  it.each(["fail", "auto_sold_out", "auto_hide"] as const)(
    "is false once the stock itself is gone under %s",
    (behavior) => {
      expect(
        canAddMore(product({ totalStock: 3, unitsSold: 3 }), 0, behavior),
      ).toBe(false);
    },
  );
});
