import { AppConfig } from "../types/config.types";
import { RawConfig } from "../schemas/config.schemas";

export const transformConfig = (
  rawData: RawConfig,
  apiHost: string,
): AppConfig => {
  const sumupEnabled = rawData.paymentMethods.some((m) => m.code === "SUMUP");
  // Converts http:// -> ws:// and https:// -> wss:// for the websocket endpoint.
  const websocketHost = apiHost.replace(/^http/, "ws");

  const currencyOptions: Intl.NumberFormatOptions = {
    style: "currency",
    currency: rawData.currencyCode,
    minimumFractionDigits: rawData.fractionDigitsMin,
    maximumFractionDigits: rawData.fractionDigitsMax,
  };

  return {
    ...rawData,
    apiHost,
    websocketHost,
    sumupEnabled,
    currencyOptions,
    currency: new Intl.NumberFormat(rawData.currencyLocale, currencyOptions),
  };
};
