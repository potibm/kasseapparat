export interface PurchaseItemLike {
  quantity: number;
  product?: { name?: string } | null;
}

/** "2× Mate, Beer, T-Shirt" – the same summary the POS history shows. */
export const summarizePurchaseItems = (
  items: PurchaseItemLike[] | null | undefined,
): string =>
  (items ?? [])
    .map((item) => {
      const name = item.product?.name || "Unknown product";
      return item.quantity > 1 ? `${item.quantity}× ${name}` : name;
    })
    .join(", ");
