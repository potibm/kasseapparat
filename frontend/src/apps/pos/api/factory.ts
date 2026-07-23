import { getValidated, postValidated } from "@core/api/client";
import {
  ProductSchema,
  GuestSchema,
  Guest,
  Purchase,
  PurchaseSchema,
  ProductInterest,
  ProductInterestSchema,
} from "./schemas";
import { ApiCreatePayloadPurchase } from "./types";
import { z } from "zod";

export const createPosApiClient = (apiBaseUrl: string) => {
  const buildUrl = (endpoint: string) => {
    const base = apiBaseUrl.replace(/\/$/, "");
    const path = endpoint.replace(/^\//, "");
    return `${base}/${path}`;
  };

  return {
    // Fetch all visible products
    fetchProducts: async () => {
      const url = buildUrl(
        "/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true",
      );
      return getValidated(url, z.array(ProductSchema));
    },
    // Fetch guests for a specific product
    fetchGuestlistByProductId: async (
      productId: number,
      query: string,
    ): Promise<Guest[]> => {
      const url = buildUrl(
        `/products/${productId}/guests?q=${encodeURIComponent(query)}`,
      );
      const GuestListSchema = z.preprocess(
        (val) => (val === null ? [] : val),
        z.array(GuestSchema),
      );
      return getValidated(url, GuestListSchema);
    },
    // Store a new purchase
    storePurchase: async (
      payload: ApiCreatePayloadPurchase,
    ): Promise<Purchase> => {
      return postValidated(buildUrl("/purchases"), payload, PurchaseSchema);
    },
    // Fetch all confirmed and pending purchases for a user
    fetchPurchases: async (username: string): Promise<Purchase[]> => {
      const url = buildUrl(
        `/purchases?createdById=${encodeURIComponent(username)}&status=confirmed&status=pending`,
      );
      return getValidated(url, z.array(PurchaseSchema));
    },
    // Refund a purchase by ID
    refundPurchaseById: async (purchaseId: string): Promise<Purchase> => {
      const url = buildUrl(`/purchases/${purchaseId}/refund`);
      return postValidated(url, {}, PurchaseSchema);
    },
    // Add interest in a product
    addProductInterest: async (productId: number): Promise<ProductInterest> => {
      return postValidated(
        buildUrl(`/productInterests`),
        { productId },
        ProductInterestSchema,
      );
    },
  };
};
