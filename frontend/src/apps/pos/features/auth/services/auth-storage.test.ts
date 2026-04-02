import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import {
  getInitialSession,
  getInitialUser,
  storeSession,
  storeUser,
  clearAuthStorage,
  AUTH_KEYS,
} from "./auth-storage";
import { AuthUser as AuthUserType } from "../types/auth.types";
import { faker } from "@faker-js/faker";
import { createMockUserData } from "@core/api/auth.schemas.mock";

// mocks
const { mockDebug, mockWarn, mockError } = vi.hoisted(() => ({
  mockDebug: vi.fn(),
  mockWarn: vi.fn(),
  mockError: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    debug: mockDebug,
    warn: mockWarn,
    error: mockError,
  }),
}));

const VALID_JWT_MOCK = faker.internet.jwt();

describe("Auth Storage Service", () => {
  let getItemSpy: ReturnType<typeof vi.spyOn>;
  let setItemSpy: ReturnType<typeof vi.spyOn>;
  let removeItemSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    vi.clearAllMocks();

    localStorage.clear();

    getItemSpy = vi.spyOn(Storage.prototype, "getItem");
    setItemSpy = vi.spyOn(Storage.prototype, "setItem");
    removeItemSpy = vi.spyOn(Storage.prototype, "removeItem");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("getInitialSession()", () => {
    it("should return the session if localStorage data is valid", () => {
      const mockDateStr = "2025-01-01T12:00:00.000Z";
      localStorage.setItem(AUTH_KEYS.TOKEN, VALID_JWT_MOCK);
      localStorage.setItem(AUTH_KEYS.EXPIRY, mockDateStr);

      const session = getInitialSession();

      expect(session.token).toBe(VALID_JWT_MOCK);
      expect(session.expiryDate).toBeInstanceOf(Date);
      expect(session.expiryDate?.toISOString()).toBe(mockDateStr);
      expect(mockDebug).toHaveBeenCalledWith("LocalStorage Session restored", {
        token: VALID_JWT_MOCK,
        expiryDate: session.expiryDate,
      });
    });

    it("should return nulls if Zod parsing fails (e.g. invalid token format)", () => {
      localStorage.setItem(AUTH_KEYS.TOKEN, "invalid-token");
      localStorage.setItem(AUTH_KEYS.EXPIRY, "2025-01-01T12:00:00.000Z");

      const session = getInitialSession();

      expect(session).toEqual({ token: null, expiryDate: null });
    });

    it("should return nulls if an error is thrown (catch block)", () => {
      getItemSpy.mockImplementation(() => {
        throw new Error("Access to Storage denied");
      });

      const session = getInitialSession();

      expect(session).toEqual({ token: null, expiryDate: null });
    });
  });

  describe("getInitialUser()", () => {
    const mockUserData = createMockUserData();

    it("should return the user if localStorage data is valid and schema passes", () => {
      localStorage.setItem(AUTH_KEYS.USER, JSON.stringify(mockUserData));

      const user = getInitialUser();

      expect(user).toEqual(mockUserData);
      expect(mockDebug).toHaveBeenCalledWith(
        "LocalStorage Userdata restored",
        mockUserData,
      );
    });

    it("should return null if there is no data in localStorage", () => {
      const user = getInitialUser();
      expect(user).toBeNull();
    });

    it("should return null and warn if JSON is valid but Zod validation fails", () => {
      localStorage.setItem(AUTH_KEYS.USER, JSON.stringify({ wrong: "data" }));

      const user = getInitialUser();

      expect(user).toBeNull();
      expect(mockWarn).toHaveBeenCalledWith(
        "LocalStorage Userdata invalid. Clearing...",
      );
    });

    it("should return null if JSON parsing throws an error (catch block)", () => {
      // corrupt json
      localStorage.setItem(AUTH_KEYS.USER, "invalid-json-{");

      const user = getInitialUser();

      expect(user).toBeNull();
    });
  });

  describe("storeSession()", () => {
    it("should store the token and the ISO string of the expiry date", () => {
      const mockDate = new Date("2025-01-01T12:00:00.000Z");

      storeSession(VALID_JWT_MOCK, mockDate);

      expect(setItemSpy).toHaveBeenCalledWith(AUTH_KEYS.TOKEN, VALID_JWT_MOCK);
      expect(setItemSpy).toHaveBeenCalledWith(
        AUTH_KEYS.EXPIRY,
        "2025-01-01T12:00:00.000Z",
      );
    });
  });

  describe("storeUser()", () => {
    it("should stringify and store the user data", () => {
      const mockUser = {
        id: 42,
        username: "test_user",
      } as unknown as AuthUserType;

      storeUser(mockUser);

      expect(setItemSpy).toHaveBeenCalledWith(
        AUTH_KEYS.USER,
        JSON.stringify(mockUser),
      );
    });
  });

  describe("clearAuthStorage()", () => {
    it("should remove all keys defined in AUTH_KEYS from localStorage", () => {
      clearAuthStorage();

      Object.values(AUTH_KEYS).forEach((key) => {
        expect(removeItemSpy).toHaveBeenCalledWith(key);
      });

      // as we have 3 keys in AUTH_KEYS, it should be called exactly 3 times
      expect(removeItemSpy).toHaveBeenCalledTimes(
        Object.keys(AUTH_KEYS).length,
      );
    });
  });
});
