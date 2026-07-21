import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useProducts } from "./useProducts";
import {
  Product as ProductType,
  ProductInterest as ProductInterestType,
} from "../../../api/schemas";
import { createMockProduct } from "@pos/api/schemas.mocks";
import React from "react";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";

// mocks
const mockFetchProducts = vi.fn();
const mockAddProductInterest = vi.fn();

vi.mock("@pos/api/usePosApi", () => ({
  usePosApi: () => ({
    fetchProducts: mockFetchProducts,
    addProductInterest: mockAddProductInterest,
  }),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
  }),
}));

const mockShowToast = vi.fn();
vi.mock("@pos/features/ui/toast/hooks/useToast", () => ({
  useToast: () => ({
    showToast: mockShowToast,
  }),
}));

// fixtures
const mockConfig: AppConfig = {
  version: "1.0.0",
  apiHost: "https://api.example.com",
  apiBaseUrl: "https://api.example.com/api/v3",
  websocketHost: "wss://api.example.com",
  websocketBaseUrl: "wss://api.example.com/api/v3",
  locale: "en",
  currencyCode: "USD",
  currencyLocale: "en-US",
  currency: new Intl.NumberFormat("en-US"),
  currencyOptions: {},
  dateLocale: "en-US",
  dateOptions: {},
  vatRates: [],
  paymentMethods: [],
  sumupEnabled: false,
  authMode: "proxy",
};

const mockProducts = [
  createMockProduct({ id: 1, name: "Product A" }),
  createMockProduct({ id: 2, name: "Product B" }),
] as ProductType[];

const wrapper = ({ children }: { children: React.ReactNode }) =>
  React.createElement(ConfigContext.Provider, { value: mockConfig }, children);

describe("useProducts Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Initialization (loadProducts)", () => {
    it("should fetch and load products automatically on mount", async () => {
      mockFetchProducts.mockResolvedValue(mockProducts);

      const { result } = renderHook(() => useProducts(), { wrapper });

      expect(result.current.loading).toBe(true);

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(mockFetchProducts).toHaveBeenCalled();
      expect(result.current.products).toEqual(mockProducts);
      expect(mockShowToast).not.toHaveBeenCalled();
    });

    it("should trigger onError and set loading to false if fetching throws an Error object", async () => {
      const apiError = new Error("Network offline");
      mockFetchProducts.mockRejectedValue(apiError);

      const { result } = renderHook(() => useProducts(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.products).toBeNull();
      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        message: "There was an error fetching the products: Network offline",
        severity: "error",
      });
    });

    it("should trigger onError with a fallback message if fetching throws a non-Error (unknown)", async () => {
      mockFetchProducts.mockRejectedValue("Some weirdstring error");

      const { result } = renderHook(() => useProducts(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        message: "An unknown error has occurred",
        severity: "error",
      });
    });
  });

  describe("addInterest()", () => {
    it("should call the API and then reload the products", async () => {
      mockFetchProducts.mockResolvedValue(mockProducts);
      mockAddProductInterest.mockResolvedValue(
        undefined as unknown as ProductInterestType,
      );

      const { result } = renderHook(() => useProducts(), { wrapper });

      await waitFor(() => expect(result.current.loading).toBe(false));

      vi.mocked(mockFetchProducts).mockClear();

      await act(async () => {
        await result.current.addInterest(99, "Test Product");
      });

      expect(mockAddProductInterest).toHaveBeenCalledWith(99);

      expect(mockFetchProducts).toHaveBeenCalledTimes(1);
    });

    it("should trigger onError if adding interest fails with an Error object", async () => {
      mockFetchProducts.mockResolvedValue(mockProducts);
      mockAddProductInterest.mockRejectedValue(new Error("Item not found"));

      const { result } = renderHook(() => useProducts(), { wrapper });

      await waitFor(() => expect(result.current.loading).toBe(false));

      await act(async () => {
        await result.current.addInterest(99, "A missing product");
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        message: "Error on saving the interest: Item not found",
        autoClose: false,
        severity: "error",
      });
    });

    it("should trigger onError with a fallback message if adding interest throws a non-Error", async () => {
      mockFetchProducts.mockResolvedValue(mockProducts);
      mockAddProductInterest.mockRejectedValue(12345);

      const { result } = renderHook(() => useProducts(), { wrapper });

      await waitFor(() => expect(result.current.loading).toBe(false));

      await act(async () => {
        await result.current.addInterest(99, "Another strange product");
      });

      expect(mockShowToast).toHaveBeenCalledWith({
        autoClose: false,
        message: "An unknown error has occurred",
        severity: "error",
      });
    });
  });
});
