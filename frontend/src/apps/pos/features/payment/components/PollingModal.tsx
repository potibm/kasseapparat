import React, { useEffect, useState, useMemo } from "react";
import { Modal, ModalBody, ModalFooter, Spinner } from "flowbite-react";
import { HiClock, HiCheckCircle, HiXCircle, HiBan } from "react-icons/hi";
import { getCurrentReaderId } from "@core/localstorage/helper/reader";
import Button from "@pos/components/Button";
import { usePaymentWebSocket } from "../hooks/usePaymentWebSocket";

interface PurchaseData {
  id: string;
  status: string;
}

interface PaymentConfirmation {
  id: string;
}

/**
 * Props definition for the PollingModal component.
 */
interface PollingModalProps {
  purchase: PurchaseData;
  onClose: () => void;
  onConfirmed: (data: PaymentConfirmation) => void;
  onComplete: (success: boolean) => void;
}

/**
 * PollingModal handles the visual representation of the payment process.
 * It subscribes to the usePaymentWebSocket hook for real-time updates.
 */
export const PollingModal: React.FC<PollingModalProps> = ({
  purchase,
  onClose,
  onConfirmed,
  onComplete,
}) => {
  // Connect to the WebSocket logic via our custom hook
  const purchaseId = useMemo(() => purchase.id, [purchase.id]);
  const { status, error, lastMessageAt, cancel, isConnected } =
    usePaymentWebSocket(purchaseId);

  const [isAborting, setIsAborting] = useState(false);
  const [now, setNow] = useState(() => Date.now());

  // Local timer to refresh the "seconds ago" display every second
  useEffect(() => {
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(timer);
  }, []);

  const ageInSeconds = Math.max(0, Math.round((now - lastMessageAt) / 1000));

  /**
   * Memoized UI configuration based on the current payment status.
   * This centralizes colors, icons, and labels.
   */
  const statusConfig = useMemo(() => {
    switch (status) {
      case "confirmed":
        return {
          textColor: "text-green-600",
          icon: <HiCheckCircle className="inline mr-1 h-8 w-8" />,
          label: "Payment Successful",
          spinnerColor: null,
        };
      case "failed":
      case "connection_lost":
      case "timeout":
        return {
          textColor: "text-red-600",
          icon: <HiXCircle className="inline mr-1 h-8 w-8" />,
          label: "Payment Failed",
          spinnerColor: null,
        };
      case "cancelled":
        return {
          textColor: "text-gray-500",
          icon: <HiBan className="inline mr-1 h-8 w-8" />,
          label: "Payment Cancelled",
          spinnerColor: null,
        };
      default: // status === "pending"
        return {
          textColor: "text-yellow-600",
          icon: <HiClock className="inline mr-1 h-5 w-5" />,
          label: `Waiting for Terminal...`,
          spinnerColor: "warning" as const,
        };
    }
  }, [status]);

  /**
   * Effect to trigger external callbacks when the terminal process finishes.
   */
  useEffect(() => {
    if (status === "confirmed") {
      onComplete(true);
      // Brief delay to allow the user to see the success state before closing
      const timer = setTimeout(() => onConfirmed({ id: purchase.id }), 2500);
      return () => clearTimeout(timer);
    }

    const terminalStates: string[] = [
      "failed",
      "cancelled",
      "connection_lost",
      "timeout",
    ];
    if (terminalStates.includes(status)) {
      onComplete(false);
    }
  }, [status, purchase.id, onComplete, onConfirmed]);

  /**
   * Triggers the cancellation process via WebSocket.
   */
  const handleAbort = () => {
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
            <div
              className={`text-lg font-medium flex items-center justify-center gap-2 ${statusConfig.textColor}`}
            >
              {statusConfig.icon}
              <span>
                {statusConfig.label}
                {status === "pending" && <> ({ageInSeconds} sec ago)</>}
              </span>
            </div>
            {status === "pending" && statusConfig.spinnerColor && (
              <Spinner
                color={statusConfig.spinnerColor}
                size="xl"
                className="mb-3"
              />
            )}
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
          {(status === "failed" ||
            status === "cancelled" ||
            status === "connection_lost" ||
            status === "timeout") && <Button onClick={onClose}>Close</Button>}
        </div>
      </ModalFooter>
    </Modal>
  );
};

export default PollingModal;
