import React, { useState } from "react";
import { HiOutlineViewList } from "react-icons/hi";
import { useConfig } from "@core/config/hooks/useConfig";
import { availableStock, canAddMore } from "@pos/features/stock/utils/stock";
import GuestlistModal from "../../../guestlist/components/GuestlistModal";
import ProductInterestModal from "./ProductInterestModal";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../../api/schemas";

/** At or below this many units left, the stock pill turns orange. */
export const LOW_STOCK_THRESHOLD = 3;

interface ProductProps {
  product: ProductType;
  addToCart: (
    product: ProductType,
    quantity: number,
    listItem: GuestType | null,
  ) => void;
  hasListItem: (guest: GuestType) => boolean;
  quantityByProductInCart: (product: ProductType) => number;
  addProductInterest: (product: ProductType) => Promise<void>;
  /** Extra classes for the card, e.g. to start a new grid row. */
  className?: string;
  /**
   * The product sold out under auto_hide and is no longer returned by the
   * backend; its slot is kept as a disabled placeholder (see vanishedProducts).
   */
  vanished?: boolean;
}

type StockTone = "normal" | "low" | "none";

const stockPillClass: Record<StockTone, string> = {
  normal: "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-200",
  low: "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-100",
  none: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-100",
};

/**
 * One product in the POS grid. The whole card is the button: one tap adds the
 * product, opens its guest list or – when sold out – asks to register interest.
 */
const Product: React.FC<ProductProps> = ({
  product,
  addToCart,
  hasListItem,
  quantityByProductInCart,
  addProductInterest,
  className = "",
  vanished = false,
}) => {
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

  if (vanished) {
    return (
      <button
        type="button"
        disabled
        data-testid={"product-card-" + product.id}
        data-vanished="true"
        aria-label={product.name + " is sold out"}
        className={`relative flex h-full flex-col justify-between rounded-xl border-2 border-dashed border-gray-300 bg-gray-50 p-3 text-left disabled:cursor-not-allowed dark:border-gray-700 dark:bg-gray-800/60 ${className}`}
      >
        <span className="text-base font-semibold leading-tight text-balance text-gray-500">
          {product.name}
        </span>
        <span className="flex items-end justify-between gap-2">
          <span className="text-2xl font-bold text-gray-400">{price}</span>
          <span className="rounded-full bg-red-100 px-2 py-0.5 text-xs font-semibold text-red-800 dark:bg-red-900 dark:text-red-100">
            Sold out
          </span>
        </span>
      </button>
    );
  }

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

  let border = "border-transparent";
  if (quantityInCart > 0) {
    border = "border-cyan-700 dark:border-cyan-500";
  } else if (product.soldOut) {
    border = "border-gray-200 dark:border-gray-700";
  }

  const muted = product.soldOut || isBlocked;

  return (
    <>
      <button
        type="button"
        data-testid={"product-card-" + product.id}
        aria-label={ariaLabel}
        disabled={isBlocked}
        onClick={handleClick}
        className={`pos-pressable relative flex h-full flex-col justify-between rounded-xl border-2 p-3 text-left shadow-sm disabled:cursor-not-allowed ${border} ${
          muted ? "bg-gray-50 dark:bg-gray-800/60" : "bg-white dark:bg-gray-800"
        } ${className}`}
      >
        <span
          className={`text-base font-semibold leading-tight text-balance ${
            muted ? "text-gray-500" : "text-gray-900 dark:text-white"
          }`}
        >
          {product.name}
        </span>

        <span className="flex items-end justify-between gap-2">
          <span
            className={`text-2xl font-bold ${
              muted ? "text-gray-500" : "text-gray-900 dark:text-white"
            }`}
          >
            {price}
          </span>

          {product.soldOut ? (
            <span className="rounded-full bg-gray-200 px-2 py-0.5 text-xs font-semibold text-gray-700 dark:bg-gray-700 dark:text-gray-200">
              Sold out · {product.soldOutRequestCount} interested
            </span>
          ) : (
            <span className="flex flex-wrap justify-end gap-1">
              {hasGuestlist && (
                <span className="flex items-center gap-1 rounded-full bg-cyan-100 px-2 py-0.5 text-xs font-semibold text-cyan-800 dark:bg-cyan-900 dark:text-cyan-100">
                  <HiOutlineViewList
                    className="h-3.5 w-3.5"
                    aria-hidden="true"
                  />
                  List
                </span>
              )}
              {hasLimitedStock && (
                <span
                  data-testid={"stock-pill-" + product.id}
                  data-tone={stockTone}
                  className={`rounded-full px-2 py-0.5 text-xs font-semibold ${stockPillClass[stockTone]}`}
                >
                  {stockLeft} left
                </span>
              )}
            </span>
          )}
        </span>

        {quantityInCart > 0 && (
          <span
            data-testid={"cart-badge-" + product.id}
            aria-hidden="true"
            className="absolute -top-2 -right-2 flex h-8 min-w-8 items-center justify-center rounded-full bg-cyan-700 px-2 text-sm font-bold text-white ring-2 ring-white dark:ring-gray-900"
          >
            {quantityInCart}
          </span>
        )}
      </button>

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
};

export default Product;
