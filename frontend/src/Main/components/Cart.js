import React, { useEffect, useState, useRef } from "react";
import { HiXCircle } from "react-icons/hi";
import { Table, Tooltip } from "flowbite-react";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";
import "animate.css";
import MyButton from "./MyButton";

function Cart({ cart, removeFromCart, removeAllFromCart, checkoutCart }) {
  const currency = useConfig().currency;

  const [flash, setFlash] = useState(false);
  const flashCount = useRef(0);

  const triggerFlash = () => {
    setFlash(true);
    setTimeout(() => {
      setFlash(false);
    }, 500);
  };

  useEffect(() => {
    // not 100% sure why this is called twice
    if (flashCount.current < 2) {
      flashCount.current++;
      return;
    }
    triggerFlash();
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
        base: "px-4 py-2",
      },
    },
    body: {
      cell: {
        base: "px-4 py-2",
      },
    },
  };

  return (
    <div>
      <Table
        striped
        theme={compactTableTheme}
        className={`table-fixed ${flash ? "animate__animated animate__pulse" : ""}`}
      >
        <Table.Head>
          <Table.HeadCell className="w-5/12">Product</Table.HeadCell>
          <Table.HeadCell className="w-2/12 text-right">
            <Tooltip content="Quantity">Qnt</Tooltip>
          </Table.HeadCell>
          <Table.HeadCell className="w-2/12 text-right">
            Total Price
          </Table.HeadCell>
          <Table.HeadCell className="w-3/12">Remove</Table.HeadCell>
        </Table.Head>
        <Table.Body>
          {cart.map((cartElement) => (
            <Table.Row key={cartElement.id}>
              <Table.Cell className="whitespace-normal px-4 py-2">
                {cartElement.name}
                {
                  // iterate over cartElement.listItems and display them
                  cartElement.listItems.map((listItem) => (
                    <div key={listItem.id} className="text-xs text-gray-500">
                      {displayListItem(listItem)}
                    </div>
                  ))
                }
              </Table.Cell>
              <Table.Cell className="text-right">
                {cartElement.quantity}
              </Table.Cell>
              <Table.Cell className="text-right">
                {currency.format(cartElement.totalPrice)}
              </Table.Cell>
              <Table.Cell className="">
                <MyButton
                  color="failure"
                  onClick={() => removeFromCart(cartElement)}
                >
                  <HiXCircle />
                </MyButton>
              </Table.Cell>
            </Table.Row>
          ))}
          <Table.Row>
            <Table.Cell colSpan={2} className="uppercase font-bold">
              Total
            </Table.Cell>
            <Table.Cell className="font-bold text-right">
              {currency.format(
                cart.reduce((total, item) => total + item.totalPrice, 0),
              )}
            </Table.Cell>
            <Table.Cell>
              {cart.length ? (
                <MyButton color="failure" onClick={() => removeAllFromCart()}>
                  <HiXCircle />
                </MyButton>
              ) : (
                <MyButton disabled color="failure">
                  <HiXCircle />
                </MyButton>
              )}
            </Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>

      <MyButton
        {...(cart.length === 0 && { disabled: true })}
        color="success"
        className="w-full mt-2 uppercase"
        onClick={checkoutCart}
      >
        Checkout&nbsp;
        {cart.length > 0 &&
          currency.format(
            cart.reduce((total, item) => total + item.totalPrice, 0),
          )}
      </MyButton>
    </div>
  );
}

Cart.propTypes = {
  cart: PropTypes.array.isRequired,
  removeFromCart: PropTypes.func.isRequired,
  removeAllFromCart: PropTypes.func.isRequired,
  checkoutCart: PropTypes.func.isRequired,
};

export default Cart;
