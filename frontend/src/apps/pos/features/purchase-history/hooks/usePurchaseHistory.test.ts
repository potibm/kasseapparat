import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { usePurchaseHistory } from "./usePurchaseHistory";
import { fetchPurchases, refundPurchaseById } from "../../../utils/api";
import { createMockPurchase } from "@pos/utils/api.schemas.mocks";
import Decimal from "decimal.js";

// mocks
vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: () => ({
    currency: new Intl.NumberFormat(),
  }),
}));

vi.mock("../../../utils/api", () => ({
  fetchPurchases: vi.fn(),
  refundPurchaseById: vi.fn(),
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
const mockApiHost = "https://api.example.com";
const mockGetToken = vi.fn(async () => "fake-token");
const mockUserId = 42;

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

describe("usePurchaseHistory Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Initialization (loadHistory)", () => {
    it("should NOT fetch anything and return an empty array if userId is falsy", async () => {
      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, 0),
      );

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(fetchPurchases).not.toHaveBeenCalled();
      expect(result.current.history).toEqual([]);
    });

    it("should fetch and load purchase history automatically on mount", async () => {
      vi.mocked(fetchPurchases).mockResolvedValue(mockPurchases);

      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, mockUserId),
      );

      expect(result.current.loading).toBe(true);

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(fetchPurchases).toHaveBeenCalledWith(
        mockApiHost,
        "fake-token",
        mockUserId,
      );
      expect(result.current.history).toEqual(mockPurchases);
      expect(mockShowToast).not.toHaveBeenCalled();
    });

    it("should trigger onError and set empty history if fetching throws an Error object", async () => {
      const apiError = new Error("Database unreachable");
      vi.mocked(fetchPurchases).mockRejectedValue(apiError);

      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, mockUserId),
      );

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
      vi.mocked(fetchPurchases).mockRejectedValue("Weird backend crash string");

      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, mockUserId),
      );

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
      vi.mocked(fetchPurchases).mockResolvedValue(mockPurchases);
      vi.mocked(refundPurchaseById).mockResolvedValue(createMockPurchase());

      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, mockUserId),
      );

      await waitFor(() => expect(result.current.loading).toBe(false));

      vi.mocked(fetchPurchases).mockClear();

      await act(async () => {
        await result.current.refundPurchase("purchase-123");
      });

      expect(refundPurchaseById).toHaveBeenCalledWith(
        mockApiHost,
        "fake-token",
        "purchase-123",
      );
      expect(fetchPurchases).toHaveBeenCalledTimes(1);
    });

    it("should trigger onError AND re-throw the error if refunding fails with an Error", async () => {
      vi.mocked(fetchPurchases).mockResolvedValue(mockPurchases);

      const refundError = new Error("Refund denied by bank");
      vi.mocked(refundPurchaseById).mockRejectedValue(refundError);

      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, mockUserId),
      );

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
      vi.mocked(fetchPurchases).mockResolvedValue(mockPurchases);
      vi.mocked(refundPurchaseById).mockRejectedValue({ some: "weird object" });

      const { result } = renderHook(() =>
        usePurchaseHistory(mockApiHost, mockGetToken, mockUserId),
      );

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
