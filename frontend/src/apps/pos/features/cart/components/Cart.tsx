import React, { useEffect, useState, useRef } from "react";
import { HiXCircle } from "react-icons/hi";
import {
  Table,
  Tooltip,
  TableHeadCell,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from "flowbite-react";
import { useConfig } from "../../../../../core/config/providers/ConfigProvider";
import "animate.css";
import Button from "../../../components/Button";
import CheckoutButtons from "./_internal/CheckoutButtons";
import { Cart as CartObject } from "../services/Cart";
import {
  CartItem as CartItemType,
  PaymentMethodData as PaymentMethodDataType,
} from "../types/cart.types";
import { Product as ProductType } from "../../../utils/api.schemas";
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

const Cart: React.FC<CartProps> = ({
  cart,
  removeFromCart,
  removeAllFromCart,
  checkoutCart,
  checkoutProcessing,
}) => {
  const { currency } = useConfig();

  const [flash, setFlash] = useState(false);

  const prevCartTotalQuantity = useRef(cart.totalQuantity);
  const isFirstRender = useRef(true);

  const triggerFlash = () => {
    requestAnimationFrame(() => {
      setFlash(true);
      setTimeout(() => {
        setFlash(false);
      }, 500);
    });
  };

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }

    if (cart.totalQuantity !== prevCartTotalQuantity.current) {
      triggerFlash();
    }

    prevCartTotalQuantity.current = cart.totalQuantity;
  }, [cart]);

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

  return (
    <div>
      <Table
        striped
        theme={compactTableTheme}
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
            <TableHeadCell className="w-[30%] text-right">Remove</TableHeadCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {cart.items.map((cartElement: CartItemType) => (
            <CartRow
              key={cartElement.id}
              cartElement={cartElement}
              currency={currency}
              removeFromCart={removeFromCart}
            />
          ))}
          <TableRow>
            <TableCell colSpan={2} className="uppercase font-bold">
              Total
            </TableCell>
            <TableCell className="font-bold text-right">
              {currency.format(cart.totalGross.toNumber())}
            </TableCell>
            <TableCell className="flex justify-end">
              {cart.isEmpty ? (
                <Button disabled color="failure" aria-label="Clear cart">
                  <HiXCircle />
                </Button>
              ) : (
                <Button
                  color="failure"
                  onClick={() => removeAllFromCart()}
                  aria-label="Clear cart"
                >
                  <HiXCircle />
                </Button>
              )}
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <CheckoutButtons
        cart={cart}
        checkoutProcessing={checkoutProcessing}
        handleCheckoutCart={checkoutCart}
      />
    </div>
  );
};

export default Cart;
