import React, { useState, useEffect, useRef } from "react";
import { HiReceiptRefund, HiRefresh } from "react-icons/hi";
import { Spinner, Table, TableRow, TableCell, TableBody } from "flowbite-react";
import { useConfig } from "@core/config/hooks/useConfig";
import { Purchase } from "../../../api/schemas";
import Button from "../../../components/Button";
import { RefundModal } from "./_internal/RefundModal";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";

const log = createLogger("Purchase");

interface PurchaseHistoryProps {
  history: Purchase[] | null;
  loading: boolean;
  removeFromPurchaseHistory: (purchase: Purchase) => Promise<void>;
  resumePolling: (purchase: Purchase) => void;
  cartEmpty: boolean;
}

const compactTableTheme = {
  head: {
    cell: {
      base: "px-2 py-1",
    },
  },
  body: {
    cell: {
      base: "px-2 py-1",
    },
  },
};

const PurchaseHistory: React.FC<PurchaseHistoryProps> = ({
  history,
  loading,
  removeFromPurchaseHistory,
  resumePolling,
  cartEmpty,
}) => {
  const { currency, dateLocale, dateOptions } = useConfig();
  const { showToast } = useToast();

  // State
  const [modalState, setModalState] = useState<{
    show: boolean;
    purchase: Purchase | null;
  }>({
    show: false,
    purchase: null,
  });
  const [processingRefund, setProcessingRefund] = useState(false);
  const [flash, setFlash] = useState(false);

  const lastTopIdRef = useRef<number | string | null>(null);

  useEffect(() => {
    if (!history || history.length === 0) return;

    const currentTopId = history[0].id;
    if (lastTopIdRef.current === null) {
      lastTopIdRef.current = currentTopId;
    } else if (lastTopIdRef.current !== currentTopId) {
      let timerId: ReturnType<typeof setTimeout>;

      const rafId = requestAnimationFrame(() => {
        setFlash(true);
        timerId = setTimeout(() => setFlash(false), 500);
      });

      lastTopIdRef.current = currentTopId;
      return () => {
        if (rafId) cancelAnimationFrame(rafId);
        if (timerId) clearTimeout(timerId);
      };
    }
  }, [history]);

  const handleRefund = async (purchase: Purchase) => {
    setProcessingRefund(true);
    try {
      await removeFromPurchaseHistory(purchase);
      log.info("Purchase refunded successfully", { purchaseId: purchase.id });
      setModalState({ show: false, purchase: null });
    } finally {
      setProcessingRefund(false);
    }
  };

  const formatDate = (date: string | Date) =>
    new Date(date).toLocaleString(dateLocale, dateOptions);

  // The sidebar only has room for the time; the full date stays in the labels.
  const formatTime = (date: string | Date) =>
    new Date(date).toLocaleTimeString(dateLocale, {
      hour: "2-digit",
      minute: "2-digit",
    });

  return (
    <div className="mt-3 shrink-0">
      <RefundModal
        show={modalState.show}
        purchase={modalState.purchase}
        processing={processingRefund}
        onClose={() => setModalState({ show: false, purchase: null })}
        onConfirm={handleRefund}
      />

      <div className="px-2 pb-1 text-left text-xs font-semibold uppercase tracking-wide text-gray-600 dark:text-gray-400">
        Recent sales
      </div>
      <Table
        striped
        data-testid="purchase-history-table"
        theme={compactTableTheme}
        className={`table-fixed dark:text-gray-200 ${flash ? "animate__animated animate__pulse" : ""}`}
      >
        <TableBody>
          {loading && (!history || history.length === 0) && (
            <TableRow>
              <TableCell colSpan={3}>
                Loading... <Spinner size="sm" />
              </TableCell>
            </TableRow>
          )}
          {!loading && history?.length === 0 && (
            <TableRow>
              <TableCell colSpan={3}>No purchases yet.</TableCell>
            </TableRow>
          )}
          {history &&
            history.length > 0 &&
            history.slice(0, 3).map((purchase) => (
              <TableRow key={purchase.id}>
                <TableCell
                  className="whitespace-nowrap text-left"
                  title={formatDate(purchase.createdAt)}
                >
                  {formatTime(purchase.createdAt)}
                </TableCell>
                <TableCell className="text-right">
                  {currency.format(purchase.totalGrossPrice.toNumber())}
                </TableCell>
                <TableCell className="flex justify-end">
                  {purchase.status === "confirmed" && (
                    <Button
                      color="failure"
                      size="sm"
                      aria-label={`Refund purchase from ${formatDate(purchase.createdAt)}`}
                      onClick={() => setModalState({ show: true, purchase })}
                      data-testid={`refund-purchase-${purchase.id}`}
                    >
                      <HiReceiptRefund />
                    </Button>
                  )}
                  {purchase.status === "pending" && (
                    <Button
                      color="warning"
                      size="sm"
                      aria-label={`Resume purchase`}
                      onClick={() => {
                        if (cartEmpty) {
                          resumePolling(purchase);
                        } else {
                          showToast({
                            severity: "warning",
                            message:
                              "Please complete or clear the current cart before resuming this purchase.",
                          });
                        }
                      }}
                    >
                      <HiRefresh className="animate-spin" />
                    </Button>
                  )}
                </TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </div>
  );
};

export default PurchaseHistory;
