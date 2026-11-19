import React, { useEffect, useState, useRef } from "react";
import {
  Table,
  Tooltip,
  TableHeadCell,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from "flowbite-react";
import { useConfig } from "@core/config/hooks/useConfig";
import "animate.css";
import Button from "../../../components/Button";
import CheckoutButtons from "./_internal/CheckoutButtons";
import { Cart as CartObject } from "../services/Cart";
import {
  CartItem as CartItemType,
  PaymentMethodData as PaymentMethodDataType,
} from "../types/cart.types";
import { Product as ProductType } from "../../../api/schemas.ts";
import CartRow from "./_internal/CartRow.tsx";

interface CartProps {
  cart: CartObject;
  removeFromCart: (item: ProductType) => void;
  removeAllFromCart: () => void;
  checkoutCart: (
    paymentMethodCode: string,
    paymentMethodData: PaymentMethodDataType,
  ) => Promise<void>;
  checkoutProcessing: string | null;
}

const compactTableTheme = {
  head: {
    cell: {
      base: "px-2 py-1",
    },
  },
  body: {
    cell: {
      base: "px-2 py-1",
    },
  },
};

/**
 * The cart fills the sidebar as a flex column: a header with the "Clear" action,
 * the line items in the middle (the only part that scrolls) and the total plus the
 * checkout buttons pinned to the bottom, so they stay within reach no matter how
 * many items are in the cart. The parent must give it a bounded height.
 */
const Cart: React.FC<CartProps> = ({
  cart,
  removeFromCart,
  removeAllFromCart,
  checkoutCart,
  checkoutProcessing,
}) => {
  const { currency } = useConfig();

  const [flash, setFlash] = useState(false);

  const prevCartTotalQuantityRef = useRef(cart.totalQuantity);
  const isFirstRenderRef = useRef(true);

  useEffect(() => {
    if (isFirstRenderRef.current) {
      isFirstRenderRef.current = false;
      return;
    }

    let rafId: number;
    let timerId: ReturnType<typeof setTimeout>;

    if (cart.totalQuantity !== prevCartTotalQuantityRef.current) {
      rafId = requestAnimationFrame(() => {
        setFlash(true);
        timerId = setTimeout(() => setFlash(false), 500);
      });
    }

    prevCartTotalQuantityRef.current = cart.totalQuantity;

    return () => {
      if (rafId) cancelAnimationFrame(rafId);
      if (timerId) clearTimeout(timerId);
    };
  }, [cart.totalQuantity]);

  return (
    <div className="flex min-h-0 flex-1 flex-col" data-testid="cart">
      <div className="flex shrink-0 items-center justify-between pb-2">
        <h2 className="flex items-center gap-2 text-lg font-bold dark:text-white">
          Cart
          <span
            className="rounded-full bg-gray-300 px-2 text-sm font-semibold dark:bg-gray-700"
            data-testid="cart-quantity"
          >
            {cart.totalQuantity}
          </span>
        </h2>
        <Button
          color="alternative"
          className="text-red-700 dark:text-red-400"
          disabled={cart.isEmpty}
          onClick={() => removeAllFromCart()}
          aria-label="Remove all items from cart"
        >
          Clear
        </Button>
      </div>

      <div className="min-h-0 flex-1 overflow-y-auto" data-testid="cart-items">
        <Table
          striped
          theme={compactTableTheme}
          data-testid="cart-table"
          className={`table-fixed dark:text-gray-200 ${flash ? "animate__animated animate__pulse" : ""}`}
        >
          <TableHead>
            <TableRow>
              <TableHeadCell className="w-[40%]">Product</TableHeadCell>
              <TableHeadCell className="w-[15%] text-right">
                <Tooltip content="Quantity">Qnt</Tooltip>
              </TableHeadCell>
              <TableHeadCell className="w-[15%] text-right">
                Total Price
              </TableHeadCell>
              <TableHeadCell className="w-[30%] text-right">
                Remove
              </TableHeadCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {cart.isEmpty && (
              <TableRow>
                <TableCell colSpan={4} className="py-4 text-center">
                  Cart is empty
                </TableCell>
              </TableRow>
            )}
            {cart.items.map((cartElement: CartItemType) => (
              <CartRow
                key={cartElement.id}
                cartElement={cartElement}
                currency={currency}
                removeFromCart={removeFromCart}
              />
            ))}
          </TableBody>
        </Table>
      </div>

      <div
        className="shrink-0 border-t border-gray-300 pt-2 dark:border-gray-700"
        data-testid="cart-footer"
      >
        <div className="flex items-baseline justify-between px-2 dark:text-white">
          <span className="text-sm font-bold uppercase">Total</span>
          <span className="text-2xl font-bold" data-testid="cart-total-value">
            {currency.format(cart.totalGross.toNumber())}
          </span>
        </div>

        <CheckoutButtons
          cart={cart}
          checkoutProcessing={checkoutProcessing}
          handleCheckoutCart={checkoutCart}
        />
      </div>
    </div>
  );
};

export default Cart;
