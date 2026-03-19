import { describe, it, expect, vi, beforeEach } from "vitest";
import * as Sentry from "@sentry/react";
import {
  fetchProducts,
  fetchGuestlistByProductId,
  storePurchase,
  fetchPurchases,
  refundPurchaseById,
  addProductInterest,
} from "./api";
import {
  Product as ProductType,
  Guest as GuestType,
  Purchase as PurchaseType,
} from "./api.schemas";
import Decimal from "decimal.js";
import { ApiCreatePayloadPurchase } from "./api.types";

// mock external dependencies
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

const createMockProduct = (overrides?: Partial<ProductType>): ProductType => {
  return {
    id: 1,
    name: "Test Product",
    netPrice: new Decimal("10.00"),
    grossPrice: new Decimal("12.00"),
    vatRate: new Decimal("20.00"),
    vatAmount: new Decimal("2.00"),
    wrapAfter: false,
    hidden: false,
    soldOut: false,
    apiExport: true,
    pos: 1,
    totalStock: 100,
    guestlists: null,
    unitsSold: 0,
    soldOutRequestCount: 0,
    ...overrides,
  };
};

const createMockGuest = (overrides?: Partial<GuestType>): GuestType => {
  return {
    id: 1,
    name: "Test Guest",
    code: null,
    listName: "Test List",
    additionalGuests: 0,
    arrivalNote: null,
    attendedGuests: 0,
    ...overrides,
  };
};

const createRawMockPurchase = (
  overrides?: Partial<PurchaseType>,
): PurchaseType => {
  return {
    id: "123e4567-e89b-12d3-a456-426614174000",
    createdAt: new Date().toISOString(),
    createdById: 1,
    createdBy: {
      id: 1,
      username: "testuser",
      email: "test@example.com",
      admin: false,
    },
    paymentMethod: "cash",
    totalNetPrice: new Decimal("80.00"),
    totalGrossPrice: new Decimal("100.00"),
    totalVatAmount: new Decimal("20.00"),
    sumupTransactionId: null,
    sumupClientTransactionId: null,
    status: "pending",
    purchaseItems: [],
    ...overrides,
  };
};

describe("Api Service", () => {
  const apiHost = "https://api.example.com";
  const fakeToken = "fake-token";

  beforeEach(() => {
    vi.restoreAllMocks();
    // Stub global fetch
    vi.stubGlobal("fetch", vi.fn());
  });

  describe("fetchProducts", () => {
    it("should return products on successful response", async () => {
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

      const result = await fetchProducts(apiHost, fakeToken);
      const url = `${apiHost}/api/v2/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "GET",
          headers: expect.objectContaining({
            Authorization: "Bearer " + fakeToken,
          }),
        }),
      );

      expect(result).toEqual(mockProducts);
    });
  });

  describe("fetchGuestlistByProductId", () => {
    it("should return guests on successful response", async () => {
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

      const result = await fetchGuestlistByProductId(
        apiHost,
        fakeToken,
        12,
        "Hans",
      );
      const url = `${apiHost}/api/v2/products/12/guests?q=Hans`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "GET",
          headers: expect.objectContaining({
            "Content-Type": "application/json",
            Authorization: "Bearer " + fakeToken,
          }),
        }),
      );

      expect(result).toEqual(mockGuests);
    });

    it("should return empty array if no guests found", async () => {
      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => null,
        } as Response),
      );

      const result = await fetchGuestlistByProductId(
        apiHost,
        fakeToken,
        12,
        "Hans",
      );

      expect(result).toEqual([]);
    });
  });

  describe("storePurchase", () => {
    it("should return purchase on successful response ", async () => {
      const mockPurchase = createRawMockPurchase();
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

      const result = await storePurchase(
        apiHost,
        fakeToken,
        createPurchasePayload,
      );
      const url = `${apiHost}/api/v2/purchases`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "POST",
          headers: expect.objectContaining({
            "Content-Type": "application/json",
            Authorization: "Bearer " + fakeToken,
          }),
          body: JSON.stringify(createPurchasePayload),
        }),
      );

      expect(result).toEqual(mockPurchase);
    });
  });

  describe("fetchPurchases", () => {
    it("should return purchases on successful response", async () => {
      const mockPurchases: PurchaseType[] = [
        createRawMockPurchase({ id: "123e4567-e89b-12d3-a456-426614174000" }),
        createRawMockPurchase({ id: "123e4567-e89b-12d3-a456-426614174001" }),
      ];
      const requestPurchase = convertDecimalsToStrings(mockPurchases);

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => requestPurchase,
        } as Response),
      );

      const result = await fetchPurchases(apiHost, fakeToken, 1);
      const url = `${apiHost}/api/v2/purchases?createdById=1&status=confirmed`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "GET",
          headers: expect.objectContaining({
            "Content-Type": "application/json",
            Authorization: "Bearer " + fakeToken,
          }),
        }),
      );

      expect(result).toEqual(mockPurchases);
    });
  });

  describe("refundPurchaseById", () => {
    it("should return purchase on successful response", async () => {
      const purchaseId = "123e4567-e89b-12d3-a456-426614174000";
      const mockPurchase = createRawMockPurchase({
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

      const result = await refundPurchaseById(apiHost, fakeToken, purchaseId);
      const url = `${apiHost}/api/v2/purchases/${purchaseId}/refund`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "POST",
          headers: expect.objectContaining({
            "Content-Type": "application/json",
            Authorization: "Bearer " + fakeToken,
          }),
        }),
      );

      expect(result).toEqual(mockPurchase);
    });
  });

  describe("addProductInterest", () => {
    it("should return id on successful response", async () => {
      const productInterest = { id: 98 };

      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: true,
          json: async () => productInterest,
        } as Response),
      );

      const result = await addProductInterest(apiHost, fakeToken, 12);
      const url = `${apiHost}/api/v2/productInterests`;

      expect(fetch).toHaveBeenCalledWith(
        url,
        expect.objectContaining({
          method: "POST",
          headers: expect.objectContaining({
            "Content-Type": "application/json",
            Authorization: "Bearer " + fakeToken,
          }),
          body: JSON.stringify({ productId: 12 }),
        }),
      );

      expect(result).toEqual(productInterest);
    });
  });

  describe("error handling", () => {
    it("should throw an error if Zod validation fails", async () => {
      vi.mocked(fetch).mockResolvedValue({
        ok: true,
        json: async () => ({ wrong_key: "garbage" }),
      } as Response);

      await expect(fetchProducts(apiHost, fakeToken)).rejects.toThrow(
        "API Response format mismatch",
      );
    });

    it("should throw an error on non-ok response and call sentry", async () => {
      vi.stubGlobal(
        "fetch",
        vi.fn().mockResolvedValue({
          ok: false,
          status: 500,
          statusText: "Internal Server Error",
          json: async () => ({ message: "Server error" }),
        } as Response),
      );

      await expect(fetchProducts(apiHost, fakeToken)).rejects.toThrow(
        "Server error",
      );

      expect(Sentry.captureException).toHaveBeenCalled();
    });
  });
});
