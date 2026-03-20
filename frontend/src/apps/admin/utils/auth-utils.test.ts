import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import * as Sentry from "@sentry/react";
import {
  initializeSession,
  getSession,
  updateToken,
  getSessionToken,
} from "./auth-utils";
import { refreshJwtToken } from "@core/api/auth";
import { faker } from "@faker-js/faker";

vi.mock("@sentry/react", () => ({
  captureException: vi.fn(),
}));

vi.mock("@core/api/auth", () => ({
  refreshJwtToken: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

const MOCK_INITIAL_DATA = {
  ID: 1,
  username: "admin",
  role: "superadmin",
  gravatarUrl: "https://gravatar.com/avatar/123",
};

const VALID_TOKEN = faker.internet.jwt();
const LOCALSTORAGE_KEY = "kasseapparat.admin.auth";

describe("Auth Session Management", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    // Set fixed system time: Jan 1st 2024, 12:00:00 PM UTC
    vi.setSystemTime(new Date("2024-01-01T12:00:00Z"));
    localStorage.clear();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe("initializeSession & calculateExpiryDate", () => {
    it("should calculate the session expiry correctly and store it in localStorage", () => {
      const expiresIn = 3600; // 1 hour
      initializeSession(MOCK_INITIAL_DATA, VALID_TOKEN, expiresIn);

      const stored = JSON.parse(localStorage.getItem(LOCALSTORAGE_KEY)!);

      // Calculation: Now + (3600 - 30) seconds
      // 12:00:00 + 3570s = 12:59:30
      expect(stored.expiryDate).toBe("2024-01-01T12:59:30.000Z");
      expect(stored.token).toBe(VALID_TOKEN);
      expect(stored.username).toBe("admin");
    });
  });

  describe("getSession", () => {
    it("should return null if there is no data in storage", () => {
      expect(getSession()).toBeNull();
    });

    it("should return the session data if the stored data is valid", () => {
      initializeSession(MOCK_INITIAL_DATA, VALID_TOKEN, 3600);
      const session = getSession();
      expect(session).not.toBeNull();
      expect(session?.username).toBe("admin");
    });

    it("should clear storage and return null if stored data is invalid JSON", () => {
      localStorage.setItem(LOCALSTORAGE_KEY, "not-a-json");

      const session = getSession();
      expect(session).toBeNull();
      expect(localStorage.getItem(LOCALSTORAGE_KEY)).toBeNull();
    });

    it("should clear storage and return null if Zod validation fails", () => {
      // Manipulated data (e.g., missing ID)
      localStorage.setItem(
        LOCALSTORAGE_KEY,
        JSON.stringify({ username: "hack" }),
      );

      const session = getSession();
      expect(session).toBeNull();
      expect(localStorage.getItem(LOCALSTORAGE_KEY)).toBeNull();
    });
  });

  describe("updateToken (Async Refresh Logic)", () => {
    it("should successfully refresh the token", async () => {
      // Initialize session first
      initializeSession(MOCK_INITIAL_DATA, VALID_TOKEN, 3600);

      const newToken = faker.internet.jwt();
      expect(newToken).not.toBe(VALID_TOKEN);

      // Prepare API mock
      vi.mocked(refreshJwtToken).mockResolvedValue({
        access_token: newToken,
        refresh_token: "new-refresh-token",
        token_type: "Bearer",
        expires_in: 3600,
      });

      await updateToken();

      expect(getSessionToken()).toBe(newToken);
      expect(refreshJwtToken).toHaveBeenCalled();
    });

    it("should throw an error when no previous session exists", async () => {
      vi.mocked(refreshJwtToken).mockResolvedValue({
        access_token: faker.internet.jwt(),
        refresh_token: "new-refresh-token",
        token_type: "Bearer",
        expires_in: 3600,
      });

      try {
        await updateToken();
        expect.fail("Should have thrown an error");
      } catch (error: unknown) {
        expect(error).toBeInstanceOf(Error);
        if (error instanceof Error) {
          expect(error.message).toBe(
            "Token refresh error. Please log in again.",
          );

          expect(error.cause).toBeInstanceOf(Error);
          if (error.cause instanceof Error) {
            expect(error.cause.message).toBe(
              "No existing admin data found when updating session...",
            );
          }
        }
      }
    });

    it("should capture exception with Sentry and clear the session on error", async () => {
      initializeSession(MOCK_INITIAL_DATA, VALID_TOKEN, 3600);
      const error = new Error("Network fail");
      vi.mocked(refreshJwtToken).mockRejectedValue(error);

      await expect(updateToken()).rejects.toThrow("Token refresh error");

      expect(Sentry.captureException).toHaveBeenCalledWith(
        error,
        expect.any(Object),
      );
      expect(getSession()).toBeNull(); // Ensure session was cleared
    });

    it("should prevent race conditions by using a singleton promise", async () => {
      initializeSession(MOCK_INITIAL_DATA, VALID_TOKEN, 3600);

      type RefreshResponse = Awaited<ReturnType<typeof refreshJwtToken>>;

      let resolveApi!: (value: RefreshResponse) => void;
      const apiPromise = new Promise<RefreshResponse>((res) => {
        resolveApi = res;
      });

      vi.mocked(refreshJwtToken).mockReturnValue(apiPromise);

      // Start two calls simultaneously
      const call1 = updateToken();
      const call2 = updateToken();

      // It must return the exact same Promise instance
      expect(call1).toStrictEqual(call2);

      // Resolve the API mock
      const newToken = faker.internet.jwt();
      expect(newToken).not.toBe(VALID_TOKEN);
      resolveApi({
        access_token: newToken,
        refresh_token: "new-refresh-token",
        token_type: "Bearer",
        expires_in: 60,
      });
      await call1;

      // The actual API call should only happen once
      expect(refreshJwtToken).toHaveBeenCalledTimes(1);
    });
  });
});
