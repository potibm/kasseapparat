import { AppConfig } from "../types/config.types";
import { RawConfig } from "../schemas/config.schemas";
import { buildApiBaseUrl } from "../constants";

export const transformConfig = (
  rawData: RawConfig,
  apiHost: string,
): AppConfig => {
  const sumupEnabled = rawData.paymentMethods.some((m) => m.code === "SUMUP");
  const websocketHost = apiHost.replace(/^http/, "ws");
  const apiBaseUrl = buildApiBaseUrl(apiHost);
  const websocketBaseUrl = buildApiBaseUrl(websocketHost);

  const currencyOptions: Intl.NumberFormatOptions = {
    style: "currency",
    currency: rawData.currencyCode,
    minimumFractionDigits: rawData.fractionDigitsMin,
    maximumFractionDigits: rawData.fractionDigitsMax,
  };

  return {
    ...rawData,
    apiHost,
    apiBaseUrl,
    websocketHost,
    websocketBaseUrl,
    sumupEnabled,
    currencyOptions,
    currency: new Intl.NumberFormat(rawData.currencyLocale, currencyOptions),
  };
};
