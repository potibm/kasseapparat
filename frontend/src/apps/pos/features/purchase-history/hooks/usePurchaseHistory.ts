// src/apps/pos/features/purchase-history/hooks/usePurchaseHistory.ts
import { useState, useEffect, useCallback } from "react";
import { fetchPurchases, refundPurchaseById } from "../../../utils/api";
import { Purchase as PurchaseType } from "../../../utils/api.schemas";

export const usePurchaseHistory = (
  apiHost: string,
  getToken: () => Promise<string>,
  userId: number,
  onError: (msg: string) => void,
) => {
  const [history, setHistory] = useState<PurchaseType[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const loadHistory = useCallback(async () => {
    if (!userId) return;
    setLoading(true);
    try {
      const token = await getToken();
      const purchases = await fetchPurchases(apiHost, token, userId);

      setHistory(purchases);
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "Error while loading the purchase history: " + error.message
          : "An unknown error has occurred";

      onError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [apiHost, getToken, userId, onError]);

  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  const refund = async (purchaseId: string) => {
    try {
      const token = await getToken();
      await refundPurchaseById(apiHost, token, purchaseId);
      await loadHistory();
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "Error while refunding the purchase: " + error.message
          : "An unknown error has occurred";

      onError(errorMessage);
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
