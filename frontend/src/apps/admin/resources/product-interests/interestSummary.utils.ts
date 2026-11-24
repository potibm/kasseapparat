export interface ProductWithInterest {
  id: number;
  name: string;
  soldOutRequestCount?: number;
}

/** Products somebody asked for, most requested first (ties by name). */
export const productsByInterest = <P extends ProductWithInterest>(
  products: P[],
): P[] =>
  products
    .filter((product) => (product.soldOutRequestCount ?? 0) > 0)
    .sort(
      (a, b) =>
        (b.soldOutRequestCount ?? 0) - (a.soldOutRequestCount ?? 0) ||
        a.name.localeCompare(b.name),
    );
