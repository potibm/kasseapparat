import { describe, expect, it } from "vitest";
import { groupProductsBySection } from "./productSections";

const product = (id: number, sectionId: number) => ({ id, sectionId });

describe("groupProductsBySection", () => {
  it("orders sections by position and keeps the product order", () => {
    const groups = groupProductsBySection(
      [product(1, 2), product(2, 1), product(3, 2)],
      [
        { id: 2, name: "Merch", pos: 2 },
        { id: 1, name: " Tickets ", pos: 1 },
      ],
    );

    expect(groups).toEqual([
      { key: "section-1", name: "Tickets", products: [product(2, 1)] },
      {
        key: "section-2",
        name: "Merch",
        products: [product(1, 2), product(3, 2)],
      },
    ]);
  });

  it("breaks ties in position by id", () => {
    const groups = groupProductsBySection(
      [product(1, 5), product(2, 4)],
      [
        { id: 5, name: "B", pos: 1 },
        { id: 4, name: "A", pos: 1 },
      ],
    );

    expect(groups.map((g) => g.key)).toEqual(["section-4", "section-5"]);
  });

  it("leaves out empty sections", () => {
    const groups = groupProductsBySection(
      [product(1, 1)],
      [
        { id: 1, name: "Tickets", pos: 1 },
        { id: 2, name: "Merch", pos: 2 },
      ],
    );

    expect(groups.map((g) => g.key)).toEqual(["section-1"]);
  });

  it("collects products of unknown sections in a trailing group without a name", () => {
    const groups = groupProductsBySection(
      [product(1, 9), product(2, 1), product(3, 0)],
      [{ id: 1, name: "Tickets", pos: 1 }],
    );

    expect(groups).toEqual([
      { key: "section-1", name: "Tickets", products: [product(2, 1)] },
      {
        key: "section-none",
        name: "",
        products: [product(1, 9), product(3, 0)],
      },
    ]);
  });

  it("puts everything in one group when there are no sections", () => {
    expect(groupProductsBySection([product(1, 1)], [])).toEqual([
      { key: "section-none", name: "", products: [product(1, 1)] },
    ]);
  });
});
