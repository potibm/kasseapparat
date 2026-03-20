import { describe, it, expect, beforeEach, vi } from "vitest";
import { fetchUtils } from "react-admin";
import * as Sentry from "@sentry/react";
import { GetListParams, GetOneParams, CreateParams } from "react-admin";

// 1. ADD THIS IMPORT: We need it to read the arguments it was called with
import jsonServerProvider from "ra-data-json-server";

// 2. USE vi.hoisted(): This tells Vitest to hoist these variables
// alongside the vi.mock() calls, preventing the ReferenceError.
const { mockBaseProvider } = vi.hoisted(() => ({
  mockBaseProvider: {
    getList: vi.fn(),
    getOne: vi.fn(),
    getMany: vi.fn(),
    getManyReference: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    updateMany: vi.fn(),
    delete: vi.fn(),
    deleteMany: vi.fn(),
  },
}));

// 3. SIMPLIFIED MOCK: Just return the hoisted provider
vi.mock("ra-data-json-server", () => ({
  default: vi.fn(() => mockBaseProvider),
}));

vi.mock("react-admin", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-admin")>();
  return {
    ...actual,
    fetchUtils: {
      fetchJson: vi.fn(),
    },
    addRefreshAuthToDataProvider: vi.fn((provider) => provider),
  };
});

vi.mock("@sentry/react", () => ({
  captureException: vi.fn(),
}));

vi.mock("../utils/auth-utils", () => ({
  getSessionToken: vi.fn(),
}));

vi.mock("./refresh-token", () => ({
  refreshToken: vi.fn(),
}));

// 4. IMPORTANT: Import the file under test AFTER the mocks are defined
import dataProvider from "./data-provider";
import { getSessionToken } from "../utils/auth-utils";

// 5. EXTRACT THE CLIENT: Since dataProvider is now imported and executed,
// jsonServerProvider was called. We simply extract the 2nd argument (index 1).

interface HttpClientOptions extends fetchUtils.Options {
  isUpload?: boolean;
}

type HttpClient = (
  url: string,
  options?: HttpClientOptions,
) => ReturnType<typeof fetchUtils.fetchJson>;
const capturedHttpClient = vi.mocked(jsonServerProvider).mock
  .calls[0][1] as HttpClient;

const mockFetchResponse = (jsonPayload: unknown) =>
  ({ json: jsonPayload }) as Awaited<ReturnType<typeof fetchUtils.fetchJson>>;

// --- 2. TESTS ---

