import { describe, it, expect, beforeEach, vi } from "vitest";
import * as Sentry from "@sentry/react";
import { getJwtToken, logout as logoutApi } from "@core/api/auth";
import {
  getSession,
  clearSession,
  initializeSession,
} from "../utils/auth-utils";
import {
  LoginResponse as LoginResponseType,
  SimpleResponse,
} from "@core/api/auth.schemas";

// --- 1. HOIST ENVIRONMENT VARIABLES ---
vi.hoisted(() => {
  vi.stubEnv("VITE_API_HOST", "https://localhost:3000");
});

// --- 2. MOCKS ---
vi.mock("@sentry/react", () => ({
  setUser: vi.fn(),
  captureException: vi.fn(),
}));

vi.mock("@core/api/auth", () => ({
  getJwtToken: vi.fn(),
  logout: vi.fn(),
}));

vi.mock("../utils/auth-utils", () => ({
  getSession: vi.fn(),
  clearSession: vi.fn(),
  initializeSession: vi.fn(),
}));

// Mock react-admin so the wrapper just returns the raw provider for easy testing
vi.mock("react-admin", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-admin")>();
  return {
    ...actual,
    addRefreshAuthToAuthProvider: vi.fn((provider) => provider),
  };
});

vi.mock("./refresh-token", () => ({
  refreshToken: vi.fn(),
}));

// Import the provider AFTER the mocks are set up
import authProvider from "./auth-provider";

// --- 3. TYPES & FIXTURES ---
type SessionData = NonNullable<ReturnType<typeof getSession>>;

const mockLoginResponse = {
  id: 123,
  username: "admin_user",
  role: "superadmin",
  gravatarUrl: "https://gravatar.com/avatar/test",
  access_token: "secret-jwt-token",
  expires_in: 3600,
  // Add any other required fields from LoginResponseType here or cast via unknown
} as unknown as LoginResponseType;

const mockSessionData = {
  ID: 123,
  username: "admin_user",
  role: "superadmin",
  gravatarUrl: "https://gravatar.com/avatar/test",
} as unknown as SessionData;

// --- 4. TESTS ---
describe("Auth Provider", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("login()", () => {
    it("should authenticate successfully, set Sentry user, and initialize session", async () => {
      vi.mocked(getJwtToken).mockResolvedValue(mockLoginResponse);

      await authProvider.login({ username: "test", password: "password123" });

      // 1. Checks if the API was called with the correct host and credentials
      expect(getJwtToken).toHaveBeenCalledWith(
        "https://localhost:3000",
        "test",
        "password123",
      );

      // 2. Checks Sentry initialization
      expect(Sentry.setUser).toHaveBeenCalledWith({
        id: "123",
        username: "admin_user",
      });

      // 3. Checks if the session was correctly mapped and initialized
      expect(initializeSession).toHaveBeenCalledWith(
        {
          ID: 123,
          username: "admin_user",
          role: "superadmin",
          gravatarUrl: "https://gravatar.com/avatar/test",
        },
        "secret-jwt-token",
        3600,
      );
    });

    it("should throw a formatted error with a cause if login fails", async () => {
      const apiError = new Error("Invalid credentials");
      vi.mocked(getJwtToken).mockRejectedValue(apiError);

      try {
        await authProvider.login({ username: "test", password: "wrong" });
        expect.fail("Should have thrown an error");
      } catch (error: unknown) {
        expect(error).toBeInstanceOf(Error);
        const err = error as Error;
        expect(err.message).toBe(
          "There was an error logging you in: Invalid credentials",
        );
        expect(err.cause).toBe(apiError);
      }
    });
  });

  describe("logout()", () => {
    it("should call the logout API and clear the session", async () => {
      vi.mocked(logoutApi).mockResolvedValue({} as SimpleResponse);

      await authProvider.logout({});

      expect(logoutApi).toHaveBeenCalledWith("https://localhost:3000");
      expect(clearSession).toHaveBeenCalledTimes(1);
    });

    it("should capture API errors to Sentry but ALWAYS clear the session (finally block)", async () => {
      const error = new Error("Network timeout during logout");
      vi.mocked(logoutApi).mockRejectedValue(error);

      await authProvider.logout({});

      // Check Sentry tagging
      expect(Sentry.captureException).toHaveBeenCalledWith(error, {
        tags: { auth: "logout" },
      });
      // The most important part: ensure clearSession ran despite the API failure
      expect(clearSession).toHaveBeenCalledTimes(1);
    });
  });

  describe("checkError()", () => {
    it("should throw an error and clear session on 401 Unauthorized", async () => {
      await expect(authProvider.checkError({ status: 401 })).rejects.toThrow(
        "Authentication error. Please log in again.",
      );
      expect(clearSession).toHaveBeenCalledTimes(1);
    });

    it("should throw an error and clear session on 403 Forbidden", async () => {
      await expect(authProvider.checkError({ status: 403 })).rejects.toThrow();
      expect(clearSession).toHaveBeenCalledTimes(1);
    });

    it("should resolve and do nothing for other status codes (e.g. 500, 404)", async () => {
      await expect(
        authProvider.checkError({ status: 500 }),
      ).resolves.toBeUndefined();
      expect(clearSession).not.toHaveBeenCalled();
    });
  });

  describe("checkAuth()", () => {
    it("should resolve successfully if a session exists", async () => {
      vi.mocked(getSession).mockReturnValue(mockSessionData);
      await expect(authProvider.checkAuth({})).resolves.toBeUndefined();
    });

    it("should throw an error if no session exists", async () => {
      vi.mocked(getSession).mockReturnValue(null);
      await expect(authProvider.checkAuth({})).rejects.toThrow(
        "No session found",
      );
    });
  });

  describe("getPermissions()", () => {
    it("should return the user role if a session exists", async () => {
      vi.mocked(getSession).mockReturnValue(mockSessionData);
      const permissions = await authProvider.getPermissions!({});
      expect(permissions).toBe("superadmin");
    });

    it("should throw an error if no session exists", async () => {
      vi.mocked(getSession).mockReturnValue(null);
      await expect(authProvider.getPermissions!({})).rejects.toThrow(
        "No session found",
      );
    });
  });

  describe("getIdentity()", () => {
    it("should format and return the UserIdentity if a session exists", async () => {
      vi.mocked(getSession).mockReturnValue(mockSessionData);
      const identity = await authProvider.getIdentity?.();

      expect(identity).toEqual({
        id: 123,
        fullName: "admin_user",
        avatar: "https://gravatar.com/avatar/test",
      });
    });

    it("should throw an error if no session exists", async () => {
      vi.mocked(getSession).mockReturnValue(null);
      await expect(authProvider.getIdentity?.()).rejects.toThrow(
        "No session found",
      );
    });
  });
});
