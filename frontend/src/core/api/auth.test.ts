import { describe, it, expect, vi, beforeEach } from "vitest";
import * as Sentry from "@sentry/react";
import { getJwtToken, ApiError, changePassword } from "./auth"; // adjust path

const minimalValidJwt =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.8VKCTiBegJPuPIZlp0wbV0Sbdn5BS6TE5DCx6oYNc5o";

// 1. Mock External Dependencies
vi.mock("@sentry/react", () => ({
  captureException: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    error: vi.fn(),
  }),
}));

describe("Auth Service", () => {
  const apiHost = "https://api.example.com";

  beforeEach(() => {
    vi.restoreAllMocks();
    // Stub global fetch
    vi.stubGlobal("fetch", vi.fn());
  });

  describe("getJwtToken", () => {
    it("should return data when the response is successful and valid", async () => {
      const mockResponse = {
        access_token: minimalValidJwt,
        token_type: "Bearer",
        expires_in: 3600,
        role: "user",
        username: "john_doe",
        gravatarUrl: "https://www.gravatar.com/avatar/123",
        id: 123,
      };

      // Mock fetch success
      vi.mocked(fetch).mockResolvedValue({
        ok: true,
        json: async () => mockResponse,
      } as Response);

      const result = await getJwtToken(apiHost, "user", "pass");

      expect(fetch).toHaveBeenCalledWith(
        `${apiHost}/api/v2/auth/login`,
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({ login: "user", password: "pass" }),
        }),
      );
      expect(result).toEqual(mockResponse);
    });

    it('should throw an ApiError and NOT notify Sentry for "invalid credentials"', async () => {
      vi.mocked(fetch).mockResolvedValue({
        ok: false,
        status: 401,
        url: ".../login",
        json: async () => ({ message: "Invalid credentials" }),
      } as Response);

      await expect(getJwtToken(apiHost, "u", "p")).rejects.toThrow(ApiError);

      // Should not report common auth errors to Sentry
      expect(Sentry.captureException).not.toHaveBeenCalled();
    });

    it("should report to Sentry for unexpected server errors (500)", async () => {
      vi.mocked(fetch).mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        url: ".../login",
        json: async () => ({ error: "Database connection failed" }),
      } as Response);

      await expect(getJwtToken(apiHost, "u", "p")).rejects.toThrow(
        "Database connection failed",
      );

      // Sentry should be called for critical/unexpected errors
      expect(Sentry.captureException).toHaveBeenCalled();
    });

    it("should throw an error if Zod validation fails", async () => {
      vi.mocked(fetch).mockResolvedValue({
        ok: true,
        json: async () => ({ wrong_key: "garbage" }), // Invalid according to LoginResponseSchema
      } as Response);

      await expect(getJwtToken(apiHost, "u", "p")).rejects.toThrow(
        "Invalid response format from Auth API",
      );
    });
  });

  describe("changePassword", () => {
    it("should return data when the response is successful and valid", async () => {
      const mockResponse = {
        message: "Password changed successfully",
        status: "success",
      };

      vi.mocked(fetch).mockResolvedValue({
        ok: true,
        json: async () => mockResponse,
      } as Response);

      await changePassword(apiHost, "123", "token", "pass");

      const callBody = JSON.parse(
        vi.mocked(fetch).mock.calls[0][1]?.body as string,
      );
      expect(callBody.userId).toBe(123); // string '123' became number 123
    });
  });
});
