import React from "react";
import { useConfig } from "../../../../../../core/config/providers/config-provider";
import Button from "../../../../components/Button";
import { Spinner } from "flowbite-react";
import { getCurrentReaderId } from "../../../../../../core/localstorage/helper/reader";
import PropTypes from "prop-types";
import { Cart } from "../../services/Cart";

interface CheckoutButtonsProps {
  cart: Cart;
  checkoutProcessing: string | boolean;
  handleCheckoutCart: (code: string, data: any) => void;
}

const CheckoutButtons = ({
  cart,
  checkoutProcessing,
  handleCheckoutCart,
}: CheckoutButtonsProps) => {
  const { currency, paymentMethods } = useConfig();
  const sumUpReaderId = getCurrentReaderId();

  const cartValue = cart.totalGross;

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
        <Button
          key={paymentMethod.code}
          {...((cart.isEmpty ||
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
          {!cart.isEmpty && currency.format(cartValue)}
          {checkoutProcessing === paymentMethod.code && (
            <Spinner color="gray" className="ml-3" />
          )}
        </Button>
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