describe("Data Provider", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Resource Resolution & Alias Mapping", () => {
    it("should map sumupReaders to sumup/readers when calling base methods", () => {
      dataProvider.getList("sumupReaders", {
        pagination: { page: 1, perPage: 10 },
      } as GetListParams);

      expect(mockBaseProvider.getList).toHaveBeenCalledWith(
        "sumup/readers",
        expect.any(Object),
      );
    });

    it("should leave unknown resources untouched", () => {
      dataProvider.getOne("products", { id: 1 } as GetOneParams);

      expect(mockBaseProvider.getOne).toHaveBeenCalledWith(
        "products",
        expect.any(Object),
      );
    });
  });

  describe("httpClient (Internal Fetch Wrapper)", () => {
    it("should set default headers and Content-Type", async () => {
      vi.mocked(fetchUtils.fetchJson).mockResolvedValue(mockFetchResponse({}));

      await capturedHttpClient("http://localhost:3001/api/v2/test");

      const calledOptions = vi.mocked(fetchUtils.fetchJson).mock.calls[0][1];
      const headers = calledOptions?.headers as Headers;

      expect(headers.get("Content-Type")).toBe("application/json");
      expect(headers.get("Accept")).toBe("application/json");
    });

    it("should skip Content-Type if isUpload is true", async () => {
      vi.mocked(fetchUtils.fetchJson).mockResolvedValue(mockFetchResponse({}));

      await capturedHttpClient("http://localhost:3001/api/v2/test", {
        isUpload: true,
      });

      const calledOptions = vi.mocked(fetchUtils.fetchJson).mock.calls[0][1];
      const headers = calledOptions?.headers as Headers;
      expect(headers.get("Content-Type")).toBeNull();
    });

    it("should inject the Authorization token if a session exists", async () => {
      vi.mocked(getSessionToken).mockReturnValue("valid-test-token");
      vi.mocked(fetchUtils.fetchJson).mockResolvedValue(mockFetchResponse({}));

      await capturedHttpClient("http://localhost:3001/api/v2/test");

      const calledOptions = vi.mocked(fetchUtils.fetchJson).mock.calls[0][1];
      const headers = calledOptions?.headers as Headers;
      expect(headers.get("Authorization")).toBe("Bearer valid-test-token");
    });
  });

  describe("httpClient Error Handling & Sentry", () => {
    it("should throw and report critical errors to Sentry", async () => {
      const error = new Error("Database connection failed");
      vi.mocked(fetchUtils.fetchJson).mockRejectedValue(error);
      vi.mocked(getSessionToken).mockReturnValue("secret-token");

      await expect(
        capturedHttpClient("http://localhost:3001/api/v2/test"),
      ).rejects.toThrow("Database connection failed");

      expect(Sentry.captureException).toHaveBeenCalledTimes(1);
    });

    it("should NOT report known non-critical errors to Sentry", async () => {
      const expectedError = new Error("Cookie token is empty"); // Should be caught by the filter
      vi.mocked(fetchUtils.fetchJson).mockRejectedValue(expectedError);

      await expect(
        capturedHttpClient("http://localhost:3001/api/v2/test"),
      ).rejects.toThrow("Cookie token is empty");

      // The filter logic should prevent Sentry from being called
      expect(Sentry.captureException).not.toHaveBeenCalled();
    });

    it("should scrub the Authorization header before sending extra data to Sentry", async () => {
      const error = new Error("Some critical API fail");
      vi.mocked(fetchUtils.fetchJson).mockRejectedValue(error);
      vi.mocked(getSessionToken).mockReturnValue(
        "super-secret-token-do-not-log",
      );

      await expect(
        capturedHttpClient("http://localhost:3001/api/v2/test"),
      ).rejects.toThrow();

      // Check the exact payload sent to Sentry
      expect(Sentry.captureException).toHaveBeenCalledWith(
        error,
        expect.objectContaining({
          extra: expect.objectContaining({
            request: expect.objectContaining({
              headers: expect.not.objectContaining({
                Authorization: expect.any(String),
              }),
            }),
          }),
        }),
      );
    });
  });

  describe("Custom Methods (upload & refund)", () => {
    it("should handle upload correctly", async () => {
      vi.mocked(fetchUtils.fetchJson).mockResolvedValue(
        mockFetchResponse({
          success: true,
        }),
      );

      const fileBody = new FormData();
      const result = await dataProvider.upload("images", {
        data: fileBody as unknown,
      } as CreateParams);

      expect(result.data).toEqual({ success: true });

      // Verify httpClient was called with correct parameters
      const [url, options] = vi.mocked(fetchUtils.fetchJson).mock.calls[0];
      const headers = options?.headers as Headers;

      expect(url).toBe("http://localhost:3001/api/v2/images");
      expect(options?.method).toBe("POST");
      expect(options?.body).toBe(fileBody);
      // isUpload flag is stripped before passing to fetchJson, but we know it bypassed Content-Type
      expect(headers.get("Content-Type")).toBeNull();
    });

    it("should throw an error if refund is called on a non-purchase resource", async () => {
      await expect(
        dataProvider.refund("products", { id: 1, data: {} } as GetOneParams), // Invalid resource
      ).rejects.toThrow("Refund is not supported for resource: products");
    });

    it("should throw an error if refund is called without an id", async () => {
      await expect(
        dataProvider.refund("purchases", { data: {} }), // Missing ID
      ).rejects.toThrow("Refund requires an id");
    });

    it("should perform a refund successfully", async () => {
      vi.mocked(fetchUtils.fetchJson).mockResolvedValue(
        mockFetchResponse({
          refunded: true,
        }),
      );

      const payload = { reason: "defective" };
      const result = await dataProvider.refund("purchases", {
        id: 123,
        data: payload,
      } as GetOneParams);

      expect(result.data).toEqual({ refunded: true });

      const [url, options] = vi.mocked(fetchUtils.fetchJson).mock.calls[0];
      expect(url).toBe("http://localhost:3001/api/v2/purchases/123/refund");
      expect(options?.method).toBe("POST");
      expect(options?.body).toBe(JSON.stringify(payload));
    });
  });
});
