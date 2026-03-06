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
import PropTypes from "prop-types";
import { useConfig } from "../../../../../core/config/providers/config-provider";
import "animate.css";
import Button from "../../../components/Button";
import Decimal from "decimal.js";
import CheckoutButtons from "./_internal/CheckoutButtons";
import { getCartTotalQuantity } from "./services/cart.logic";

const Cart = ({ cart, removeFromCart, removeAllFromCart, checkoutCart }) => {
  const { currency } = useConfig();

  const [flash, setFlash] = useState(false);
  const [checkoutProcessing, setCheckoutProcessing] = useState(null);

  const prevCartTotalQuantity = useRef(getCartTotalQuantity(cart));
  const isFirstRender = useRef(true);

  const triggerFlash = () => {
    requestAnimationFrame(() => {
      setFlash(true);
      setTimeout(() => {
        setFlash(false);
      }, 500);
    });
  };

  const handleCheckoutCart = async (paymentMethodCode, paymentMethodData) => {
    if (checkoutProcessing) {
      return;
    }
    setCheckoutProcessing(paymentMethodCode);
    checkoutCart(paymentMethodCode, paymentMethodData).then(() => {
      setCheckoutProcessing(null);
    });
  };

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }

    if (getCartTotalQuantity(cart) !== prevCartTotalQuantity.current) {
      triggerFlash();
    }

    prevCartTotalQuantity.current = getCartTotalQuantity(cart);
  }, [cart]);

  const displayListItem = (listItem) => {
    if (listItem.code !== null) {
      return listItem.code;
    } else {
      return listItem.name;
    }
  };

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
          {cart.map((cartElement) => (
            <TableRow key={cartElement.id}>
              <TableCell className="whitespace-normal px-4 py-2">
                {cartElement.name}
                {
                  // iterate over cartElement.listItems and display them
                  cartElement.listItems.map((listItem) => (
                    <div key={listItem.id} className="text-xs text-gray-500">
                      {displayListItem(listItem)}
                    </div>
                  ))
                }
              </TableCell>
              <TableCell className="text-right">
                {cartElement.quantity}
              </TableCell>
              <TableCell className="text-right">
                {currency.format(cartElement.totalGrossPrice)}
              </TableCell>
              <TableCell className="flex justify-end">
                <Button
                  color="failure"
                  onClick={() => removeFromCart(cartElement)}
                >
                  <HiXCircle />
                </Button>
              </TableCell>
            </TableRow>
          ))}
          <TableRow>
            <TableCell colSpan={2} className="uppercase font-bold">
              Total
            </TableCell>
            <TableCell className="font-bold text-right">
              {currency.format(
                cart.reduce(
                  (total, item) => total.add(item.totalGrossPrice),
                  new Decimal(0),
                ),
              )}
            </TableCell>
            <TableCell className="flex justify-end">
              {cart.length ? (
                <Button color="failure" onClick={() => removeAllFromCart()}>
                  <HiXCircle />
                </Button>
              ) : (
                <Button disabled color="failure">
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
        handleCheckoutCart={handleCheckoutCart}
      />
    </div>
  );
};

Cart.propTypes = {
  cart: PropTypes.array.isRequired,
  removeFromCart: PropTypes.func.isRequired,
  removeAllFromCart: PropTypes.func.isRequired,
  checkoutCart: PropTypes.func.isRequired,
};

export default Cart;
