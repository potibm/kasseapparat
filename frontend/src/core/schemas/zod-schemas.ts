import { z } from "zod";
import Decimal from "decimal.js";

export const JsonObjectSchema = z
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

export const DecimalSchema = z.string().transform((val, ctx) => {
  try {
    return new Decimal(val);
  } catch {
    ctx.addIssue({
      code: "custom",
      message: "Invalid decimal value",
      params: {
        input: val,
        format: "decimal",
      },
    });
    return z.NEVER;
  }
});
