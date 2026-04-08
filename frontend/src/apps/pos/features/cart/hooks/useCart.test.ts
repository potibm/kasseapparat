import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useCart } from "./useCart";
import { storePurchase } from "../../../utils/api";
import {
  Product as ProductType,
  Purchase as PurchaseType,
} from "../../../utils/api.schemas";
import {
  createMockProduct,
  createMockPurchase,
} from "@pos/utils/api.schemas.mocks";
import { PaymentMethodData } from "../types/cart.types";
import Decimal from "decimal.js";
import {
  getPurchaseErrorType,
  getErrorMessage,
  PurchaseErrorType,
} from "../services/PurchaseErrorHandler";

// --- 1. MOCKS ---
vi.mock("../../../utils/api", () => ({
  storePurchase: vi.fn(),
}));

vi.mock("../services/PurchaseErrorHandler", () => ({
  getPurchaseErrorType: vi.fn(),
  getErrorMessage: vi.fn(),
}));

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: () => ({
    currency: new Intl.NumberFormat(),
  }),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: vi.fn(() => ({
    debug: vi.fn(),
    warn: vi.fn(),
    info: vi.fn(),
    error: vi.fn(),
  })),
}));

const mockShowToast = vi.fn();
vi.mock("@pos/features/ui/toast/hooks/useToast", () => ({
  useToast: () => ({
    showToast: mockShowToast,
  }),
}));

// --- 2. FIXTURES (Dummy Data) ---
const mockApiHost = "https://api.example.com";
const mockGetToken = vi.fn(async () => "fake-token");

const mockProduct: ProductType = createMockProduct();
const mockPaymentData = { type: "empty" } as PaymentMethodData;

