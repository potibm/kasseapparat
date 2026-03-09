import * as Sentry from "@sentry/react";
import Decimal from "decimal.js";

import { Product } from "../features/product-list/types/product.types";
import { Guest } from "../features/guestlist/types/guest.types";
import {
  ApiGetResponseProduct,
  ApiCreateResponsePurchase,
  ApiGetResponsePurchase,
  ApiCreateResponseProductInterest,
} from "./api.types";

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

// Authenticated GET helper
const get = async <T>(url: string, token: string): Promise<T> => {
  const response = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
  if (!response.ok) await handleFetchError(response);
  return response.json();
};

// Authenticated POST helper
const post = async <T>(
  url: string,
  token: string,
  body: object,
): Promise<T> => {
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });
  if (!response.ok) await handleFetchError(response);
  return response.json();
};

// Fetch all visible products
export const fetchProducts = async (
  apiHost: string,
  jwtToken: string,
): Promise<Product[]> => {
  const url = `${apiHost}/api/v2/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`;
  const data = await get<ApiGetResponseProduct[]>(url, jwtToken);

  return data.map((p) => ({
    ...p,
    netPrice: new Decimal(p.netPrice),
    grossPrice: new Decimal(p.grossPrice),
    vatRate: new Decimal(p.vatRate),
    vatAmount: new Decimal(p.vatAmount),
  }));
};

// Fetch guests for a specific product
export const fetchGuestlistByProductId = async (
  apiHost: string,
  jwtToken: string,
  productId: number,
  query: string,
): Promise<Guest[]> => {
  const url = `${apiHost}/api/v2/products/${productId}/guests?q=${query}`;
  return get(url, jwtToken);
};

// Store a new purchase
export const storePurchase = async (
  apiHost: string,
  jwtToken: string,
  payload: object,
): Promise<ApiCreateResponsePurchase> => {
  return post(`${apiHost}/api/v2/purchases`, jwtToken, payload);
};

// Fetch all confirmed purchases for a user
export const fetchPurchases = async (
  apiHost: string,
  jwtToken: string,
  userId: number,
): Promise<ApiGetResponsePurchase[]> => {
  const url = `${apiHost}/api/v2/purchases?createdById=${encodeURIComponent(userId)}&status=confirmed`;
  return get(url, jwtToken);
};

// Refund a purchase by ID
export const refundPurchaseById = async (
  apiHost: string,
  jwtToken: string,
  purchaseId: string,
): Promise<ApiGetResponsePurchase> => {
  const url = `${apiHost}/api/v2/purchases/${purchaseId}/refund`;
  return post(url, jwtToken, {});
};

// Add interest in a product
export const addProductInterest = async (
  apiHost: string,
  jwtToken: string,
  productId: number,
): Promise<ApiCreateResponseProductInterest> => {
  return post(`${apiHost}/api/v2/productInterests`, jwtToken, { productId });
};
