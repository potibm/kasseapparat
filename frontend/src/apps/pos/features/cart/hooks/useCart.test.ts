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

// --- 1. MOCKS ---
vi.mock("../../../utils/api", () => ({
  storePurchase: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: vi.fn(() => ({
    debug: vi.fn(),
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
      const confirmedPurchase = createMockPurchase({ status: "confirmed" });
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
        type: "success",
        message: "Purchase completed successfully!",
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
});
