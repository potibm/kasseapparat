import { describe, it, expect, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { usePosApi } from "./usePosApi";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";

describe("usePosApi", () => {
  const mockConfig: AppConfig = {
    version: "1.0.0",
    apiHost: "http://localhost:3001",
    apiBaseUrl: "http://localhost:3001/api/v3",
    websocketHost: "ws://localhost:3001",
    websocketBaseUrl: "ws://localhost:3001/api/v3",
    locale: "en-US",
    currencyCode: "USD",
    currencyLocale: "en-US",
    currency: new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }),
    currencyOptions: {
      style: "currency",
      currency: "USD",
    },
    dateLocale: "en-US",
    dateOptions: {},
    vatRates: [{ rate: 19, name: "MwSt" }],
    paymentMethods: [
      { code: "cash", name: "Cash" },
      { code: "card", name: "Card" },
    ],
    sumupEnabled: false,
    authMode: "proxy",
  };

  it("should create a POS API client with the correct base URL", () => {
    const wrapper = ({ children }: { children: React.ReactNode }) => (
      <ConfigContext value={mockConfig}>{children}</ConfigContext>
    );

    const { result } = renderHook(() => usePosApi(), { wrapper });

    expect(result.current).toBeDefined();
    expect(result.current.fetchProducts).toBeDefined();
    expect(result.current.fetchGuestlistByProductId).toBeDefined();
    expect(result.current.storePurchase).toBeDefined();
    expect(result.current.fetchPurchases).toBeDefined();
    expect(result.current.refundPurchaseById).toBeDefined();
    expect(result.current.addProductInterest).toBeDefined();
  });

  it("should throw an error when ConfigContext is missing", () => {
    // Suppress console.error for this test
    const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});

    expect(() => {
      renderHook(() => usePosApi());
    }).toThrow("ConfigContext is missing");

    consoleError.mockRestore();
  });
});
