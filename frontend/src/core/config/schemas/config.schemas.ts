import { z } from "zod";

export const ConfigSchema = z.object({
  version: z.string().default("1.0.0"),
  sentryDSN: z.url().or(z.literal("")).optional(),
  sentryTraceSampleRate: z.number().min(0).max(1).optional(),
  sentryReplaySessionSampleRate: z.number().min(0).max(1).optional(),
  sentryReplayErrorSampleRate: z.number().min(0).max(1).optional(),
  currencyLocale: z.string().default("dk-DK"),
  currencyCode: z.string().default("DKK"),
  fractionDigitsMin: z.number().default(0),
  fractionDigitsMax: z.number().default(2),
  dateLocale: z.string().default("en-US"),
  dateOptions: z
    .string()
    .transform((str) => {
      try {
        return JSON.parse(str);
      } catch {
        return {};
      }
    })
    .default("{}"),
  vatRates: z
    .string()
    .transform((str) => {
      try {
        return JSON.parse(str);
      } catch {
        return {};
      }
    })
    .default("{}"),
  paymentMethods: z
    .array(
      z.object({
        code: z.string(),
        name: z.string(),
      }),
    )
    .optional()
    .default([]),
  locale: z.string().default("dk-DK"),
});

export type RawConfig = z.infer<typeof ConfigSchema>;
