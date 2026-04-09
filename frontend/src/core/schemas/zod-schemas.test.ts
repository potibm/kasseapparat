import { describe, it, expect } from "vitest";
import { DecimalSchema } from "./zod-schemas";
import Decimal from "decimal.js";

describe("DecimalSchema", () => {
  it("should parse a valid decimal string correctly", () => {
    const input = "123.45";
    const result = DecimalSchema.safeParse(input);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data).toBeInstanceOf(Decimal);
      expect(result.data.toString()).toBe("123.45");
    }
  });

  it("should also handle negative numbers and integers", () => {
    const input = "-42";
    const result = DecimalSchema.safeParse(input);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.toString()).toBe("-42");
    }
  });

  it("should fail for invalid decimal strings and provide a custom error", () => {
    const input = "no-number";
    const result = DecimalSchema.safeParse(input);

    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            code: "custom",
            message: "Invalid decimal value",
            params: {
              input: "no-number",
              format: "decimal",
            },
          }),
        ]),
      );
    }
  });

  it("should fail when the input is not a string (or undefined)", () => {
    const result = DecimalSchema.safeParse(123.45);

    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].code).toBe("invalid_type");
    }
  });
});
