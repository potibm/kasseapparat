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
  dateLocale: string;
  dateOptions: Intl.DateTimeFormatOptions;
  vatRates: Record<string, number>;
  paymentMethods: Array<{ code: string; name: string }>;
  sumupEnabled: boolean;
  environmentMessage?: string;
}
