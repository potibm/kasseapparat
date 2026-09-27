import React from "react";
import { TableRow, TableCell } from "flowbite-react";
import { HiXCircle } from "react-icons/hi";
import Button from "../../../../components/Button";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../../api/schemas";
import { CartItem as CartItemType } from "../../types/cart.types";
import { availableStock } from "@pos/features/stock/utils/stock";

interface CartRowProps {
  cartElement: CartItemType;
  currency: { format: (val: number) => string };
  removeFromCart: (item: ProductType) => void;
  isOverStock: boolean;
}

const CartRow: React.FC<CartRowProps> = ({
  cartElement,
  currency,
  removeFromCart,
  isOverStock,
}) => {
  const displayListItem = (listItem: GuestType) => {
    return listItem.code ?? listItem.name;
  };

  const available = availableStock(cartElement, 0);
  const overBy = cartElement.quantity - available;

  return (
    <TableRow
      key={cartElement.id}
      data-testid={"cart-product-" + cartElement.id}
      className={isOverStock ? "bg-red-100 dark:bg-red-900/40" : undefined}
    >
      <TableCell className="whitespace-normal px-4 py-2">
        {cartElement.name}
        {cartElement.listItems.map((listItem: GuestType) => (
          <div key={listItem.id} className="text-xs text-gray-500">
            {displayListItem(listItem)}
          </div>
        ))}
        {isOverStock && (
          <div
            className="text-xs font-medium text-red-700 dark:text-red-300"
            data-testid={"cart-over-stock-" + cartElement.id}
          >
            {overBy > 0
              ? `Only ${available} left, ${overBy} over stock`
              : `Only ${available} left`}
          </div>
        )}
      </TableCell>
      <TableCell className="text-right">{cartElement.quantity}</TableCell>
      <TableCell className="text-right">
        {currency.format(cartElement.totalGrossPrice.toNumber())}
      </TableCell>
      <TableCell className="flex justify-end">
        <Button
          color="failure"
          aria-label={`Remove ${cartElement.name} from cart`}
          onClick={() => removeFromCart(cartElement)}
        >
          <HiXCircle />
        </Button>
      </TableCell>
    </TableRow>
  );
};

export default CartRow;
