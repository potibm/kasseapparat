import * as Sentry from "@sentry/react";
import { z } from "zod";
import { refreshJwtToken } from "@core/api/auth";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Auth");

const API_HOST = import.meta.env.VITE_API_HOST ?? "https://localhost:3000";
const LOCALSTORAGE_KEY = "kasseapparat.admin.auth";

export const InitialSessionDataSchema = z.object({
  ID: z.number(),
  username: z.string(),
  role: z.string(),
  gravatarUrl: z.url(),
});

export const SessionDataSchema = InitialSessionDataSchema.extend({
  token: z.jwt(),
  expiryDate: z.iso.datetime(),
});

export type InitialSessionData = z.infer<typeof InitialSessionDataSchema>;

export type SessionData = z.infer<typeof SessionDataSchema>;

let refreshingPromise: Promise<void> | null = null;

export const initializeSession = (
  data: InitialSessionData,
  token: string,
  expiresIn: number,
): void => {
  const session: SessionData = {
    ...data,
    token,
    expiryDate: calculateExpiryDate(expiresIn),
  };

  storeSession(session);
};

export const updateSession = (token: string, expiresIn: number): void => {
  const currentData = getSession();
  if (!currentData) {
    throw new Error("No existing admin data found when updating session...");
  }
  const updatedSession: SessionData = {
    ...currentData,
    token,
    expiryDate: calculateExpiryDate(expiresIn),
  };
  storeSession(updatedSession);
};

const storeSession = (data: SessionData): void => {
  localStorage.setItem(LOCALSTORAGE_KEY, JSON.stringify(data));
};

const calculateExpiryDate = (expiresIn: number): string => {
  return new Date(Date.now() + (expiresIn - 30) * 1000).toISOString();
};

export const getSession = (): SessionData | null => {
  const rawData = localStorage.getItem(LOCALSTORAGE_KEY);
  if (!rawData) return null;

  try {
    const parsed = JSON.parse(rawData);
    const result = SessionDataSchema.safeParse(parsed);

    if (result.success) {
      return result.data;
    } else {
      log.warn("Invalid admin data in localStorage, clearing...");
      clearSession();
      return null;
    }
  } catch (error) {
    log.error("Error parsing admin data from localStorage:", error);
    return null;
  }
};

export const getSessionToken = (): string | null => {
  const session = getSession();
  return session ? session.token : null;
};

export const clearSession = (): void => {
  localStorage.removeItem(LOCALSTORAGE_KEY);
};

export const updateToken = async (): Promise<void> => {
  if (refreshingPromise) {
    return refreshingPromise;
  }

  refreshingPromise = (async (): Promise<void> => {
    try {
      const res = await refreshJwtToken(API_HOST);

      updateSession(res.access_token, res.expires_in);
    } catch (error: unknown) {
      // Sentry Logging
      Sentry.captureException(error, {
        tags: { auth: "refresh_token" },
      });

      clearSession();

      throw new Error("Token refresh error. Please log in again.", {
        cause: error,
      });
    } finally {
      refreshingPromise = null;
    }
  })();

  return refreshingPromise;
};
