// src/apps/pos/features/purchase-history/hooks/usePurchaseHistory.ts
import { useState, useEffect, useCallback } from "react";
import { fetchPurchases, refundPurchaseById } from "../../../utils/api";
import { Purchase } from "../types/purchase.types";
import Decimal from "decimal.js";

export const usePurchaseHistory = (
  apiHost: string,
  getToken: () => Promise<string>,
  userId: string,
  onError: (msg: string) => void,
) => {
  const [history, setHistory] = useState<Purchase[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const loadHistory = useCallback(async () => {
    if (!userId) return;
    setLoading(true);
    try {
      const token = await getToken();
      const rawData = await fetchPurchases(apiHost, token, userId);

      const convertedData: Purchase[] = (rawData as any[]).map((p: any) => ({
        ...p,
        totalGrossPrice: new Decimal(p.totalGrossPrice),
        totalNetPrice: new Decimal(p.totalNetPrice),
        totalVatAmount: new Decimal(p.totalVatAmount),
        purchaseItems: p.purchaseItems.map((item: any) => ({
          ...item,
          netPrice: new Decimal(item.netPrice),
          grossPrice: new Decimal(item.grossPrice),
          vatAmount: new Decimal(item.vatAmount),
          totalNetPrice: new Decimal(item.totalNetPrice),
          totalGrossPrice: new Decimal(item.totalGrossPrice),
          totalVatAmount: new Decimal(item.totalVatAmount),
        })),
      }));

      setHistory(convertedData);
    } catch (error: any) {
      onError("Error while loading the purchase history: " + error.message);
    } finally {
      setLoading(false);
    }
  }, [apiHost, getToken, userId, onError]);

  // Initiales Laden beim Mounten oder wenn sich die UserID ändert
  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  const refund = async (purchaseId: string) => {
    try {
      const token = await getToken();
      await refundPurchaseById(apiHost, token, purchaseId);
      await loadHistory();
    } catch (error: any) {
      onError("Error while refunding the purchase: " + error.message);
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
