import { describe, it, expect } from "vitest";
import { JsonObjectSchema, DecimalSchema } from "./zod-schemas";
import Decimal from "decimal.js";

describe("JsonObjectSchema", () => {
  it("should parse a valid JSON object string correctly", () => {
    const input = '{"name": "Vitest", "awesome": true}';
    const result = JsonObjectSchema.safeParse(input);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data).toEqual({ name: "Vitest", awesome: true });
    }
  });

  it("should parse a valid JSON array string correctly", () => {
    const input = "[1, 2, 3]";
    const result = JsonObjectSchema.safeParse(input);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data).toEqual([1, 2, 3]);
    }
  });

  it("should fail for invalid JSON and provide a custom error", () => {
    const input = "kein-json-{]-;";
    const result = JsonObjectSchema.safeParse(input);

    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            code: "custom",
            message: "Invalid JSON",
          }),
        ]),
      );
    }
  });

  it("should use the default value when undefined is passed", () => {
    const result = JsonObjectSchema.safeParse(undefined);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data).toEqual({});
    }
  });

  it("should fail when the input is not a string (or undefined)", () => {
    const result = JsonObjectSchema.safeParse(12345);

    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].code).toBe("invalid_type");
    }
  });
});

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
