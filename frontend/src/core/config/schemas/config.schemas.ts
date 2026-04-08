import { z } from "zod";

export const ConfigSchema = z.object({
  version: z.string().default("1.0.0"),
  sentryDSN: z.url().or(z.literal("")).optional(),
  sentryTraceSampleRate: z.number().min(0).max(1).optional(),
  sentryReplaySessionSampleRate: z.number().min(0).max(1).optional(),
  sentryReplayErrorSampleRate: z.number().min(0).max(1).optional(),
  currencyLocale: z.string().default("da-DK"),
  currencyCode: z.string().default("DKK"),
  fractionDigitsMin: z.number().default(0),
  fractionDigitsMax: z.number().default(2),
  dateLocale: z.string().default("en-US"),
  dateOptions: z.record(z.string(), z.any()).default({}),
  vatRates: z
    .array(
      z.object({
        rate: z.number(),
        name: z.string(),
      }),
    )
    .default([]),
  paymentMethods: z
    .array(
      z.object({
        code: z.string(),
        name: z.string(),
      }),
    )
    .default([]),
  /** @deprecated: Will be removed in favor of dateLocale or currencyLocale */
  locale: z.string().default("da-DK"),
});

export type RawConfig = z.infer<typeof ConfigSchema>;
