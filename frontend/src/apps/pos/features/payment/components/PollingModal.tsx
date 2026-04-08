import React, { useEffect, useState, useMemo } from "react";
import { Modal, ModalBody, ModalFooter, Spinner } from "flowbite-react";
import { HiClock } from "react-icons/hi";
import { getCurrentReaderId } from "@core/localstorage/helper/local-storage-reader";
import Button from "@pos/components/Button";
import { usePaymentWebSocket } from "../hooks/usePaymentWebSocket";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { ToastSeverity } from "@pos/features/ui/toast/types/toast.types";
import { useConfig } from "@core/config/hooks/useConfig";
import { Purchase as PurchaseType } from "../../../utils/api.schemas";

const log = createLogger("Payment");

/**
 * Props definition for the PollingModal component.
 */
interface PollingModalProps {
  purchase: PurchaseType;
  onComplete: (success: boolean) => void;
}

/**
 * PollingModal handles the visual representation of the payment process.
 * It subscribes to the usePaymentWebSocket hook for real-time updates.
 */
export const PollingModal: React.FC<PollingModalProps> = ({
  purchase,
  onComplete,
}) => {
  // Connect to the WebSocket logic via our custom hook
  const purchaseId = useMemo(() => purchase.id, [purchase.id]);
  const { status, error, lastMessageAt, cancel, isConnected } =
    usePaymentWebSocket(purchaseId);

  const { showToast } = useToast();
  const { currency, paymentMethods } = useConfig();

  const paymentMethodName =
    paymentMethods.find((method) => method.code === purchase.paymentMethod)
      ?.name || purchase.paymentMethod;

  const [isAborting, setIsAborting] = useState(false);
  const [now, setNow] = useState(() => Date.now());
  const isTerminalError = [
    "failed",
    "cancelled",
    "connection_lost",
    "timeout",
  ].includes(status);

  // Local timer to refresh the "seconds ago" display every second
  useEffect(() => {
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(timer);
  }, []);

  const ageInSeconds = Math.max(0, Math.round((now - lastMessageAt) / 1000));

  useEffect(() => {
    if (status === "confirmed") {
      log.info("Purchase confirmed", purchase.id);
      showToast({
        severity: "success",
        message: `Purchase of ${currency.format(purchase.totalGrossPrice.toNumber())} using ${paymentMethodName} was successful!`,
      });
      onComplete(true);
    }

    if (isTerminalError) {
      log.error("Terminal process ended with status", purchase.id, status);

      const toastMessage = {
        severity: "error" as ToastSeverity,
        message: "Payment failed.",
        autoClose: status === "cancelled",
        blocking: status !== "cancelled",
      };
      if (status === "cancelled") {
        toastMessage.message = `Payment using ${paymentMethodName} was cancelled by the user.`;
      } else if (status === "connection_lost") {
        toastMessage.message = `Connection to the terminal lost for payment method ${paymentMethodName}.`;
      } else if (status === "timeout") {
        toastMessage.message = `Timeout at the terminal for payment method ${paymentMethodName}.`;
      }

      showToast(toastMessage);

      onComplete(false);
    }
  }, [
    status,
    isTerminalError,
    purchase.id,
    purchase.totalGrossPrice,
    onComplete,
    showToast,
    paymentMethodName,
    currency,
  ]);

  /**
   * Triggers the cancellation process via WebSocket.
   */
  const handleAbort = () => {
    log.info("User initiated purchase cancellation", { purchaseId });
    setIsAborting(true);
    cancel(getCurrentReaderId());
  };

  return (
    <Modal show={true} size="md" popup dismissible={false}>
      <ModalBody className="m-5">
        <div className="text-center">
          <h3 className="text-lg font-semibold mb-8 dark:text-gray-200">
            Purchase Status
          </h3>
          <div className="text-center space-y-4">
            <div className="text-lg font-medium flex items-center justify-center gap-2 text-yellow-600">
              <HiClock className="inline mr-1 h-5 w-5" />
              <span>Waiting for Terminal... ({ageInSeconds} sec ago)</span>
            </div>

            <Spinner color="warning" size="xl" className="mb-3" />

            {error && <p className="text-red-600 font-medium">{error}</p>}
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <div className="w-full flex justify-center">
          {status === "pending" && (
            <Button disabled={isAborting || !isConnected} onClick={handleAbort}>
              {isAborting ? "Cancelling..." : "Abort Purchase"}
              {isAborting && (
                <Spinner color="gray" size="sm" className="ml-2" />
              )}
            </Button>
          )}
          {isTerminalError && (
            <Button onClick={() => onComplete(false)}>Close</Button>
          )}
        </div>
      </ModalFooter>
    </Modal>
  );
};

export default PollingModal;
