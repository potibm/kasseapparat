import React, { useState } from "react";
import { Modal, ModalBody, ModalHeader, TextInput } from "flowbite-react";
import Decimal from "decimal.js";
import Button from "../../../components/Button";
import { useConfig } from "@core/config/hooks/useConfig";
import { parseAmount, suggestReceivedAmounts } from "../utils/change";

interface ChangeCalculatorModalProps {
  total: Decimal;
  paymentMethodName: string;
  onConfirm: () => void;
  onClose: () => void;
}

/**
 * Asks how much the customer handed over and shows the change before the payment is
 * booked. "Exact" is preselected, so an exact payment is booked with one more tap.
 * Mount it only while it is needed: it starts fresh every time.
 */
const ChangeCalculatorModal: React.FC<ChangeCalculatorModalProps> = ({
  total,
  paymentMethodName,
  onConfirm,
  onClose,
}) => {
  const { currency } = useConfig();
  const [received, setReceived] = useState<Decimal | null>(total);
  const [customInput, setCustomInput] = useState("");

  const format = (amount: Decimal) => currency.format(amount.toNumber());
  const suggestions = suggestReceivedAmounts(total);
  const enough = received !== null && received.greaterThanOrEqualTo(total);

  const choose = (amount: Decimal) => {
    setReceived(amount);
    setCustomInput("");
  };

  const handleCustomInput = (value: string) => {
    setCustomInput(value);
    setReceived(value.trim() === "" ? total : parseAmount(value));
  };

  const amountButtonClass = (amount: Decimal) =>
    `pos-pressable h-14 rounded-lg border-2 text-lg font-bold ${
      customInput === "" && received?.equals(amount)
        ? "border-cyan-700 bg-cyan-50 text-cyan-900 dark:border-cyan-500 dark:bg-cyan-950 dark:text-white"
        : "border-gray-300 bg-white text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
    }`;

  let result: React.ReactNode;
  if (received === null) {
    result = (
      <span className="text-gray-500 dark:text-gray-400">
        Enter a valid amount
      </span>
    );
  } else if (enough) {
    result = (
      <>
        <span className="text-gray-600 dark:text-gray-300">Change</span>
        <span
          data-testid="change-calculator-change"
          className="text-4xl font-bold text-green-700 dark:text-green-400"
        >
          {format(received.minus(total))}
        </span>
      </>
    );
  } else {
    result = (
      <>
        <span className="text-gray-600 dark:text-gray-300">Missing</span>
        <span
          data-testid="change-calculator-missing"
          className="text-4xl font-bold text-red-700 dark:text-red-400"
        >
          {format(total.minus(received))}
        </span>
      </>
    );
  }

  return (
    <Modal show onClose={onClose} dismissible size="lg">
      <ModalHeader>{paymentMethodName}</ModalHeader>
      <ModalBody>
        <div className="flex flex-col gap-4 dark:text-white">
          <div className="flex items-baseline justify-between">
            <span className="text-sm font-semibold tracking-wider text-gray-600 uppercase dark:text-gray-300">
              Total
            </span>
            <span
              data-testid="change-calculator-total"
              className="text-3xl font-bold"
            >
              {format(total)}
            </span>
          </div>

          <div>
            <div className="mb-2 text-sm font-semibold tracking-wider text-gray-600 uppercase dark:text-gray-300">
              Received
            </div>
            <div className="grid grid-cols-3 gap-2">
              <button
                type="button"
                className={amountButtonClass(total)}
                onClick={() => choose(total)}
              >
                Exact
              </button>
              {suggestions.map((amount) => (
                <button
                  key={amount.toString()}
                  type="button"
                  className={amountButtonClass(amount)}
                  onClick={() => choose(amount)}
                >
                  {format(amount)}
                </button>
              ))}
            </div>
            <TextInput
              className="mt-2"
              sizing="lg"
              inputMode="decimal"
              placeholder="Other amount"
              aria-label="Other amount received"
              value={customInput}
              onChange={(event) => handleCustomInput(event.target.value)}
            />
          </div>

          <div
            className="flex min-h-16 items-baseline justify-between"
            aria-live="polite"
          >
            {result}
          </div>

          <div className="grid grid-cols-2 gap-2">
            <Button
              color="light"
              className="w-full"
              onClick={onClose}
              data-testid="change-calculator-cancel"
            >
              Cancel
            </Button>
            <Button
              className="w-full"
              disabled={!enough}
              onClick={onConfirm}
              data-testid="change-calculator-confirm"
            >
              Book payment
            </Button>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};

export default ChangeCalculatorModal;
