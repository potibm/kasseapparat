// src/apps/pos/features/purchase-history/hooks/usePurchaseHistory.ts
import { useState, useEffect, useCallback } from "react";
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

  const loadHistory = useCallback(async () => {
    if (!userId) {
      setHistory([]);
      setLoading(false);
      return;
    }

    setLoading(true);
    setHistory([]);

    try {
      const token = await getToken();
      const purchases = await fetchPurchases(apiHost, token, userId);

      setHistory(purchases);
      log.debug("Purchase history fetched successfully", {
        purchaseCount: purchases.length,
      });
    } catch (error: unknown) {
      log.error(
        "Error fetching purchase history",
        error instanceof Error ? { message: error.message } : { error },
      );
      const errorMessage =
        error instanceof Error
          ? "Error while loading the purchase history: " + error.message
          : "An unknown error has occurred";

      showToast({ type: "error", message: errorMessage, autoClose: false });
      setHistory([]);
    } finally {
      setLoading(false);
    }
  }, [apiHost, getToken, userId, showToast]);

  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  const refund = async (purchaseId: string) => {
    try {
      const token = await getToken();
      const purchase = await refundPurchaseById(apiHost, token, purchaseId);
      showToast({
        type: "success",
        message: `Purchase of ${currency.format(purchase.totalGrossPrice.toNumber())} refunded successfully!`,
      });
      await loadHistory();
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "Error while refunding the purchase: " + error.message
          : "An unknown error has occurred";

      showToast({
        type: "error",
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
