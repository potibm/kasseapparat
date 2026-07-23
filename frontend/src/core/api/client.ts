import * as Sentry from "@sentry/react";
import { z } from "zod";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Api");

// Unified error handler for failed fetch responses
export const handleFetchError = async (response: Response): Promise<never> => {
  let message = `HTTP ${response.status} ${response.statusText}`;
  try {
    const data = await response.json();
    message = data?.details || data?.error || data?.message || message;
  } catch {
    // Ignore invalid JSON
  }
  const error = new Error(message);
  const path = response.url ? new URL(response.url).pathname : undefined;
  Sentry.captureException(error, {
    extra: {
      status: response.status,
      path,
    },
  });
  log.warn("API request failed", { status: response.status, path, message });
  throw error;
};

export const postValidated = async <S extends z.ZodTypeAny>(
  url: string,
  body: object,
  schema: S,
): Promise<z.infer<S>> => {
  return performFetch(url, schema, "POST", body);
};

export const getValidated = async <S extends z.ZodTypeAny>(
  url: string,
  schema: S,
): Promise<z.infer<S>> => {
  return performFetch(url, schema, "GET");
};

export const performFetch = async <S extends z.ZodTypeAny>(
  url: string,
  schema: S,
  method: string,
  body?: object,
  options?: RequestInit,
): Promise<z.infer<S>> => {
  const response = await fetch(url, {
    method,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
    credentials: options?.credentials ?? "include",
  });
  if (!response.ok) await handleFetchError(response);

  const rawData = await response.json();

  const result = schema.safeParse(rawData);
  if (!result.success) {
    log.error("Zod Validation Error", result.error);
    throw new Error("API Response format mismatch");
  }
  return result.data;
};
