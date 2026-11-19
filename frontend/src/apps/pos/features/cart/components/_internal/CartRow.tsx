import React from "react";
import { TableRow, TableCell } from "flowbite-react";
import { HiMinus, HiPlus, HiTrash } from "react-icons/hi";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../../api/schemas";
import { CartItem as CartItemType } from "../../types/cart.types";
import { OutOfStockBehavior } from "@core/config/types/config.types";
import { canAddMore } from "@pos/features/stock/utils/stock";

interface CartRowProps {
  cartElement: CartItemType;
  currency: { format: (val: number) => string };
  outOfStockBehavior: OutOfStockBehavior;
  removeFromCart: (item: ProductType) => void;
  increaseInCart: (item: ProductType) => void;
  decreaseInCart: (item: ProductType) => void;
}

const stepperButton =
  "flex h-11 w-12 items-center justify-center border border-gray-300 bg-gray-50 text-gray-900 disabled:cursor-not-allowed disabled:opacity-40 dark:border-gray-600 dark:bg-gray-800 dark:text-white";

/**
 * One cart line: name and line total on top, the quantity stepper below.
 *
 * At quantity 1 the minus button turns into a trash button that removes the line.
 * Lines with guest-list entries only get the trash button: every unit belongs to a
 * guest, so they cannot be changed unit by unit.
 */
const CartRow: React.FC<CartRowProps> = ({
  cartElement,
  currency,
  outOfStockBehavior,
  removeFromCart,
  increaseInCart,
  decreaseInCart,
}) => {
  const displayListItem = (listItem: GuestType) => {
    return listItem.code ?? listItem.name;
  };

  const hasListItems = cartElement.listItems.length > 0;
  const removesLine = hasListItems || cartElement.quantity <= 1;
  const canIncrease =
    !hasListItems &&
    canAddMore(cartElement, cartElement.quantity, outOfStockBehavior);

  return (
    <TableRow data-testid={"cart-product-" + cartElement.id}>
      <TableCell className="whitespace-normal px-2 py-3 text-left">
        <div className="flex justify-between gap-2">
          <div className="min-w-0">
            <div className="font-semibold text-gray-900 dark:text-white">
              {cartElement.name}
            </div>
            {cartElement.listItems.map((listItem: GuestType) => (
              <div key={listItem.id} className="text-xs text-gray-500">
                {displayListItem(listItem)}
              </div>
            ))}
          </div>
          <div className="whitespace-nowrap font-semibold text-gray-900 dark:text-white">
            {currency.format(cartElement.totalGrossPrice.toNumber())}
          </div>
        </div>

        <div className="mt-2 flex items-center">
          {removesLine ? (
            <button
              type="button"
              className={`${stepperButton} rounded-l-lg text-red-700 dark:text-red-400`}
              aria-label={`Remove ${cartElement.name} from cart`}
              onClick={() => removeFromCart(cartElement)}
            >
              <HiTrash className="h-5 w-5" />
            </button>
          ) : (
            <button
              type="button"
              className={`${stepperButton} rounded-l-lg`}
              aria-label={`Remove one ${cartElement.name}`}
              onClick={() => decreaseInCart(cartElement)}
            >
              <HiMinus className="h-5 w-5" />
            </button>
          )}
          <span
            className="flex h-11 w-12 items-center justify-center border-y border-gray-300 text-lg font-bold text-gray-900 dark:border-gray-600 dark:text-white"
            data-testid={"cart-quantity-" + cartElement.id}
          >
            {cartElement.quantity}
          </span>
          <button
            type="button"
            className={`${stepperButton} rounded-r-lg`}
            aria-label={`Add one ${cartElement.name}`}
            disabled={!canIncrease}
            onClick={() => increaseInCart(cartElement)}
          >
            <HiPlus className="h-5 w-5" />
          </button>
          <span className="ml-2 text-xs text-gray-500 dark:text-gray-400">
            × {currency.format(cartElement.grossPrice.toNumber())}
          </span>
        </div>
      </TableCell>
    </TableRow>
  );
};

export default CartRow;
