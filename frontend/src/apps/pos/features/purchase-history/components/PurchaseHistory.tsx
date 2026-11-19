import React, { useState, useEffect, useRef } from "react";
import { HiLockClosed, HiReceiptRefund, HiRefresh } from "react-icons/hi";
import { Spinner } from "flowbite-react";
import { useConfig } from "@core/config/hooks/useConfig";
import { Purchase } from "../../../api/schemas";
import { RefundModal } from "./_internal/RefundModal";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { REFUND_WINDOW_MS, summarizeItems } from "../utils/purchase.utils";

const log = createLogger("Purchase");

/** How often the refund windows are re-evaluated while the POS is open. */
const REFUND_WINDOW_TICK_MS = 30 * 1000;

const MAX_VISIBLE_PURCHASES = 3;

interface PurchaseHistoryProps {
  history: Purchase[] | null;
  loading: boolean;
  removeFromPurchaseHistory: (purchase: Purchase) => Promise<void>;
  resumePolling: (purchase: Purchase) => void;
  cartEmpty: boolean;
  /** Admins are exempt from the refund time limit. */
  canRefundAnytime?: boolean;
}

const actionButton =
  "flex min-h-10 shrink-0 items-center gap-1 rounded-lg border px-3 text-sm font-semibold";

const PurchaseHistory: React.FC<PurchaseHistoryProps> = ({
  history,
  loading,
  removeFromPurchaseHistory,
  resumePolling,
  cartEmpty,
  canRefundAnytime = false,
}) => {
  const { currency, dateLocale, dateOptions, paymentMethods } = useConfig();
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
  const [now, setNow] = useState(() => Date.now());

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

  // Re-evaluate the refund windows while the POS stays open.
  useEffect(() => {
    if (canRefundAnytime) return;
    const intervalId = setInterval(
      () => setNow(Date.now()),
      REFUND_WINDOW_TICK_MS,
    );
    return () => clearInterval(intervalId);
  }, [canRefundAnytime]);

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

  const handleResume = (purchase: Purchase) => {
    if (cartEmpty) {
      resumePolling(purchase);
    } else {
      showToast({
        severity: "warning",
        message:
          "Please complete or clear the current cart before resuming this purchase.",
      });
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

  const paymentMethodName = (code: string) =>
    paymentMethods?.find((method) => method.code === code)?.name ?? code;

  const isRefundable = (purchase: Purchase) =>
    canRefundAnytime ||
    now - new Date(purchase.createdAt).getTime() <= REFUND_WINDOW_MS;

  const visibleHistory = history?.slice(0, MAX_VISIBLE_PURCHASES) ?? [];

  return (
    <div className="mt-3 shrink-0 text-left">
      <RefundModal
        show={modalState.show}
        purchase={modalState.purchase}
        processing={processingRefund}
        onClose={() => setModalState({ show: false, purchase: null })}
        onConfirm={handleRefund}
      />

      <div className="px-2 pb-1 text-xs font-semibold uppercase tracking-wide text-gray-600 dark:text-gray-400">
        Recent sales
      </div>

      <div
        data-testid="purchase-history-table"
        className={`overflow-hidden rounded-lg bg-white dark:bg-gray-800 dark:text-gray-200 ${flash ? "animate__animated animate__pulse" : ""}`}
      >
        {loading && visibleHistory.length === 0 && (
          <p className="flex items-center gap-2 px-3 py-2 text-sm">
            Loading... <Spinner size="sm" />
          </p>
        )}
        {!loading && history?.length === 0 && (
          <p className="px-3 py-2 text-sm">No purchases yet.</p>
        )}

        {visibleHistory.length > 0 && (
          <ul className="divide-y divide-gray-200 dark:divide-gray-700">
            {visibleHistory.map((purchase) => {
              const isPending = purchase.status === "pending";
              const summary = summarizeItems(purchase);

              return (
                <li
                  key={purchase.id}
                  data-testid={`purchase-${purchase.id}`}
                  className={`flex items-center gap-2 px-3 py-1.5 ${isPending ? "bg-amber-50 dark:bg-amber-900/30" : ""}`}
                >
                  <div className="min-w-0 flex-1">
                    <div className="truncate text-sm">
                      <span
                        className="font-semibold"
                        title={formatDate(purchase.createdAt)}
                      >
                        {formatTime(purchase.createdAt)}
                      </span>
                      {" · "}
                      {currency.format(purchase.totalGrossPrice.toNumber())}
                      {" · "}
                      {paymentMethodName(purchase.paymentMethod)}
                    </div>
                    {isPending ? (
                      <div className="text-xs font-semibold text-amber-800 dark:text-amber-300">
                        Pending on reader
                      </div>
                    ) : (
                      summary && (
                        <div className="truncate text-xs text-gray-600 dark:text-gray-400">
                          {summary}
                        </div>
                      )
                    )}
                  </div>

                  {purchase.status === "confirmed" &&
                    (isRefundable(purchase) ? (
                      <button
                        type="button"
                        className={`${actionButton} border-gray-300 bg-white text-gray-800 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100`}
                        aria-label={`Refund purchase from ${formatDate(purchase.createdAt)}`}
                        onClick={() => setModalState({ show: true, purchase })}
                        data-testid={`refund-purchase-${purchase.id}`}
                      >
                        <HiReceiptRefund className="h-4 w-4" aria-hidden="true" />
                        Refund
                      </button>
                    ) : (
                      <span
                        className="flex shrink-0 items-center px-2 text-gray-400"
                        title="Refunds are only possible within 15 minutes. Ask an admin."
                        aria-label="Refund window closed"
                        role="img"
                        data-testid={`refund-closed-${purchase.id}`}
                      >
                        <HiLockClosed className="h-4 w-4" aria-hidden="true" />
                      </span>
                    ))}

                  {isPending && (
                    <button
                      type="button"
                      className={`${actionButton} border-amber-400 bg-amber-100 text-amber-900 dark:border-amber-500 dark:bg-amber-800 dark:text-amber-100`}
                      aria-label="Resume purchase"
                      onClick={() => handleResume(purchase)}
                    >
                      <HiRefresh className="h-4 w-4" aria-hidden="true" />
                      Resume
                    </button>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </div>
  );
};

export default PurchaseHistory;
