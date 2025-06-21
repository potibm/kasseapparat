import React from "react";
import { useConfig } from "../../../provider/ConfigProvider";
import MyButton from "../MyButton";
import Decimal from "decimal.js";
import { Spinner } from "flowbite-react";
import { getCurrentReaderId } from "../../../helper/ReaderCookie";
import PropTypes from "prop-types";

const CheckoutButtons = ({ cart, checkoutProcessing, handleCheckoutCart }) => {
  const { currency, paymentMethods } = useConfig();
  const sumUpReaderId = getCurrentReaderId();

  const cartValue = cart.reduce(
    (total, item) => total.add(item.totalGrossPrice),
    new Decimal(0),
  );

  const paymentMethodIsActive = (paymentMethodCode, cartValue) => {
    if (paymentMethodCode === "SUMUP") {
      return sumUpReaderId !== undefined && cartValue.greaterThan(0);
    }

    return true;
  };

  const getPaymentMethodData = (paymentMethodCode) => {
    if (paymentMethodCode === "SUMUP") {
      return {
        sumupReaderId: sumUpReaderId,
      };
    }

    return {};
  };

  return (
    <>
      {paymentMethods.map((paymentMethod) => (
        <MyButton
          key={paymentMethod.code}
          {...((cart.length === 0 ||
            checkoutProcessing ||
            !paymentMethodIsActive(paymentMethod.code, cartValue)) && {
            disabled: true,
          })}
          className="w-full mt-2 uppercase"
          onClick={() =>
            handleCheckoutCart(
              paymentMethod.code,
              getPaymentMethodData(paymentMethod.code),
            )
          }
        >
          {paymentMethod.name}&nbsp;
          {cart.length > 0 && currency.format(cartValue)}
          {checkoutProcessing === paymentMethod.code && (
            <Spinner color="gray" className="ml-3" />
          )}
        </MyButton>
      ))}
    </>
  );
};

CheckoutButtons.propTypes = {
  cart: PropTypes.array.isRequired,
  checkoutProcessing: PropTypes.string.isRequired,
  handleCheckoutCart: PropTypes.func.isRequired,
};

export default CheckoutButtons;
