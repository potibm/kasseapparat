import { useConfig } from "../../../../../../core/config/providers/ConfigProvider";
import Button from "../../../../components/Button";
import { Spinner } from "flowbite-react";
import { getCurrentReaderId } from "../../../../../../core/localstorage/helper/reader";
import { Cart } from "../../services/Cart";
import Decimal from "decimal.js";
import { PaymentMethodData } from "../../types/cart.types";

interface CheckoutButtonsProps {
  cart: Cart;
  checkoutProcessing: string | null;
  handleCheckoutCart: (code: string, data: PaymentMethodData) => void;
}

const CheckoutButtons = ({
  cart,
  checkoutProcessing,
  handleCheckoutCart,
}: CheckoutButtonsProps) => {
  const { currency, paymentMethods } = useConfig();
  const sumUpReaderId = getCurrentReaderId();

  const cartValue = cart.totalGross;

  const paymentMethodIsActive = (
    paymentMethodCode: string,
    cartValue: Decimal,
  ) => {
    if (paymentMethodCode === "SUMUP") {
      return sumUpReaderId !== undefined && cartValue.greaterThan(0);
    }

    return true;
  };

  const getPaymentMethodData = (
    paymentMethodCode: string,
  ): PaymentMethodData => {
    if (paymentMethodCode === "SUMUP") {
      return {
        type: "sumup",
        sumupReaderId: String(sumUpReaderId),
      };
    }

    return {
      type: "empty",
    };
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
          {!cart.isEmpty && currency.format(cartValue.toNumber())}
          {checkoutProcessing === paymentMethod.code && (
            <Spinner color="gray" className="ml-3" />
          )}
        </Button>
      ))}
    </>
  );
};

export default CheckoutButtons;
