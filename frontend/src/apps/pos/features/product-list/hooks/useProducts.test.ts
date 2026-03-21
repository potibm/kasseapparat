import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useProducts } from "./useProducts";
import { fetchProducts, addProductInterest } from "../../../utils/api";
import {
  Product as ProductType,
  ProductInterest as ProductInterestType,
} from "../../../utils/api.schemas";
import { createMockProduct } from "@pos/utils/api.schemas.mocks";

// mocks
vi.mock("../../../utils/api", () => ({
  fetchProducts: vi.fn(),
  addProductInterest: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
  }),
}));

// fixtures
const mockApiHost = "https://api.example.com";
const mockGetToken = vi.fn(async () => "fake-token");
const mockOnError = vi.fn();

const mockProducts = [
  createMockProduct({ id: 1, name: "Product A" }),
  createMockProduct({ id: 2, name: "Product B" }),
] as ProductType[];

describe("useProducts Hook", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Initialization (loadProducts)", () => {
    it("should fetch and load products automatically on mount", async () => {
      vi.mocked(fetchProducts).mockResolvedValue(mockProducts);

      const { result } = renderHook(() =>
        useProducts(mockApiHost, mockGetToken, mockOnError),
      );

      expect(result.current.loading).toBe(true);

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(fetchProducts).toHaveBeenCalledWith(mockApiHost, "fake-token");
      expect(result.current.products).toEqual(mockProducts);
      expect(mockOnError).not.toHaveBeenCalled();
    });

    it("should trigger onError and set loading to false if fetching throws an Error object", async () => {
      const apiError = new Error("Network offline");
      vi.mocked(fetchProducts).mockRejectedValue(apiError);

      const { result } = renderHook(() =>
        useProducts(mockApiHost, mockGetToken, mockOnError),
      );

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.products).toBeNull();
      expect(mockOnError).toHaveBeenCalledWith(
        "There was an error fetching the products: Network offline",
      );
    });

    it("should trigger onError with a fallback message if fetching throws a non-Error (unknown)", async () => {
      vi.mocked(fetchProducts).mockRejectedValue("Some weird string error");

      const { result } = renderHook(() =>
        useProducts(mockApiHost, mockGetToken, mockOnError),
      );

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(mockOnError).toHaveBeenCalledWith("An unknown error has occurred");
    });
  });

  describe("addInterest()", () => {
    it("should call the API and then reload the products", async () => {
      vi.mocked(fetchProducts).mockResolvedValue(mockProducts);
      vi.mocked(addProductInterest).mockResolvedValue(
        undefined as unknown as ProductInterestType,
      );

      const { result } = renderHook(() =>
        useProducts(mockApiHost, mockGetToken, mockOnError),
      );

      await waitFor(() => expect(result.current.loading).toBe(false));

      vi.mocked(fetchProducts).mockClear();

      await act(async () => {
        await result.current.addInterest(99);
      });

      expect(addProductInterest).toHaveBeenCalledWith(
        mockApiHost,
        "fake-token",
        99,
      );

      expect(fetchProducts).toHaveBeenCalledTimes(1);
    });

    it("should trigger onError if adding interest fails with an Error object", async () => {
      vi.mocked(fetchProducts).mockResolvedValue(mockProducts);
      vi.mocked(addProductInterest).mockRejectedValue(
        new Error("Item not found"),
      );

      const { result } = renderHook(() =>
        useProducts(mockApiHost, mockGetToken, mockOnError),
      );

      await waitFor(() => expect(result.current.loading).toBe(false));

      await act(async () => {
        await result.current.addInterest(99);
      });

      expect(mockOnError).toHaveBeenCalledWith(
        "Error on saving the interest: Item not found",
      );
    });

    it("should trigger onError with a fallback message if adding interest throws a non-Error", async () => {
      vi.mocked(fetchProducts).mockResolvedValue(mockProducts);
      vi.mocked(addProductInterest).mockRejectedValue(12345);

      const { result } = renderHook(() =>
        useProducts(mockApiHost, mockGetToken, mockOnError),
      );

      await waitFor(() => expect(result.current.loading).toBe(false));

      await act(async () => {
        await result.current.addInterest(99);
      });

      expect(mockOnError).toHaveBeenCalledWith("An unknown error has occurred");
    });
  });
});
