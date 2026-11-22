import { describe, expect, it } from "vitest";
import { buildSectionBlocks } from "./productBlocks";

const product = (id: number, pos: number, productGroupId?: number | null) => ({
  id,
  pos,
  productGroupId,
});

const male = { id: 10, sectionId: 1, name: "Male", pos: 20 };
const female = { id: 11, sectionId: 1, name: "Female", pos: 30 };

describe("buildSectionBlocks", () => {
  it("puts single products into one card grid", () => {
    const blocks = buildSectionBlocks([product(1, 1), product(2, 2)], []);

    expect(blocks).toEqual([
      {
        kind: "products",
        key: "products-1",
        products: [product(1, 1), product(2, 2)],
      },
    ]);
  });

  it("orders single products and groups by position and merges neighbours", () => {
    const blocks = buildSectionBlocks(
      [
        product(1, 10),
        product(2, 1, 10),
        product(3, 2, 10),
        product(4, 1, 11),
        product(5, 40),
        product(6, 5),
      ],
      [male, female],
    );

    expect(blocks.map((block) => block.kind)).toEqual([
      "products",
      "groups",
      "products",
    ]);

    const [first, groups, last] = blocks;
    expect(first.kind === "products" && first.products).toEqual([
      product(6, 5),
      product(1, 10),
    ]);
    expect(
      groups.kind === "groups" &&
        groups.groups.map((entry) => [
          entry.group.name,
          entry.products.map((p) => p.id),
        ]),
    ).toEqual([
      ["Male", [2, 3]],
      ["Female", [4]],
    ]);
    expect(last.kind === "products" && last.products).toEqual([product(5, 40)]);
  });

  it("shows products of an unknown group as single products", () => {
    const blocks = buildSectionBlocks(
      [product(1, 1, 99), product(2, 2, null)],
      [male],
    );

    expect(blocks).toEqual([
      {
        kind: "products",
        key: "products-1",
        products: [product(1, 1, 99), product(2, 2, null)],
      },
    ]);
  });

  it("leaves out groups without products", () => {
    expect(buildSectionBlocks([], [male, female])).toEqual([]);
  });
});
