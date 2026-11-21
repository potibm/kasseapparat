import { describe, it, expect } from "vitest";
import { createMockProduct } from "@pos/api/schemas.mocks";
import { keepVanishedProducts, DisplayProduct } from "./vanishedProducts";

const limited = (id: number) =>
  createMockProduct({ id, name: "Limited " + id, totalStock: 10 });
const unlimited = (id: number) =>
  createMockProduct({ id, name: "Unlimited " + id, totalStock: 0 });

const ids = (products: DisplayProduct[]) =>
  products.map((p) => (p.vanished ? `${p.id}*` : `${p.id}`));

describe("keepVanishedProducts", () => {
  it("returns the fetched list on the first load", () => {
    const fetched = [limited(1), limited(2)];

    expect(keepVanishedProducts(null, fetched, "auto_hide")).toBe(fetched);
  });

  it.each(["ignore", "fail", "auto_sold_out"] as const)(
    "does not keep anything under %s",
    (behavior) => {
      const previous = [limited(1), limited(2)];
      const fetched = [limited(2)];

      expect(keepVanishedProducts(previous, fetched, behavior)).toBe(fetched);
    },
  );

  it("keeps a sold-out product as a placeholder in its old slot", () => {
    const previous = [limited(1), limited(2), limited(3)];
    const fetched = [limited(1), limited(3)];

    const result = keepVanishedProducts(previous, fetched, "auto_hide");

    expect(ids(result)).toEqual(["1", "2*", "3"]);
  });

  it("keeps several placeholders in their old slots", () => {
    const previous = [limited(1), limited(2), limited(3), limited(4)];
    const fetched = [limited(2), limited(4)];

    const result = keepVanishedProducts(previous, fetched, "auto_hide");

    expect(ids(result)).toEqual(["1*", "2", "3*", "4"]);
  });

  it("keeps a placeholder at the end of the grid", () => {
    const previous = [limited(1), limited(2)];
    const fetched = [limited(1)];

    expect(ids(keepVanishedProducts(previous, fetched, "auto_hide"))).toEqual([
      "1",
      "2*",
    ]);
  });

  it("drops an unlimited product, which an admin must have hidden", () => {
    const previous = [limited(1), unlimited(2), limited(3)];
    const fetched = [limited(1), limited(3)];

    expect(ids(keepVanishedProducts(previous, fetched, "auto_hide"))).toEqual([
      "1",
      "3",
    ]);
  });

  it("keeps an existing placeholder across further refreshes", () => {
    const first = keepVanishedProducts(
      [limited(1), limited(2), limited(3)],
      [limited(1), limited(3)],
      "auto_hide",
    );

    const second = keepVanishedProducts(
      first,
      [limited(1), limited(3)],
      "auto_hide",
    );

    expect(ids(second)).toEqual(["1", "2*", "3"]);
  });

  it("replaces the placeholder when the product comes back", () => {
    const withPlaceholder = keepVanishedProducts(
      [limited(1), limited(2), limited(3)],
      [limited(1), limited(3)],
      "auto_hide",
    );

    const restocked = keepVanishedProducts(
      withPlaceholder,
      [limited(1), limited(2), limited(3)],
      "auto_hide",
    );

    expect(ids(restocked)).toEqual(["1", "2", "3"]);
  });

  it("uses the fresh data of products that are still there", () => {
    const previous = [
      createMockProduct({ id: 1, totalStock: 10, unitsSold: 2 }),
    ];
    const fresh = createMockProduct({ id: 1, totalStock: 10, unitsSold: 5 });

    const [result] = keepVanishedProducts(previous, [fresh], "auto_hide");

    expect(result).toBe(fresh);
  });
});
