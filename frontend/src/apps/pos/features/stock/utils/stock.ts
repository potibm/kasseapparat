import { Product } from "../../../api/schemas";
import { OutOfStockBehavior } from "@core/config/types/config.types";

/**
 * Stock rules shared by the cart and the product list. These mirror the backend's
 * definition so the POS never offers a quantity the server will reject, but the
 * backend remains the authority: it re-checks inside the purchase transaction.
 */

/** A product with no total stock is unlimited and is never restricted. */
function isUnlimited(product: Pick<Product, "totalStock">): boolean {
  return product.totalStock <= 0;
}

/** How many units of a product can still be added to the cart. */
export function availableStock(
  product: Pick<Product, "totalStock" | "unitsSold">,
  quantityInCart: number = 0,
): number {
  if (isUnlimited(product)) {
    return Number.POSITIVE_INFINITY;
  }

  return Math.max(0, product.totalStock - product.unitsSold - quantityInCart);
}

/** Whether the behaviour stops a sale that exceeds the stock. */
export function restrictsStock(behavior: OutOfStockBehavior): boolean {
  return (
    behavior === "fail" ||
    behavior === "auto_sold_out" ||
    behavior === "auto_hide"
  );
}

/**
 * Whether one more unit of a product can go into the cart.
 *
 * This is the single rule the product card and the cart share, so the add button can
 * never read as enabled while the cart would refuse the same quantity. Under "ignore"
 * the sale stays allowed, so it is always true there.
 */
export function canAddMore(
  product: Pick<Product, "totalStock" | "unitsSold">,
  quantityInCart: number,
  behavior: OutOfStockBehavior,
): boolean {
  if (!restrictsStock(behavior) || isUnlimited(product)) {
    return true;
  }

  return availableStock(product, quantityInCart) > 0;
}
