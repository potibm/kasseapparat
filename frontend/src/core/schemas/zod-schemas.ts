import { z } from "zod";
import Decimal from "decimal.js";

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
