// src/apps/pos/features/purchase-history/hooks/usePurchaseHistory.ts
import { useState, useEffect, useCallback, useMemo } from "react";
import { fetchPurchases, refundPurchaseById } from "../../../utils/api";
import { Purchase as PurchaseType } from "../../../utils/api.schemas";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { useConfig } from "@core/config/hooks/useConfig";

const log = createLogger("Purchase");

export const usePurchaseHistory = (
  apiHost: string,
  getToken: () => Promise<string>,
  userId: number,
) => {
  const [history, setHistory] = useState<PurchaseType[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const { showToast } = useToast();
  const { currency } = useConfig();

  /**
   * Load history of purchases for the current user.
   * @param isSilent when true, the function will not show loading states or error toasts. Useful for background refreshes.
   */
  const loadHistory = useCallback(
    async (isSilent = false) => {
      if (!userId) {
        log.warn("No user ID provided, cannot load purchase history");
        if (!isSilent) {
          setHistory([]);
          setLoading(false);
        }
        return;
      }

      if (!isSilent) {
        setLoading(true);
        if (!history) setHistory([]);
      }

      try {
        const token = await getToken();
        const purchases = await fetchPurchases(apiHost, token, userId);

        setHistory(purchases);
        log.debug("Purchase history fetched successfully", {
          purchaseCount: purchases.length,
          silent: isSilent,
        });
      } catch (error: unknown) {
        log.error(
          "Error fetching purchase history",
          error instanceof Error ? { message: error.message } : { error },
        );

        if (!isSilent) {
          const errorMessage =
            error instanceof Error
              ? "Error while loading the purchase history: " + error.message
              : "An unknown error has occurred";

          showToast({
            severity: "error",
            message: errorMessage,
            autoClose: false,
          });
          setHistory([]);
        }
      } finally {
        if (!isSilent) {
          setLoading(false);
        }
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [apiHost, getToken, userId, showToast],
  );

  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  const hasPendingPurchases = useMemo(() => {
    return history?.some((p) => p.status === "pending") ?? false;
  }, [history]);

  useEffect(() => {
    if (!hasPendingPurchases) return;

    log.debug("Pending purchases detected. Starting background polling...");

    const intervalId = setInterval(() => {
      loadHistory(true);
    }, 5000);

    return () => {
      log.debug("Stopping background polling.");
      clearInterval(intervalId);
    };
  }, [hasPendingPurchases, loadHistory]);

  const refund = async (purchaseId: string) => {
    try {
      const token = await getToken();
      const purchase = await refundPurchaseById(apiHost, token, purchaseId);
      showToast({
        severity: "success",
        message: `Purchase of ${currency.format(purchase.totalGrossPrice.toNumber())} refunded successfully!`,
      });
      await loadHistory();
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "Error while refunding the purchase: " + error.message
          : "An unknown error has occurred";

      showToast({
        severity: "error",
        message: errorMessage,
        autoClose: false,
        blocking: true,
      });
      throw error;
    }
  };

  return {
    history,
    loading,
    refreshHistory: loadHistory,
    refundPurchase: refund,
  };
};
