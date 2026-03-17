import React from "react";
import { TableRow, TableCell } from "flowbite-react";
import { HiXCircle } from "react-icons/hi";
import Button from "../../../../components/Button";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../../utils/api.schemas";
import { CartItem as CartItemType } from "../../types/cart.types";

interface CartRowProps {
  cartElement: CartItemType;
  currency: { format: (val: number) => string };
  removeFromCart: (item: ProductType) => void;
}

const CartRow: React.FC<CartRowProps> = ({
  cartElement,
  currency,
  removeFromCart,
}) => {
  const displayListItem = (listItem: GuestType) => {
    return listItem.code ?? listItem.name;
  };

  return (
    <TableRow key={cartElement.id}>
      <TableCell className="whitespace-normal px-4 py-2">
        {cartElement.name}
        {cartElement.listItems.map((listItem: GuestType) => (
          <div key={listItem.id} className="text-xs text-gray-500">
            {displayListItem(listItem)}
          </div>
        ))}
      </TableCell>
      <TableCell className="text-right">{cartElement.quantity}</TableCell>
      <TableCell className="text-right">
        {currency.format(cartElement.totalGrossPrice.toNumber())}
      </TableCell>
      <TableCell className="flex justify-end">
        <Button color="failure" onClick={() => removeFromCart(cartElement)}>
          <HiXCircle />
        </Button>
      </TableCell>
    </TableRow>
  );
};

export default CartRow;
