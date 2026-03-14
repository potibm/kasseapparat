export interface AppConfig {
  version: string;
  apiHost: string;
  websocketHost: string;
  sentryDSN?: string;
  sentryTraceSampleRate?: number;
  sentryReplaySessionSampleRate?: number;
  sentryReplayErrorSampleRate?: number;
  locale: string;
  currencyCode: string;
  currencyLocale: string;
  currency: Intl.NumberFormat;
  currencyOptions: Intl.NumberFormatOptions;
  dateLocale: string;
  dateOptions: Intl.DateTimeFormatOptions;
  vatRates: Array<{ rate: number; name: string }>;
  paymentMethods: Array<{ code: string; name: string }>;
  sumupEnabled: boolean;
  environmentMessage?: string;
}
