import { ProductGroup } from "@pos/api/schemas";

export interface GroupWithProducts<P> {
  group: ProductGroup;
  products: P[];
}

export type SectionBlock<P> =
  | { kind: "products"; key: string; products: P[] }
  | { kind: "groups"; key: string; groups: GroupWithProducts<P>[] };

type Item<P> =
  | { kind: "product"; pos: number; product: P }
  | { kind: "group"; pos: number; entry: GroupWithProducts<P> };

interface GroupableProduct {
  id: number;
  pos: number;
  productGroupId?: number | null;
}

/**
 * Lays out the products of one section (ADR 006). Single products and product groups
 * share one order by position. Consecutive single products form a card grid,
 * consecutive groups a size matrix with one row per group.
 *
 * Products keep the order they came in within their group. A product whose group is
 * unknown (for example because the groups could not be loaded) is shown as a single
 * product, so it stays sellable.
 */
export const buildSectionBlocks = <P extends GroupableProduct>(
  products: P[],
  groups: ProductGroup[],
): SectionBlock<P>[] => {
  const groupById = new Map(groups.map((group) => [group.id, group]));
  const entries = new Map<number, GroupWithProducts<P>>();
  const items: Item<P>[] = [];

  for (const product of products) {
    addItem(items, entries, groupById, product);
  }

  // Array.prototype.sort is stable: equal positions keep the order of the products.
  items.sort((a, b) => a.pos - b.pos);

  const blocks: SectionBlock<P>[] = [];
  for (const item of items) {
    appendBlock(blocks, item);
  }

  return blocks;
};

const addItem = <P extends GroupableProduct>(
  items: Item<P>[],
  entries: Map<number, GroupWithProducts<P>>,
  groupById: Map<number, ProductGroup>,
  product: P,
): void => {
  const group =
    product.productGroupId == null
      ? undefined
      : groupById.get(product.productGroupId);

  if (!group) {
    items.push({ kind: "product", pos: product.pos, product });
    return;
  }

  let entry = entries.get(group.id);
  if (!entry) {
    entry = { group, products: [] };
    entries.set(group.id, entry);
    items.push({ kind: "group", pos: group.pos, entry });
  }
  entry.products.push(product);
};

const appendBlock = <P extends GroupableProduct>(
  blocks: SectionBlock<P>[],
  item: Item<P>,
): void => {
  const last = blocks.at(-1);

  if (item.kind === "product") {
    if (last?.kind === "products") {
      last.products.push(item.product);
      return;
    }
    blocks.push({
      kind: "products",
      key: "products-" + item.product.id,
      products: [item.product],
    });
    return;
  }

  if (last?.kind === "groups") {
    last.groups.push(item.entry);
    return;
  }

  blocks.push({
    kind: "groups",
    key: "groups-" + item.entry.group.id,
    groups: [item.entry],
  });
};
