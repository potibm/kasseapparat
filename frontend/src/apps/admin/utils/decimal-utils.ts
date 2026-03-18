import Decimal from "decimal.js";

export const parseDecimal = (value: string | null): string | null => {
  if (typeof value !== "string" || !value) return null;

  const cleaned = value
    .trim()
    .replace(",", ".")
    .replaceAll(/[^\d.]/g, "");
  const [integerPart, ...fractionalParts] = cleaned.split(".");

  return fractionalParts.length > 0
    ? `${integerPart}.${fractionalParts.join("")}`
    : cleaned;
};

export const formatDecimal = (value: unknown): string => {
  if (typeof value === "string" || typeof value === "number") {
    return String(value).replace(".", ",");
  }
  return "";
};

export const decimalValidator = (value: unknown) => {
  if (value === null || value === undefined || value === "") return undefined;
  if (typeof value !== "string" && typeof value !== "number")
    return "Invalid input type";

  try {
    const decimalValue = new Decimal(value);
    if (decimalValue.isNaN()) return "Invalid number";
    if (decimalValue.isNegative()) return "Negative number";
    return undefined;
  } catch {
    return "Invalid number";
  }
};
