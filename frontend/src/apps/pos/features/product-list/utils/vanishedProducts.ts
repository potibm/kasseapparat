import { Product } from "../../../api/schemas";
import { OutOfStockBehavior } from "@core/config/types/config.types";

/** A product as shown in the grid; `vanished` marks a kept placeholder slot. */
export type DisplayProduct = Product & { vanished?: boolean };

/**
 * Under `auto_hide` the backend stops returning a product as soon as its stock
 * is sold. Simply dropping it would shift every following card by one slot –
 * right under the cashier's finger, so the next quick tap hits the wrong
 * product. Instead the product keeps its slot as a disabled "sold out"
 * placeholder until the POS is reloaded.
 *
 * Only products with limited stock are kept: an unlimited product cannot sell
 * out, so if it disappears an admin hid it on purpose.
 */
export function keepVanishedProducts(
  previous: DisplayProduct[] | null,
  fetched: Product[],
  behavior: OutOfStockBehavior,
): DisplayProduct[] {
  if (!previous || behavior !== "auto_hide") {
    return fetched;
  }

  const fetchedIds = new Set(fetched.map((product) => product.id));
  const result: DisplayProduct[] = [...fetched];

  // Walk the previous order, so each placeholder goes back to its old index.
  previous.forEach((product, index) => {
    if (fetchedIds.has(product.id)) return;
    if (!product.vanished && product.totalStock <= 0) return;

    result.splice(Math.min(index, result.length), 0, {
      ...product,
      vanished: true,
    });
  });

  return result;
}