// --- 3. TESTS ---
describe("useCart Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Cart Manipulation (add, remove, clear)", () => {
    it("should initialize with an empty cart and default states", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      expect(result.current.cart.isEmpty).toBe(true);
      expect(result.current.isPolling).toBe(false);
      expect(result.current.pendingPurchase).toBeNull();
      expect(result.current.checkoutProcessing).toBeNull();
    });

    it("should add, remove, and clear products", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      act(() => {
        result.current.add(mockProduct, 2, null);
      });
      expect(result.current.cart.isEmpty).toBe(false);

      act(() => {
        result.current.remove(mockProduct);
      });
      expect(result.current.cart.isEmpty).toBe(true);

      act(() => {
        result.current.add(mockProduct, 1, null);
      });
      act(() => {
        result.current.clear();
      });
      expect(result.current.cart.isEmpty).toBe(true);
    });
  });

  describe("finalizeCheckout()", () => {
    it("should reset states and clear cart on success", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      act(() => {
        result.current.add(mockProduct, 1, null);
      });

      act(() => {
        result.current.finalizeCheckout(true);
      });

      expect(result.current.isPolling).toBe(false);
      expect(result.current.pendingPurchase).toBeNull();
      expect(result.current.checkoutProcessing).toBeNull();
      expect(result.current.cart.isEmpty).toBe(true);
    });

    it("should reset states but KEEP the cart on failure", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      act(() => {
        result.current.add(mockProduct, 1, null);
      });

      act(() => {
        result.current.finalizeCheckout(false);
      });

      expect(result.current.isPolling).toBe(false);
      expect(result.current.pendingPurchase).toBeNull();
      expect(result.current.checkoutProcessing).toBeNull();
      expect(result.current.cart.isEmpty).toBe(false);
    });
  });

  describe("checkout() API Interaction", () => {
    it("should handle an immediately confirmed purchase", async () => {
      const confirmedPurchase = createMockPurchase({
        status: "confirmed",
        totalGrossPrice: Decimal(979.66),
      });
      vi.mocked(storePurchase).mockResolvedValue(confirmedPurchase);

      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      act(() => {
        result.current.add(mockProduct, 1, null);
      });

      const expectedPayload = result.current.cart.toApiPayload(
        "cash",
        mockPaymentData,
      );

      // perform checkout
      let purchaseResult;
      await act(async () => {
        purchaseResult = await result.current.checkout("cash", mockPaymentData);
      });

      expect(storePurchase).toHaveBeenCalledWith(
        mockApiHost,
        "fake-token",
        expectedPayload,
      );
      expect(purchaseResult).toEqual(confirmedPurchase);

      expect(result.current.isPolling).toBe(false);
      expect(result.current.checkoutProcessing).toBeNull();
      expect(result.current.cart.isEmpty).toBe(true);

      expect(mockShowToast).toHaveBeenCalledWith({
        severity: "success",
        message: "Payment of 979.66 successful!",
      });
    });

    it("should handle a pending purchase (triggering polling)", async () => {
      const pendingPurchase = createMockPurchase({ status: "pending" });
      vi.mocked(storePurchase).mockResolvedValue(pendingPurchase);

      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      act(() => {
        result.current.add(mockProduct, 1, null);
      });

      let purchaseResult;
      await act(async () => {
        purchaseResult = await result.current.checkout(
          "sumup",
          mockPaymentData,
        );
      });

      expect(purchaseResult).toEqual(pendingPurchase);

      // On pending: polling starts, purchase is saved, cart remains full
      expect(result.current.isPolling).toBe(true);
      expect(result.current.pendingPurchase).toEqual(pendingPurchase);
      expect(result.current.checkoutProcessing).toBe("sumup");
      expect(result.current.cart.isEmpty).toBe(false);
    });

    it("should throw an error for unknown purchase status", async () => {
      const validPurchase = createMockPurchase();
      const weirdPurchase = {
        ...validPurchase,
        status: "aliens_attacked",
      } as unknown as PurchaseType;

      vi.mocked(storePurchase).mockResolvedValue(weirdPurchase);

      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      await act(async () => {
        await expect(
          result.current.checkout("cash", mockPaymentData),
        ).rejects.toThrow("Unknown purchase status: aliens_attacked");
      });
    });

    it("should reset checkoutProcessing and rethrow if API fails", async () => {
      const networkError = new Error("Network timeout");
      vi.mocked(storePurchase).mockRejectedValue(networkError);

      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      await act(async () => {
        await expect(
          result.current.checkout("cash", mockPaymentData),
        ).rejects.toThrow("Network timeout");
      });

      // After an error, it should not be stuck in "Processing" state
      expect(result.current.checkoutProcessing).toBeNull();
    });
  });

  describe("checkout() Specific Error Handling", () => {
    it("should show a blocking toast when error type is READER_BUSY", async () => {
      const mockError = new Error("Terminal is busy");
      vi.mocked(storePurchase).mockRejectedValue(mockError);
      vi.mocked(getPurchaseErrorType).mockReturnValue(
        "READER_BUSY" as unknown as PurchaseErrorType,
      );
      vi.mocked(getErrorMessage).mockReturnValue("Card reader is busy.");

      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      await act(async () => {
        await expect(
          result.current.checkout("sumup", mockPaymentData),
        ).rejects.toThrow("Terminal is busy");
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        severity: "error",
        message: "Card reader is busy.",
        blocking: true, // Das ist der kritische Pfad für die Coverage!
      });
    });

    it("should show a non-blocking toast for other errors", async () => {
      const mockError = new Error("General error");
      vi.mocked(storePurchase).mockRejectedValue(mockError);
      vi.mocked(getPurchaseErrorType).mockReturnValue(
        "GENERAL_ERROR" as unknown as PurchaseErrorType,
      );
      vi.mocked(getErrorMessage).mockReturnValue("Something went wrong.");

      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));

      await act(async () => {
        await expect(
          result.current.checkout("sumup", mockPaymentData),
        ).rejects.toThrow("General error");
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        severity: "error",
        message: "Something went wrong.",
        blocking: false,
      });
    });
  });

  describe("resumePolling()", () => {
    it("should resume polling for a pending purchase", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));
      const pendingPurchase = createMockPurchase({
        status: "pending",
        paymentMethod: "sumup",
      });

      act(() => {
        result.current.resumePolling(pendingPurchase);
      });

      expect(result.current.isPolling).toBe(true);
      expect(result.current.pendingPurchase).toEqual(pendingPurchase);
      expect(result.current.checkoutProcessing).toBe("sumup");
    });

    it("should early return and ignore if already polling", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));
      const purchase1 = createMockPurchase({
        status: "pending",
        id: "purchase-1",
      });
      const purchase2 = createMockPurchase({
        status: "pending",
        id: "purchase-2",
      });

      // Starte initiales Polling
      act(() => {
        result.current.resumePolling(purchase1);
      });

      // Versuch, zweites Polling zu starten
      act(() => {
        result.current.resumePolling(purchase2);
      });

      // Erwartung: Zustand bleibt beim ersten Purchase
      expect(result.current.pendingPurchase?.id).toBe("purchase-1");
    });

    it("should early return if purchase status is not pending", () => {
      const { result } = renderHook(() => useCart(mockApiHost, mockGetToken));
      const confirmedPurchase = createMockPurchase({ status: "confirmed" });

      act(() => {
        result.current.resumePolling(confirmedPurchase);
      });

      expect(result.current.isPolling).toBe(false);
      expect(result.current.pendingPurchase).toBeNull();
      expect(result.current.checkoutProcessing).toBeNull();
    });
  });
});
