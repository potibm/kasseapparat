import React from "react";
import { useProductTap, ProductTapProps, StockTone } from "./useProductTap";

interface VariantButtonProps extends ProductTapProps {
  /** Show the price on the button; off when the whole row shares one price. */
  showPrice: boolean;
  /** Sold out under auto_hide and kept as a placeholder (see vanishedProducts). */
  vanished?: boolean;
}

const stockTextClass: Record<StockTone, string> = {
  normal: "text-gray-600 dark:text-gray-300",
  low: "text-orange-700 dark:text-orange-300",
  none: "text-red-700 dark:text-red-300",
};

/**
 * One variant in the size matrix of a product group, e.g. "XL". It behaves like a
 * product card: one tap adds the variant, opens its guest list or – when sold out –
 * asks to register interest.
 */
const VariantButton: React.FC<VariantButtonProps> = ({
  showPrice,
  vanished = false,
  ...tapProps
}) => {
  const { product } = tapProps;
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
  } = useProductTap(tapProps);

  const label = product.variant?.trim() || product.name;
  const unavailable = vanished || product.soldOut || isBlocked;

  let status: string | null = null;
  if (vanished || product.soldOut) {
    status = "sold out";
  } else if (hasLimitedStock) {
    status = stockLeft + " left";
  } else if (hasGuestlist) {
    status = "List";
  }

  let border = "border-gray-300 dark:border-gray-600";
  if (quantityInCart > 0) {
    border = "border-cyan-700 dark:border-cyan-500";
  } else if (unavailable) {
    border = "border-gray-200 dark:border-gray-700";
  }

  return (
    <>
      <button
        type="button"
        data-testid={"product-card-" + product.id}
        data-vanished={vanished ? "true" : undefined}
        aria-label={vanished ? product.name + " is sold out" : ariaLabel}
        disabled={vanished || isBlocked}
        onClick={handleClick}
        className={`pos-pressable relative flex h-[4.5rem] flex-col items-center justify-center gap-0.5 rounded-lg border-2 px-1 disabled:cursor-not-allowed ${border} ${
          unavailable
            ? "bg-gray-100 text-gray-500 dark:bg-gray-800/60"
            : "bg-white text-gray-900 dark:bg-gray-800 dark:text-white"
        } ${vanished ? "border-dashed" : ""}`}
      >
        <span
          className={`max-w-full truncate text-xl leading-none font-bold ${
            unavailable ? "line-through" : ""
          }`}
        >
          {label}
        </span>
        {showPrice && !unavailable && (
          <span className="text-xs font-semibold">{price}</span>
        )}
        {status && (
          <span
            data-testid={"variant-status-" + product.id}
            className={`text-xs font-medium ${
              unavailable ? "text-gray-500" : stockTextClass[stockTone]
            }`}
          >
            {status}
          </span>
        )}
        {quantityInCart > 0 && (
          <span
            data-testid={"cart-badge-" + product.id}
            aria-hidden="true"
            className="absolute -top-2 -right-2 flex h-7 min-w-7 items-center justify-center rounded-full bg-cyan-700 px-1.5 text-xs font-bold text-white ring-2 ring-white dark:ring-gray-900"
          >
            {quantityInCart}
          </span>
        )}
      </button>
      {!vanished && modals}
    </>
  );
};

export default VariantButton;
