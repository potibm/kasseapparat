import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { usePurchaseHistory } from "./usePurchaseHistory";
import { createMockPurchase } from "@pos/api/schemas.mocks";
import Decimal from "decimal.js";
import React from "react";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";

// mocks
const mockFetchPurchases = vi.fn();
const mockRefundPurchaseById = vi.fn();

vi.mock("@pos/api/usePosApi", () => ({
  usePosApi: () => ({
    fetchPurchases: mockFetchPurchases,
    refundPurchaseById: mockRefundPurchaseById,
  }),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
    warn: vi.fn(),
  }),
}));

const mockShowToast = vi.fn();
vi.mock("@pos/features/ui/toast/hooks/useToast", () => ({
  useToast: () => ({
    showToast: mockShowToast,
  }),
}));

// fixture data
const mockConfig: AppConfig = {
  version: "1.0.0",
  apiHost: "https://api.example.com",
  apiBaseUrl: "https://api.example.com/api/v2",
  websocketHost: "wss://api.example.com",
  websocketBaseUrl: "wss://api.example.com/api/v2",
  locale: "en",
  currencyCode: "USD",
  currencyLocale: "en-US",
  currency: new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }),
  currencyOptions: {},
  dateLocale: "en-US",
  dateOptions: {},
  vatRates: [],
  paymentMethods: [],
  sumupEnabled: false,
  authMode: "proxy",
};

const mockUsername = "testuser";

const mockPurchases = [
  createMockPurchase({
    id: "p-1",
    totalGrossPrice: new Decimal("10.00"),
    status: "confirmed",
  }),
  createMockPurchase({
    id: "p-2",
    totalGrossPrice: new Decimal("25.50"),
    status: "pending",
  }),
];

const wrapper = ({ children }: { children: React.ReactNode }) =>
  React.createElement(ConfigContext.Provider, { value: mockConfig }, children);

describe("usePurchaseHistory Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Initialization (loadHistory)", () => {
    it("should NOT fetch anything and return an empty array if username is falsy", async () => {
      const { result } = renderHook(() => usePurchaseHistory(""), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(mockFetchPurchases).not.toHaveBeenCalled();
      expect(result.current.history).toEqual([]);
    });

    it("should fetch and load purchase history automatically on mount", async () => {
      mockFetchPurchases.mockResolvedValue(mockPurchases);

      const { result } = renderHook(() => usePurchaseHistory(mockUsername), {
        wrapper,
      });

      expect(result.current.loading).toBe(true);

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(mockFetchPurchases).toHaveBeenCalledWith(mockUsername);
      expect(result.current.history).toEqual(mockPurchases);
      expect(mockShowToast).not.toHaveBeenCalled();
    });

    it("should trigger onError and set empty history if fetching throws an Error object", async () => {
      const apiError = new Error("Database unreachable");
      mockFetchPurchases.mockRejectedValue(apiError);

      const { result } = renderHook(() => usePurchaseHistory(mockUsername), {
        wrapper,
      });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.history).toEqual([]);
      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        message:
          "Error while loading the purchase history: Database unreachable",
        severity: "error",
      });
    });

    it("should trigger onError with a fallback message if fetching throws a non-Error", async () => {
      mockFetchPurchases.mockRejectedValue("Weird backend crash string");

      const { result } = renderHook(() => usePurchaseHistory(mockUsername), {
        wrapper,
      });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.history).toEqual([]);
      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        message: "An unknown error has occurred",
        severity: "error",
      });
    });
  });

  describe("refundPurchase()", () => {
    it("should call the refund API and then reload the history", async () => {
      mockFetchPurchases.mockResolvedValue(mockPurchases);
      mockRefundPurchaseById.mockResolvedValue(createMockPurchase());

      const { result } = renderHook(() => usePurchaseHistory(mockUsername), {
        wrapper,
      });

      await waitFor(() => expect(result.current.loading).toBe(false));

      vi.mocked(mockFetchPurchases).mockClear();

      await act(async () => {
        await result.current.refundPurchase("purchase-123");
      });

      expect(mockRefundPurchaseById).toHaveBeenCalledWith("purchase-123");
      expect(mockFetchPurchases).toHaveBeenCalledTimes(1);
    });

    it("should trigger onError AND re-throw the error if refunding fails with an Error", async () => {
      mockFetchPurchases.mockResolvedValue(mockPurchases);

      const refundError = new Error("Refund denied by bank");
      mockRefundPurchaseById.mockRejectedValue(refundError);

      const { result } = renderHook(() => usePurchaseHistory(mockUsername), {
        wrapper,
      });

      await waitFor(() => expect(result.current.loading).toBe(false));

      await act(async () => {
        await expect(
          result.current.refundPurchase("purchase-123"),
        ).rejects.toThrow("Refund denied by bank");
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        blocking: true,
        message: "Error while refunding the purchase: Refund denied by bank",
        severity: "error",
      });
    });

    it("should trigger onError AND re-throw with fallback if refunding throws a non-Error", async () => {
      mockFetchPurchases.mockResolvedValue(mockPurchases);
      mockRefundPurchaseById.mockRejectedValue({ some: "weird object" });

      const { result } = renderHook(() => usePurchaseHistory(mockUsername), {
        wrapper,
      });

      await waitFor(() => expect(result.current.loading).toBe(false));

      await act(async () => {
        await expect(
          result.current.refundPurchase("purchase-123"),
        ).rejects.toEqual({ some: "weird object" });
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        blocking: true,
        message: "An unknown error has occurred",
        severity: "error",
      });
    });
  });
});
