import * as Sentry from "@sentry/react";
import { z } from "zod";

import {
  Product,
  ProductSchema,
  Purchase,
  PurchaseSchema,
  Guest,
  GuestSchema,
  ProductInterest,
  ProductInterestSchema,
} from "./api.schemas";

// Unified error handler for failed fetch responses
const handleFetchError = async (response: Response): Promise<never> => {
  let message = `HTTP ${response.status} ${response.statusText}`;
  try {
    const data = await response.json();
    message = data?.details || data?.error || data?.message || message;
  } catch {
    // Ignore invalid JSON
  }
  const error = new Error(message);
  Sentry.captureException(error, {
    extra: {
      url: response.url,
      status: response.status,
    },
  });
  throw error;
};

const postValidated = async <S extends z.ZodTypeAny>(
  url: string,
  token: string,
  body: object,
  schema: S,
): Promise<z.infer<S>> => {
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });
  if (!response.ok) await handleFetchError(response);

  const rawData = await response.json();

  const result = schema.safeParse(rawData);
  if (!result.success) {
    console.error("Zod Validation Error:", result.error);
    throw new Error("API Response format mismatch");
  }
  return result.data;
};

const getValidated = async <S extends z.ZodTypeAny>(
  url: string,
  token: string,
  schema: S,
): Promise<z.infer<S>> => {
  const response = await fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) await handleFetchError(response);

  const rawData = await response.json();

  const result = schema.safeParse(rawData);
  if (!result.success) {
    console.error("Zod Validation Error:", result.error);
    throw new Error("API Response format mismatch");
  }
  return result.data;
};

// Fetch all visible products
export const fetchProducts = async (
  apiHost: string,
  jwtToken: string,
): Promise<Product[]> => {
  const url = `${apiHost}/api/v2/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`;

  return getValidated(url, jwtToken, z.array(ProductSchema));
};

// Fetch guests for a specific product
export const fetchGuestlistByProductId = async (
  apiHost: string,
  jwtToken: string,
  productId: number,
  query: string,
): Promise<Guest[]> => {
  const url = `${apiHost}/api/v2/products/${productId}/guests?q=${query}`;

  const GuestListSchema = z.preprocess(
    (val) => (val === null ? [] : val),
    z.array(GuestSchema),
  );

  return getValidated(url, jwtToken, GuestListSchema);
};

// Store a new purchase
export const storePurchase = async (
  apiHost: string,
  jwtToken: string,
  payload: object,
): Promise<Purchase> => {
  return postValidated(
    `${apiHost}/api/v2/purchases`,
    jwtToken,
    payload,
    PurchaseSchema,
  );
};

// Fetch all confirmed purchases for a user
export const fetchPurchases = async (
  apiHost: string,
  jwtToken: string,
  userId: number,
): Promise<Purchase[]> => {
  const url = `${apiHost}/api/v2/purchases?createdById=${encodeURIComponent(userId)}&status=confirmed`;
  return getValidated(url, jwtToken, z.array(PurchaseSchema));
};

// Refund a purchase by ID
export const refundPurchaseById = async (
  apiHost: string,
  jwtToken: string,
  purchaseId: string,
): Promise<Purchase> => {
  const url = `${apiHost}/api/v2/purchases/${purchaseId}/refund`;
  return postValidated(url, jwtToken, {}, PurchaseSchema);
};

// Add interest in a product
export const addProductInterest = async (
  apiHost: string,
  jwtToken: string,
  productId: number,
): Promise<ProductInterest> => {
  return postValidated(
    `${apiHost}/api/v2/productInterests`,
    jwtToken,
    { productId },
    ProductInterestSchema,
  );
};
