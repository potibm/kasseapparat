import { describe, it, expect } from "vitest";
import { transformConfig } from "./config.transform"; // adjust path as needed
import { RawConfig } from "../schemas/config.schemas";

describe("transformConfig", () => {
  // Mock data for a standard configuration
  const mockRawData: RawConfig = {
    version: "1.0.0",
    currencyCode: "EUR",
    currencyLocale: "de-DE",
    fractionDigitsMin: 2,
    fractionDigitsMax: 2,
    paymentMethods: [{ code: "CREDIT_CARD", name: "Credit Card" }],
    locale: "de-DE",
    dateLocale: "de-DE",
    dateOptions: {},
    vatRates: [],
  };

  const mockApiHost = "https://api.example.com";

  it("should correctly assign apiHost and transform websocketHost", () => {
    const result = transformConfig(mockRawData, mockApiHost);

    expect(result.apiHost).toBe("https://api.example.com");
    expect(result.websocketHost).toBe("wss://api.example.com");
  });

  it("should transform http to ws for insecure connections", () => {
    const result = transformConfig(mockRawData, "http://localhost:3000");
    expect(result.websocketHost).toBe("ws://localhost:3000");
  });

  it("should set sumupEnabled to true when SUMUP is present in paymentMethods", () => {
    const dataWithSumup: RawConfig = {
      ...mockRawData,
      paymentMethods: [{ code: "SUMUP", name: "SumUp Terminal" }],
    };

    const result = transformConfig(dataWithSumup, mockApiHost);
    expect(result.sumupEnabled).toBe(true);
  });

  it("should set sumupEnabled to false when SUMUP is missing", () => {
    const result = transformConfig(mockRawData, mockApiHost);
    expect(result.sumupEnabled).toBe(false);
  });

  it("should create a correctly configured Intl.NumberFormat instance", () => {
    const result = transformConfig(mockRawData, mockApiHost);

    // Test formatting (EUR in de-DE)
    // Note: We use a regex to handle different space characters (like non-breaking spaces)
    const formatted = result.currency.format(1234.5);
    expect(formatted).toMatch(/1.234,50\s*€/);

    expect(result.currencyOptions.currency).toBe("EUR");
    expect(result.currencyOptions.minimumFractionDigits).toBe(2);
  });

  it("should preserve all original rawData fields using the spread operator", () => {
    const result = transformConfig(mockRawData, mockApiHost);
    expect(result.version).toBe("1.0.0");
    expect(result.locale).toBe("de-DE");
    expect(result.dateLocale).toBe("de-DE");
  });
});
