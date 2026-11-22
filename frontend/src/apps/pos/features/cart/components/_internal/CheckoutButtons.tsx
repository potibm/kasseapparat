import { useState } from "react";
import { useConfig } from "@core/config/hooks/useConfig";
import Button from "../../../../components/Button";
import { Spinner } from "flowbite-react";
import { getCurrentReaderId } from "@core/localstorage/helper/local-storage-reader";
import { Cart } from "../../services/Cart";
import Decimal from "decimal.js";
import { PaymentMethodData } from "../../types/cart.types";
import ChangeCalculatorModal from "@pos/features/payment/components/ChangeCalculatorModal";

type PaymentMethod = ReturnType<typeof useConfig>["paymentMethods"][number];

interface CheckoutButtonsProps {
  cart: Cart;
  checkoutProcessing: string | null;
  handleCheckoutCart: (code: string, data: PaymentMethodData) => void;
}

const CheckoutButtons: React.FC<CheckoutButtonsProps> = ({
  cart,
  checkoutProcessing,
  handleCheckoutCart,
}) => {
  const { currency, paymentMethods } = useConfig();
  const sumUpReaderId = getCurrentReaderId();
  const [calculatorFor, setCalculatorFor] = useState<PaymentMethod | null>(
    null,
  );

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

  const checkoutWith = (paymentMethod: PaymentMethod) =>
    handleCheckoutCart(
      paymentMethod.code,
      getPaymentMethodData(paymentMethod.code),
    );

  const handleClick = (paymentMethod: PaymentMethod) => {
    if (paymentMethod.changeCalculator) {
      setCalculatorFor(paymentMethod);
    } else {
      checkoutWith(paymentMethod);
    }
  };

  return (
    <div className="mt-2 grid grid-cols-2 gap-2">
      {paymentMethods.map((paymentMethod) => (
        <Button
          key={paymentMethod.code}
          data-testid={"checkout-button-" + paymentMethod.code}
          aria-label={
            cart.isEmpty
              ? `Pay with ${paymentMethod.name}`
              : `Pay ${currency.format(cartValue.toNumber())} with ${paymentMethod.name}`
          }
          {...((cart.isEmpty ||
            checkoutProcessing ||
            !paymentMethodIsActive(paymentMethod.code, cartValue)) && {
            disabled: true,
          })}
          className="w-full uppercase odd:last:col-span-2"
          onClick={() => handleClick(paymentMethod)}
        >
          <span className="flex flex-col items-center leading-tight">
            <span>{paymentMethod.name}</span>
            {!cart.isEmpty && (
              <span className="text-sm font-normal normal-case">
                {currency.format(cartValue.toNumber())}
              </span>
            )}
          </span>
          {checkoutProcessing === paymentMethod.code && (
            <Spinner color="gray" className="ml-3" />
          )}
        </Button>
      ))}
      {calculatorFor && (
        <ChangeCalculatorModal
          total={cartValue}
          paymentMethodName={calculatorFor.name}
          onClose={() => setCalculatorFor(null)}
          onConfirm={() => {
            setCalculatorFor(null);
            checkoutWith(calculatorFor);
          }}
        />
      )}
    </div>
  );
};

export default CheckoutButtons;
