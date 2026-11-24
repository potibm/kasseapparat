export type OutOfStockBehavior = "ignore" | "fail" | "auto_sold_out";

export interface AppConfig {
  version: string;
  apiHost: string;
  apiBaseUrl: string;
  websocketHost: string;
  websocketBaseUrl: string;
  sentryDSN?: string;
  sentryTraceSampleRate?: number;
  sentryReplaySessionSampleRate?: number;
  sentryReplayErrorSampleRate?: number;
  sentryEnvironment?: string;
  locale: string;
  currencyCode: string;
  currencyLocale: string;
  currency: Intl.NumberFormat;
  currencyOptions: Intl.NumberFormatOptions;
  dateLocale: string;
  dateOptions: Intl.DateTimeFormatOptions;
  vatRates: Array<{ rate: number; name: string }>;
  paymentMethods: Array<{
    code: string;
    name: string;
    changeCalculator?: boolean;
  }>;
  sumupEnabled: boolean;
  environmentMessage?: string;
  authMode: string;
  outOfStockBehavior: OutOfStockBehavior;
}
