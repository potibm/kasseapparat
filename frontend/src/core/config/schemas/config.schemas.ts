import { z } from "zod";

const JsonObjectSchema = z
  .string()
  .transform((str, ctx) => {
    try {
      return JSON.parse(str);
    } catch {
      ctx.addIssue({
        code: "custom",
        message: "Invalid JSON",
      });
      return z.NEVER;
    }
  })
  .prefault("{}");

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
  dateOptions: JsonObjectSchema,
  vatRates: JsonObjectSchema,
  paymentMethods: z
    .array(
      z.object({
        code: z.string(),
        name: z.string(),
      }),
    )
    .default([]),
  locale: z.string().default("da-DK"),
});

export type RawConfig = z.infer<typeof ConfigSchema>;
