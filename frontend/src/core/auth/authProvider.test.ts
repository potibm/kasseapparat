import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { createAuthProvider } from "./authProvider";

const mockFetch = vi.fn();
globalThis.fetch = mockFetch;

describe("authProvider", () => {
  const apiBaseUrl = "http://localhost:3001/api/v3";

  beforeEach(() => {
    vi.clearAllMocks();
    Object.defineProperty(window, "location", {
      value: {
        pathname: "/test",
        search: "?param=value",
        href: "http://localhost:3000/test?param=value",
      },
      writable: true,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("login", () => {
    it("should resolve without action", async () => {
      const provider = createAuthProvider(apiBaseUrl);
      await expect(provider.login!({})).resolves.toBeUndefined();
    });
  });

  describe("logout", () => {
    it("should call logout endpoint with correct URL", async () => {
      mockFetch.mockResolvedValueOnce({ ok: true });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.logout!({});

      expect(mockFetch).toHaveBeenCalledWith(`${apiBaseUrl}/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
    });

    it("should handle logout failure gracefully", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      const provider = createAuthProvider(apiBaseUrl);
      await expect(provider.logout!({})).resolves.toBeUndefined();
    });

    it("should clear cached user on logout", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });
      mockFetch.mockResolvedValueOnce({ ok: true });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkAuth!({});
      await provider.logout!({});

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });
      await provider.checkAuth!({});

      expect(mockFetch).toHaveBeenCalledTimes(3);
    });
  });

  describe("checkAuth", () => {
    it("should call auth/me endpoint with correct URL", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkAuth!({});

      expect(mockFetch).toHaveBeenCalledWith(`${apiBaseUrl}/auth/me`, {
        credentials: "include",
      });
    });

    it("should resolve on successful authentication", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);
      await expect(provider.checkAuth!({})).resolves.toBeUndefined();
    });

    it("should redirect to login on 401", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
      });

      const provider = createAuthProvider(apiBaseUrl);
      const promise = provider.checkAuth!({});

      await vi.waitFor(() => {
        expect(window.location.href).toContain(`${apiBaseUrl}/auth/login`);
        expect(window.location.href).toContain("returnTo=");
      });

      promise.catch(() => {});
    });

    it("should redirect to login on 403", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
      });

      const provider = createAuthProvider(apiBaseUrl);
      const promise = provider.checkAuth!({});

      await vi.waitFor(() => {
        expect(window.location.href).toContain(`${apiBaseUrl}/auth/login`);
      });

      promise.catch(() => {});
    });

    it("should reject on other HTTP errors", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
      });

      const provider = createAuthProvider(apiBaseUrl);
      await expect(provider.checkAuth!({})).rejects.toThrow("API Error");
    });

    it("should reject on network error", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      const provider = createAuthProvider(apiBaseUrl);
      await expect(provider.checkAuth!({})).rejects.toThrow("Network error");
    });

    it("should use cached user on subsequent calls", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkAuth!({});
      await provider.checkAuth!({});

      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it("should parse and cache user data", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "admin", role: "admin" }),
      });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkAuth!({});

      const permissions = await provider.getPermissions!({});
      expect(permissions).toBe("admin");
    });
  });

  describe("getPermissions", () => {
    it("should await a pending checkAuth instead of throwing", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);

      // getPermissions runs during the first render, before the in-flight
      // /auth/me request can resolve. It must wait for it rather than throw.
      const permissions = await provider.getPermissions!({});

      expect(permissions).toBe("user");
    });

    it("should redirect to login when auth/me is unauthorized", async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 401 });

      const provider = createAuthProvider(apiBaseUrl);
      const promise = provider.getPermissions!({});

      await vi.waitFor(() =>
        expect(window.location.href).toBe(
          `${apiBaseUrl}/auth/login?returnTo=%2Ftest%3Fparam%3Dvalue`,
        ),
      );

      // The promise deliberately never settles while the page navigates away.
      expect(promise).toBeInstanceOf(Promise);
    });

    it("should share one /auth/me request across concurrent callers", async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "admin" }),
      });

      const provider = createAuthProvider(apiBaseUrl);

      const [permissions, identity] = await Promise.all([
        provider.getPermissions!({}),
        provider.getIdentity!(),
        provider.checkAuth!({}),
      ]);

      expect(permissions).toBe("admin");
      expect(identity).toEqual({ id: "testuser", fullName: "testuser" });
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it("should return user role after authentication", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkAuth!({});

      const permissions = await provider.getPermissions!({});
      expect(permissions).toBe("user");
    });
  });

  describe("getIdentity", () => {
    it("should await a pending checkAuth instead of throwing", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);

      const identity = await provider.getIdentity!();

      expect(identity).toEqual({ id: "testuser", fullName: "testuser" });
    });

    it("should return user identity after authentication", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ username: "testuser", role: "user" }),
      });

      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkAuth!({});

      const identity = await provider.getIdentity!();
      expect(identity).toEqual({
        id: "testuser",
        fullName: "testuser",
      });
    });
  });

  describe("checkError", () => {
    it("should redirect to login on 401 error", async () => {
      const provider = createAuthProvider(apiBaseUrl);
      provider.checkError!({ status: 401 });

      await vi.waitFor(() => {
        expect(window.location.href).toContain(`${apiBaseUrl}/auth/login`);
      });
    });

    it("should redirect to login on 403 error", async () => {
      const provider = createAuthProvider(apiBaseUrl);
      provider.checkError!({ status: 403 });

      await vi.waitFor(() => {
        expect(window.location.href).toContain(`${apiBaseUrl}/auth/login`);
      });
    });

    it("should not redirect on other errors", async () => {
      const initialHref = window.location.href;
      const provider = createAuthProvider(apiBaseUrl);
      await provider.checkError!({ status: 500 });

      expect(window.location.href).toBe(initialHref);
    });
  });
});
