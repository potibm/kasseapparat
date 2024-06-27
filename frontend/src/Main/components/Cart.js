import React, { useEffect, useState, useRef } from "react";
import { HiXCircle } from "react-icons/hi";
import { Button, Table } from "flowbite-react";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";
import "animate.css";

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

  return (
    <div>
      <Table
        striped
        className={`table-fixed ${flash ? "animate__animated animate__pulse" : ""}`}
      >
        <Table.Head>
          <Table.HeadCell>Product</Table.HeadCell>
          <Table.HeadCell className="text-right">Quantity</Table.HeadCell>
          <Table.HeadCell className="text-right">Total Price</Table.HeadCell>
          <Table.HeadCell>Remove</Table.HeadCell>
        </Table.Head>
        <Table.Body>
          {cart.map((cartElement) => (
            <Table.Row key={cartElement.id}>
              <Table.Cell className="whitespace-nowrap">
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
              <Table.Cell className="text-right align-top">
                {cartElement.quantity}
              </Table.Cell>
              <Table.Cell className="text-right align-top">
                {currency.format(cartElement.totalPrice)}
              </Table.Cell>
              <Table.Cell className="align-top">
                <Button
                  color="failure"
                  onClick={() => removeFromCart(cartElement)}
                >
                  <HiXCircle />
                </Button>
              </Table.Cell>
            </Table.Row>
          ))}
          <Table.Row>
            <Table.Cell className="uppercase font-bold">Total</Table.Cell>
            <Table.Cell></Table.Cell>
            <Table.Cell className="font-bold text-right">
              {currency.format(
                cart.reduce((total, item) => total + item.totalPrice, 0),
              )}
            </Table.Cell>
            <Table.Cell>
              {cart.length ? (
                <Button color="failure" onClick={() => removeAllFromCart()}>
                  <HiXCircle />
                </Button>
              ) : (
                <Button disabled color="failure">
                  <HiXCircle />
                </Button>
              )}
            </Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>

      <Button
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
      </Button>
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
