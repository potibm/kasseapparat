import { describe, it, expect, vi, beforeEach } from "vitest";
import * as Sentry from "@sentry/react";
import { createPosApiClient } from "./factory";
import { Purchase as PurchaseType } from "./schemas";
import Decimal from "decimal.js";
import { ApiCreatePayloadPurchase } from "./types";
import {
  createMockProduct,
  createMockGuest,
  createMockPurchase,
} from "./schemas.mocks";

vi.mock("@sentry/react", () => ({
  captureException: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    error: vi.fn(),
    warn: vi.fn(),
  }),
}));

/* eslint-disable @typescript-eslint/no-explicit-any */
const convertDecimalsToStrings = (obj: any): any => {
  if (obj instanceof Decimal) {
    return obj.toFixed(2);
  } else if (Array.isArray(obj)) {
    return obj.map(convertDecimalsToStrings);
  } else if (obj !== null && typeof obj === "object") {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    const converted: any = {};
    for (const key in obj) {
      converted[key] = convertDecimalsToStrings(obj[key]);
    }
    return converted;
  }
  return obj;
};

describe("POS API Client", () => {
  const apiBaseUrl = "https://api.example.com/api/v2";

  beforeEach(() => {
    vi.restoreAllMocks();
    vi.stubGlobal("fetch", vi.fn());
  });

  describe("createPosApiClient", () => {
    it("should create a client with all expected methods", () => {
      const client = createPosApiClient(apiBaseUrl);

      expect(client.fetchProducts).toBeDefined();
      expect(client.fetchGuestlistByProductId).toBeDefined();
      expect(client.storePurchase).toBeDefined();
      expect(client.fetchPurchases).toBeDefined();
      expect(client.refundPurchaseById).toBeDefined();
      expect(client.addProductInterest).toBeDefined();
    });
  });

  describe("URL building", () => {
    it("should handle trailing slash in base URL", async () => {
      const client = createPosApiClient("https://api.example.com/api/v2/");

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => [],
        } as Response),
      );

      await client.fetchProducts();

      expect(fetch).toHaveBeenCalledWith(
        "https://api.example.com/api/v2/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true",
        expect.any(Object),
      );
    });

    it("should handle leading slash in endpoint", async () => {
      const client = createPosApiClient(apiBaseUrl);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => [],
        } as Response),
      );

      await client.fetchProducts();

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/products"),
        expect.any(Object),
      );
    });
  });

  describe("fetchProducts", () => {
    it("should return products on successful response", async () => {
      const client = createPosApiClient(apiBaseUrl);
      const mockProducts = [
        createMockProduct({ id: 1, name: "Product A" }),
        createMockProduct({ id: 2, name: "Product B" }),
      ];
      const requestProducts = convertDecimalsToStrings(mockProducts);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => requestProducts,
        } as Response),
      );

      const result = await client.fetchProducts();
      const url = `${apiBaseUrl}/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "GET",
        }),
      );

      expect(result).toEqual(mockProducts);
    });
  });

  describe("fetchGuestlistByProductId", () => {
    it("should return guests on successful response", async () => {
      const client = createPosApiClient(apiBaseUrl);
      const mockGuests = [
        createMockGuest({ id: 1, name: "Guest A" }),
        createMockGuest({ id: 2, name: "Guest B" }),
      ];
      const requestGuests = convertDecimalsToStrings(mockGuests);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => requestGuests,
        } as Response),
      );

      const result = await client.fetchGuestlistByProductId(12, "Hans");
      const url = `${apiBaseUrl}/products/12/guests?q=Hans`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "GET",
        }),
      );

      expect(result).toEqual(mockGuests);
    });

    it("should return empty array if no guests found", async () => {
      const client = createPosApiClient(apiBaseUrl);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => null,
        } as Response),
      );

      const result = await client.fetchGuestlistByProductId(12, "Hans");

      expect(result).toEqual([]);
    });
  });

  describe("storePurchase", () => {
    it("should return purchase on successful response", async () => {
      const client = createPosApiClient(apiBaseUrl);
      const mockPurchase = createMockPurchase();
      const requestPurchase = convertDecimalsToStrings(mockPurchase);

      const createPurchasePayload: ApiCreatePayloadPurchase = {
        paymentMethod: "cash",
        cart: [
          {
            id: 1,
            quantity: 2,
            lists: null,
            guestlists: null,
          },
        ],
        totalGrossPrice: "100.00",
        totalNetPrice: "80.00",
      };

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => requestPurchase,
        } as Response),
      );

      const result = await client.storePurchase(createPurchasePayload);
      const url = `${apiBaseUrl}/purchases`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify(createPurchasePayload),
        }),
      );

      expect(result).toEqual(mockPurchase);
    });
  });

  describe("fetchPurchases", () => {
    it("should return purchases on successful response", async () => {
      const client = createPosApiClient(apiBaseUrl);
      const mockPurchases: PurchaseType[] = [
        createMockPurchase({ id: "123e4567-e89b-12d3-a456-426614174000" }),
        createMockPurchase({ id: "123e4567-e89b-12d3-a456-426614174001" }),
      ];
      const requestPurchase = convertDecimalsToStrings(mockPurchases);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => requestPurchase,
        } as Response),
      );

      const result = await client.fetchPurchases("testuser");
      const url = `${apiBaseUrl}/purchases?createdById=testuser&status=confirmed&status=pending`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "GET",
        }),
      );

      expect(result).toEqual(mockPurchases);
    });
  });

  describe("refundPurchaseById", () => {
    it("should return purchase on successful response", async () => {
      const client = createPosApiClient(apiBaseUrl);
      const purchaseId = "123e4567-e89b-12d3-a456-426614174000";
      const mockPurchase = createMockPurchase({
        id: purchaseId,
        status: "refunded",
      });
      const requestPurchase = convertDecimalsToStrings(mockPurchase);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => requestPurchase,
        } as Response),
      );

      const result = await client.refundPurchaseById(purchaseId);
      const url = `${apiBaseUrl}/purchases/${purchaseId}/refund`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "POST",
        }),
      );

      expect(result).toEqual(mockPurchase);
    });
  });

  describe("addProductInterest", () => {
    it("should return id on successful response", async () => {
      const client = createPosApiClient(apiBaseUrl);
      const productInterest = { id: 98 };

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => productInterest,
        } as Response),
      );

      const result = await client.addProductInterest(12);
      const url = `${apiBaseUrl}/productInterests`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({ productId: 12 }),
        }),
      );

      expect(result).toEqual(productInterest);
    });
  });

  describe("error handling", () => {
    it("should throw an error if Zod validation fails", async () => {
      const client = createPosApiClient(apiBaseUrl);

      vi.mocked(fetch).mockResolvedValue({
        ok: true,
        json: async () => ({ wrong_key: "garbage" }),
      } as Response);

      await expect(client.fetchProducts()).rejects.toThrow(
        "API Response format mismatch",
      );
    });

    it("should throw an error on non-ok response and call sentry", async () => {
      const client = createPosApiClient(apiBaseUrl);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: false,
          status: 500,
          statusText: "Internal Server Error",
          json: async () => ({ message: "Server error" }),
        } as Response),
      );

      await expect(client.fetchProducts()).rejects.toThrow("Server error");

      expect(Sentry.captureException).toHaveBeenCalled();
    });
  });
});
