import { Section } from "@pos/api/schemas";

export interface ProductSection<P> {
  /** Stable React key. */
  key: string;
  /** Heading shown above the products; empty for a plain row break. */
  name: string;
  products: P[];
}

/**
 * Splits the products into the sections they belong to (ADR 006).
 *
 * Sections are ordered by position, products keep the order they came in. Products
 * whose section is unknown (for example because the sections could not be loaded)
 * end up in a trailing group without a heading. Empty sections are left out.
 */
export const groupProductsBySection = <P extends { sectionId: number }>(
  products: P[],
  sections: Section[],
): ProductSection<P>[] => {
  const knownSectionIds = new Set(sections.map((section) => section.id));

  const groups: ProductSection<P>[] = [...sections]
    .sort((a, b) => a.pos - b.pos || a.id - b.id)
    .map((section) => ({
      key: `section-${section.id}`,
      name: section.name.trim(),
      products: products.filter((product) => product.sectionId === section.id),
    }));

  const unassigned = products.filter(
    (product) => !knownSectionIds.has(product.sectionId),
  );
  if (unassigned.length > 0) {
    groups.push({ key: "section-none", name: "", products: unassigned });
  }

  return groups.filter((group) => group.products.length > 0);
};
