import React from "react";
import { HiOutlineViewList } from "react-icons/hi";
import { useProductTap, ProductTapProps, StockTone } from "./useProductTap";

interface ProductProps extends ProductTapProps {
  /** Extra classes for the card, e.g. to start a new grid row. */
  className?: string;
}

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
}) => {
  const {
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
  } = useProductTap({
    product,
    addToCart,
    hasListItem,
    quantityByProductInCart,
    addProductInterest,
  });

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

      {modals}
    </>
  );
};

export default Product;
