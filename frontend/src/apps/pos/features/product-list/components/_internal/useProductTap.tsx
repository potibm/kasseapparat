import { useState } from "react";
import { useConfig } from "@core/config/hooks/useConfig";
import { availableStock, canAddMore } from "@pos/features/stock/utils/stock";
import GuestlistModal from "../../../guestlist/components/GuestlistModal";
import ProductInterestModal from "./ProductInterestModal";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../../api/schemas";

/** At or below this many units left, the stock turns orange. */
export const LOW_STOCK_THRESHOLD = 3;

export type StockTone = "normal" | "low" | "none";

export interface ProductTapProps {
  product: ProductType;
  addToCart: (
    product: ProductType,
    quantity: number,
    listItem: GuestType | null,
  ) => void;
  hasListItem: (guest: GuestType) => boolean;
  quantityByProductInCart: (product: ProductType) => number;
  addProductInterest: (product: ProductType) => Promise<void>;
}

/**
 * What tapping a product does, shared by the product card and the variant buttons of
 * a product group: add one unit, open the guest list or – when sold out – ask to
 * register interest. `modals` has to be rendered next to the button.
 */
export const useProductTap = ({
  product,
  addToCart,
  hasListItem,
  quantityByProductInCart,
  addProductInterest,
}: ProductTapProps) => {
  const [isGuestListModalOpen, setIsGuestListModalOpen] = useState(false);
  const [isPIModalOpen, setIsPIModalOpen] = useState(false);
  const { currency, outOfStockBehavior } = useConfig();

  const quantityInCart = quantityByProductInCart(product);
  const stockLeft = availableStock(product, quantityInCart);
  const hasLimitedStock = Number.isFinite(stockLeft);
  const hasGuestlist = !!product.guestlists && product.guestlists.length > 0;

  // Whether another unit still fits. Under the "ignore" behaviour the sale stays
  // allowed, so nothing is ever disabled there.
  const canAdd = canAddMore(product, quantityInCart, outOfStockBehavior);

  // The guest list modal does its own stock checks, so it stays reachable.
  const isBlocked = !product.soldOut && !hasGuestlist && !canAdd;

  const price = currency.format(product.grossPrice.toNumber());

  let ariaLabel: string;
  if (product.soldOut) {
    ariaLabel = "Register interest in " + product.name;
  } else if (hasGuestlist) {
    ariaLabel = "Show guestlist for " + product.name;
  } else if (isBlocked) {
    ariaLabel = "No " + product.name + " left";
  } else {
    ariaLabel = "Add " + product.name + " for " + price + " to cart";
  }

  let stockTone: StockTone = "normal";
  if (stockLeft <= 0) {
    stockTone = "none";
  } else if (stockLeft <= LOW_STOCK_THRESHOLD) {
    stockTone = "low";
  }

  const handleClick = () => {
    if (product.soldOut) {
      setIsPIModalOpen(true);
    } else if (hasGuestlist) {
      setIsGuestListModalOpen(true);
    } else if (canAdd) {
      addToCart(product, 1, null);
    }
  };

  const modals = (
    <>
      {!product.soldOut && hasGuestlist && (
        <GuestlistModal
          isOpen={isGuestListModalOpen}
          onClose={() => setIsGuestListModalOpen(false)}
          product={product}
          addToCart={addToCart}
          hasListItem={hasListItem}
          quantityInCart={quantityInCart}
        />
      )}

      {product.soldOut && (
        <ProductInterestModal
          show={isPIModalOpen}
          onClose={() => setIsPIModalOpen(false)}
          product={product}
          addProductInterest={addProductInterest}
        />
      )}
    </>
  );

  return {
    quantityInCart,
    stockLeft,
    hasLimitedStock,
    hasGuestlist,
    isBlocked,
    price,
    ariaLabel,
    stockTone,
    handleClick,
    modals,
  };
};
