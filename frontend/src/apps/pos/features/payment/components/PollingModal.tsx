import React, { useEffect, useState, useMemo } from "react";
import { Modal, ModalBody, Spinner } from "flowbite-react";
import { HiCreditCard, HiExclamation } from "react-icons/hi";
import { getCurrentReaderId } from "@core/localstorage/helper/local-storage-reader";
import Button from "@pos/components/Button";
import { usePaymentWebSocket } from "../hooks/usePaymentWebSocket";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { ToastSeverity } from "@pos/features/ui/toast/types/toast.types";
import { useConfig } from "@core/config/hooks/useConfig";
import { Purchase as PurchaseType } from "../../../api/schemas";

const log = createLogger("Payment");

/** Without news from the terminal for this long, the dialog says so. */
export const STALE_AFTER_SECONDS = 15;

/** A connection that is still not up after this long is reported. */
export const CONNECTING_GRACE_SECONDS = 3;

const formatElapsed = (seconds: number) =>
  `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, "0")}`;

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
  const [startedAt] = useState(() => Date.now());
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
  const elapsedSeconds = Math.max(0, Math.floor((now - startedAt) / 1000));

  // Only problems are spelled out; a normal payment shows amount and instruction.
  let problem: string | null = null;
  if (error) {
    problem = error;
  } else if (!isConnected && elapsedSeconds >= CONNECTING_GRACE_SECONDS) {
    problem = "No connection to the payment service. Retrying…";
  } else if (ageInSeconds >= STALE_AFTER_SECONDS) {
    problem = `No news from the terminal for ${ageInSeconds} s.`;
  }

  useEffect(() => {
    if (status === "confirmed") {
      log.info("Purchase confirmed", purchase.id);
      showToast({
        severity: "success",
        message: `Payment of ${currency.format(purchase.totalGrossPrice.toNumber())} successful via ${paymentMethodName}!`,
      });
      onComplete(true);
    }

    if (isTerminalError) {
      log.error("Terminal process ended with status", purchase.id, status);

      const toastMessage = {
        severity: "error" as ToastSeverity,
        message: `Payment with ${paymentMethodName} failed.`,
        autoClose: status === "cancelled",
        blocking: status !== "cancelled",
      };
      if (status === "cancelled") {
        toastMessage.message = `Payment via ${paymentMethodName} was cancelled by the user.`;
      } else if (status === "connection_lost") {
        toastMessage.message = `Connection lost: ${paymentMethodName} terminal is unreachable.`;
      } else if (status === "timeout") {
        toastMessage.message = `Timeout: No response from ${paymentMethodName} terminal.`;
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
      <ModalBody className="p-6">
        <div className="flex flex-col gap-5 dark:text-white">
          <div className="flex items-center justify-between">
            <span className="text-sm font-semibold tracking-wider text-gray-600 uppercase dark:text-gray-300">
              {paymentMethodName}
            </span>
            <span
              data-testid="polling-elapsed"
              className="text-sm font-semibold text-gray-700 tabular-nums dark:text-gray-200"
            >
              {formatElapsed(elapsedSeconds)}
            </span>
          </div>

          <div className="text-center">
            <div
              data-testid="polling-amount"
              className="text-5xl font-bold tracking-tight"
            >
              {currency.format(purchase.totalGrossPrice.toNumber())}
            </div>
            <div className="mt-2 text-base text-gray-700 dark:text-gray-300">
              Tap or insert the card on the reader
            </div>
          </div>

          <div className="flex justify-center py-2">
            <div className="relative h-24 w-24">
              {/* expanding ring: clearly visible "waiting" signal */}
              <span
                aria-hidden="true"
                data-testid="polling-pulse"
                className="absolute inset-0 rounded-full bg-cyan-400/70 motion-safe:animate-ping dark:bg-cyan-500/60"
              />
              <div className="relative flex h-24 w-24 items-center justify-center rounded-full bg-cyan-100 ring-4 ring-cyan-300 dark:bg-cyan-900 dark:ring-cyan-700">
                <HiCreditCard
                  aria-hidden="true"
                  className="h-12 w-12 text-cyan-700 dark:text-cyan-300"
                />
              </div>
            </div>
          </div>

          {problem && (
            <div
              role="alert"
              data-testid="polling-problem"
              className="flex items-start gap-2 rounded-lg bg-amber-50 p-3 text-sm font-medium text-amber-900 dark:bg-amber-950 dark:text-amber-100"
            >
              <HiExclamation aria-hidden="true" className="h-5 w-5 shrink-0" />
              <span>{problem}</span>
            </div>
          )}

          {status === "pending" && (
            <Button
              color="red"
              outline
              size="xl"
              className="w-full"
              disabled={isAborting || !isConnected}
              onClick={handleAbort}
            >
              {isAborting ? "Cancelling…" : "Cancel payment"}
              {isAborting && (
                <Spinner color="gray" size="sm" className="ml-2" />
              )}
            </Button>
          )}
          {isTerminalError && (
            <Button
              size="xl"
              className="w-full"
              onClick={() => onComplete(false)}
            >
              Close
            </Button>
          )}
        </div>
      </ModalBody>
    </Modal>
  );
};

export default PollingModal;
